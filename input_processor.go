package prompt

import (
	"context"
	"log"
	"time"
)

type InputProcessor struct {
	in ConsoleParser
}

func (ip *InputProcessor) Run(ctx context.Context, bufCh chan []byte) {
	log.Printf("[INFO] Start running input processor")
	for {
		select {
		case <-ctx.Done():
			log.Print("[INFO] Stop input processor")
			return
		default:
			if b, err := ip.in.Read(); err == nil && !(len(b) == 1 && b[0] == 0) {
				bufCh <- b
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
