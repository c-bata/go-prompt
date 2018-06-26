// +build !windows

package prompt

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type SignalHandler struct {
	SigWinch chan struct{}
}

func NewSignalHandler() *SignalHandler {
	return &SignalHandler{
		SigWinch: make(chan struct{}),
	}
}

func (sh *SignalHandler) Run(ctx context.Context, cancel context.CancelFunc) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(
		sigchan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	sigwinch := make(chan os.Signal, 1)
	signal.Notify(sigwinch, syscall.SIGWINCH)

	for {
		select {
		case <-ctx.Done():
			log.Print("[INFO] SignalHandler: stop by context")
			return
		case s := <-sigchan:
			log.Printf("[INFO] SignalHandler: stop by %v", s.String())
			cancel()
		case <-sigwinch:
			sh.SigWinch <- struct{}{}
		}
	}
}
