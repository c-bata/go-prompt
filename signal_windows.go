// +build windows

package prompt

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func handleSignals(ctx context.Context, cancel context.CancelFunc, winsizechan chan struct{}) {
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
				cancel()

			case syscall.SIGQUIT: // kill -SIGQUIT XXXX
				cancel()
			}
		}
	}
}
