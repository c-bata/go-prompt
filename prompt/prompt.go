package prompt

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Executor func(*Buffer) string
type Completer func(*Buffer) []string

type Prompt struct {
	in        *VT100Parser
	buf       *Buffer
	renderer  *Render
	title     string
	executor  Executor
	completer Completer
}

func (p *Prompt) Run() {
	p.setUp()
	defer p.tearDown()

	bufCh := make(chan []byte, 128)
	go readBuffer(bufCh)

	exitCh := make(chan bool, 16)
	winSizeCh := make(chan *WinSize, 128)
	go handleSignals(p.in, exitCh, winSizeCh)

	for {
		select {
		case b := <-bufCh:
			p.renderer.Erase(p.buf)
			ac := p.in.GetASCIICode(b)
			if ac == nil {
				p.buf.InsertText(string(b), false, true)
			} else if ac.Key == ControlJ || ac.Key == Enter {
				res := p.executor(p.buf)
				p.renderer.BreakLine(p.buf, res)
				p.buf = NewBuffer()
			} else if ac.Key == ControlC || ac.Key == ControlD {
				return
			} else {
				InputHandler(ac, p.buf)
			}

			completions := p.completer(p.buf)
			p.renderer.Render(p.buf, completions)
		case w := <-winSizeCh:
			p.renderer.UpdateWinSize(w)
		case e := <-exitCh:
			if e {
				return
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (p *Prompt) setUp() {
	p.in.Setup()
	p.renderer.Setup()
}

func (p *Prompt) tearDown() {
	p.in.TearDown()
	p.renderer.TearDown()
}

func readBuffer(bufCh chan []byte) {
	buf := make([]byte, 1024)

	for {
		if n, err := syscall.Read(syscall.Stdin, buf); err == nil {
			bufCh <- buf[:n]
		}
	}
}

func handleSignals(in *VT100Parser, exitCh chan bool, winSizeCh chan *WinSize) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGWINCH,
	)

	for {
		s := <-sigCh
		switch s {
		// kill -SIGHUP XXXX
		case syscall.SIGHUP:
			exitCh <- true

			// kill -SIGINT XXXX or Ctrl+c
		case syscall.SIGINT:
			exitCh <- true

			// kill -SIGTERM XXXX
		case syscall.SIGTERM:
			exitCh <- true

			// kill -SIGQUIT XXXX
		case syscall.SIGQUIT:
			exitCh <- true

		case syscall.SIGWINCH:
			winSizeCh <- in.GetWinSize()
		default:
		}
	}
}

func NewPrompt(executor Executor, completer Completer, maxCompletions uint8) *Prompt {
	return &Prompt{
		in: NewVT100Parser(),
		renderer: &Render{
			Prefix:         ">>> ",
			out:            NewVT100Writer(),
			maxCompletions: maxCompletions,
		},
		title:     "Hello! this is prompt toolkit",
		buf:       NewBuffer(),
		executor:  executor,
		completer: completer,
	}
}
