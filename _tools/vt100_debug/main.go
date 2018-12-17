// +build !windows

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	prompt "github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/internal/term"
)

type KeyPress struct {
	Key   prompt.Key
	Bytes []byte
}

type Parser struct {
	Input chan []byte
	Out   chan KeyPress
}

func (p *Parser) Feed(data []byte) {
	p.Input <- data
}

func (p *Parser) process(prefix []byte, retry bool, in []byte) ([]byte, bool) {
	if retry {
		retry = false
	} else {
		prefix = append(prefix, in...)
	}
	if len(prefix) > 0 {
		isPrefixOfLongerMatch := false
		for _, s := range prompt.ASCIISequences {
			if bytes.HasPrefix(s.ASCIICode, prefix) && bytes.Compare(s.ASCIICode, prefix) != 0 {
				isPrefixOfLongerMatch = true
				break
			}
		}
		match := false
		matchedKey := prompt.GetKey(prefix)
		if matchedKey != prompt.NotDefined {
			match = true
		}
		if !isPrefixOfLongerMatch && match {
			p.Out <- KeyPress{Key: matchedKey, Bytes: prefix}
			prefix = make([]byte, 0, 5)
		} else if !isPrefixOfLongerMatch && !match {
			found := false
			retry = true

			for i := len(prefix); i > 0; i-- {
				matchedKey = prompt.GetKey(prefix[:i])
				if matchedKey != prompt.NotDefined {
					match = true
					break
				}

				if match {
					p.Out <- KeyPress{Key: matchedKey, Bytes: prefix}
					prefix = prefix[i:]
					found = true
				}
			}
			if !found {
				p.Out <- KeyPress{Key: matchedKey, Bytes: prefix}
				prefix = prefix[1:]
			}
		}
	}
	return prefix, retry
}

func (p *Parser) Start() {
	prefix := make([]byte, 0, 5)
	retry := false
	for {
		select {
		case in, ok := <-p.Input:
			if !ok {
				return
			}
			prefix, retry = p.process(prefix, retry, in)
			if retry {
				prefix, retry = p.process(prefix, retry, in)
			}
		}
	}
}

func NewParser() *Parser {
	return &Parser{
		Input: make(chan []byte, 1),
		Out:   make(chan KeyPress, 1),
	}
}

func main() {
	if err := term.SetRaw(syscall.Stdin); err != nil {
		fmt.Println(err)
		return
	}
	defer term.Restore()

	p := NewParser()
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go p.Start()

	r := bufio.NewReader(os.Stdin)
	go func() {
		for {
			c, _, err := r.ReadRune()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					break
				}
			}
			p.Feed([]byte(string(c)))
		}
		close(p.Input)
	}()

	for {
		select {
		case keypress := <-p.Out:
			if keypress.Key == prompt.NotDefined {
				fmt.Printf("Key '%s' data:'%#v'\n", string(keypress.Bytes), keypress.Bytes)
			} else {
				fmt.Printf("Key '%s' data:'%#v'\n", keypress.Key, keypress.Bytes)
				if keypress.Key == prompt.ControlC {
					return
				}
			}
			if keypress.Key == prompt.ControlC {
				return
			}
		case <-sigquit:
			close(p.Input)
			return
		}
	}
}
