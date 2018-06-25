package prompt

import (
	"context"
	"log"
	"time"
)

type InputProcessor struct {
	UserInput chan []byte
	in        ConsoleParser
}

func NewInputProcessor(in ConsoleParser) *InputProcessor {
	return &InputProcessor{
		in:        in,
		UserInput: make(chan []byte, 128),
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
		default:
			if b, err := ip.in.Read(); err == nil && !(len(b) == 1 && b[0] == 0) {
				ip.UserInput <- b
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
