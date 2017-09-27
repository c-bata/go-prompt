// +build !windows

package prompt

import (
	"log"
	"syscall"
	"time"
)

func (p *Prompt) readBuffer(bufCh chan []byte, stopCh chan struct{}) {
	buf := make([]byte, 1024)

	log.Printf("[INFO] readBuffer start")
	for {
		time.Sleep(10 * time.Millisecond)
		select {
		case <-stopCh:
			log.Print("[INFO] stop readBuffer")
			return
		default:
			if n, err := syscall.Read(syscall.Stdin, buf); err == nil {
				bufCh <- buf[:n]
			}
		}
	}
}
