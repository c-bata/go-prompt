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
		syscall.SIGWINCH,
	)

	for {
		select {
		case <-ctx.Done():
			log.Println("[INFO] stop handleSignals")
			return
		case s := <-sigchan:
			switch s {
			case syscall.SIGINT: // kill -SIGINT XXXX or Ctrl+c
				log.Println("[SIGNAL] Catch SIGINT")
				cancel()

			case syscall.SIGTERM: // kill -SIGTERM XXXX
				log.Println("[SIGNAL] Catch SIGTERM")
				cancel()

			case syscall.SIGQUIT: // kill -SIGQUIT XXXX
				log.Println("[SIGNAL] Catch SIGQUIT")
				cancel()

			case syscall.SIGWINCH:
				log.Println("[SIGNAL] Catch SIGWINCH")
				sh.SigWinch <- struct{}{}
			}
		}
	}
}
