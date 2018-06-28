// +build !windows

package prompt

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/pkg/term/termios"
)

const maxReadBytes = 1024

// PosixParser is a ConsoleParser implementation for POSIX environment.
type PosixParser struct {
	fd          int
	origTermios syscall.Termios
	initTermios sync.Once
}

// SetUp should be called before starting input
func (t *PosixParser) SetUp() error {
	_, _, e := syscall.Syscall(syscall.SYS_FCNTL, uintptr(t.fd), uintptr(syscall.F_SETFL),
		uintptr(syscall.O_ASYNC|syscall.O_NONBLOCK))
	if e != 0 {
		log.Printf("[ERROR] Cannot set non-blocking mode: %d\n", e)
		return e
	}

	_, _, e = syscall.Syscall(syscall.SYS_FCNTL, uintptr(t.fd), uintptr(syscall.F_SETOWN),
		uintptr(syscall.Getpid()))
	if runtime.GOOS != "darwin" && e != 0 {
		log.Printf("[ERROR] Cannot set F_SETOWN: %d\n", e)
		return e
	}

	if err := t.setRawMode(); err != nil {
		log.Printf("[ERROR] Cannot set raw mode: %s", err)
		return err
	}
	return nil
}

// TearDown should be called after stopping input
func (t *PosixParser) TearDown() error {
	if err := syscall.SetNonblock(t.fd, false); err != nil {
		log.Println("[ERROR] Cannot set blocking mode.")
		return err
	}
	if err := t.restoreTermios(); err != nil {
		log.Println("[ERROR] Cannot reset from raw mode.")
		return err
	}
	return nil
}

// Read returns byte array.
func (t *PosixParser) Read() ([]byte, error) {
	buf := make([]byte, maxReadBytes)
	n, err := syscall.Read(t.fd, buf)
	if err != nil {
		return []byte{}, err
	}
	return buf[:n], nil
}

func (t *PosixParser) setRawMode() error {
	var err error
	t.initTermios.Do(func() {
		log.Println("termios is initialized")
		err = termios.Tcgetattr(uintptr(t.fd), &t.origTermios)
	})
	if err != nil {
		return err
	}

	n := t.origTermios
	n.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK |
		syscall.ISTRIP | syscall.INLCR | syscall.IGNCR |
		syscall.ICRNL | syscall.IXON
	n.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG | syscall.ECHONL
	n.Cflag &^= syscall.CSIZE | syscall.PARENB
	n.Cc[syscall.VMIN] = 1
	n.Cc[syscall.VTIME] = 0
	return termios.Tcsetattr(uintptr(t.fd), termios.TCSANOW, &n)
}

func (t *PosixParser) restoreTermios() error {
	return termios.Tcsetattr(uintptr(t.fd), termios.TCSANOW, &t.origTermios)
}

// NewStandardInputParser returns ConsoleParser object to read from stdin.
func NewStandardInputParser() ConsoleParser {
	in, err := syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	return &PosixParser{
		fd: in,
	}
}

// Run start to worker goroutine.
func (ip *InputProcessor) Run(ctx context.Context) (err error) {
	log.Printf("[INFO] InputProcessor: Start running input processor")
	defer log.Print("[INFO] InputProcessor: Stop input processor")
	sigio := make(chan os.Signal, 1)
	signal.Notify(sigio, syscall.SIGIO)

	ip.in.SetUp()
	defer ip.in.TearDown()
	for {
		select {
		case <-ctx.Done():
			return
		case p := <-ip.Pause:
			if ip.pause == p {
				return // distinct until changed
			}
			ip.pause = p

			if ip.pause {
				ip.in.TearDown()
			} else {
				ip.in.SetUp()
			}
		case <-sigio:
			b, err := ip.in.Read()
			if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
				continue
			} else if err != nil {
				log.Printf("[ERROR] cannot read %s", err)
				return err
			}
			if !(len(b) == 1 && b[0] == 0) {
				ip.UserInput <- b
			}
		}
	}
}
