package prompt

import (
	"context"
	"log"
	"sync"
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

func (ip *InputProcessor) Run(ctx context.Context, wg *sync.WaitGroup) {
	log.Printf("[INFO] InputProcessor: Start running input processor")
	wg.Add(1)
	defer wg.Done()
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
