package prompt

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type InputProcessor struct {
	UserInput chan []byte
	Pause     chan bool
	in        ConsoleParser
	pause     bool
}

func NewInputProcessor(in ConsoleParser) *InputProcessor {
	return &InputProcessor{
		in:        in,
		UserInput: make(chan []byte, 128),
		Pause:     make(chan bool),
	}
}

func (ip *InputProcessor) Run(ctx context.Context) (err error) {
	log.Printf("[INFO] InputProcessor: Start running input processor")
	defer log.Print("[INFO] InputProcessor: Stop input processor")
	sigio := make(chan os.Signal, 1)
	signal.Notify(sigio, syscall.SIGIO)

	ip.in.Setup()
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
				ip.in.Setup()
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
