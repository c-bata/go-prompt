// +build windows

package prompt

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func (p *Prompt) handleSignals(tx context.Context, cancel context.CancelFunc, winSizeCh chan *WinSize) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	for {
		select {
		case <-ctx.Done():
			return
		case s := <-sigCh:
			switch s {
			case syscall.SIGINT: // kill -SIGINT XXXX or Ctrl+c
				cancel()

			case syscall.SIGTERM: // kill -SIGTERM XXXX
				log.Println("[SIGNAL] Catch SIGTERM")

			case syscall.SIGQUIT: // kill -SIGQUIT XXXX
				cancel()
			}
		}
	}
}
