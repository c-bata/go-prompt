package prompt

import (
	"context"
	"log"
	"time"
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

func (ip *InputProcessor) Run(ctx context.Context) {
	log.Printf("[INFO] InputProcessor: Start running input processor")
	defer log.Print("[INFO] InputProcessor: Stop input processor")

	ip.in.Setup()
	defer ip.in.TearDown()
	for {
		select {
		case <-ctx.Done():
			return
		case p := <-ip.Pause:
			if p == ip.pause {
				continue
			}
			ip.pause = p

			if ip.pause {
				ip.in.TearDown()
			} else {
				ip.in.Setup()
			}
		default:
			if !ip.pause {
				if b, err := ip.in.Read(); err == nil && !(len(b) == 1 && b[0] == 0) {
					ip.UserInput <- b
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
