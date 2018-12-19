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
	Input chan rune
	Out   chan KeyPress
}

func (p *Parser) Feed(r rune) {
	p.Input <- r
}

func (p *Parser) Start() {
	prefix := bytes.NewBuffer(nil)
	retry := false
	for {
		if retry {
			retry = false
		} else {
			in, ok := <-p.Input
			if !ok {
				break
			}
			prefix.WriteRune(in)
		}

		if prefix.Len() > 0 {
			isPrefixOfLongerMatch := false
			for _, s := range prompt.ASCIISequences {
				if bytes.HasPrefix(s.ASCIICode, prefix.Bytes()) && bytes.Compare(s.ASCIICode, prefix.Bytes()) != 0 {
					isPrefixOfLongerMatch = true
					break
				}
			}
			match := prompt.GetKey(prefix.Bytes())
			if !isPrefixOfLongerMatch && match != prompt.NotDefined {
				p.Out <- KeyPress{Key: match, Bytes: prefix.Bytes()}
				prefix.Reset()
			} else if !isPrefixOfLongerMatch && match == prompt.NotDefined {
				found := false
				retry = true

				prefixRunes := []rune(prefix.String())
				for i := len(prefixRunes); i > 0; i-- {
					prefixBytes := []byte(string(prefixRunes[:i]))
					match = prompt.GetKey(prefixBytes)
					if match != prompt.NotDefined {
						p.Out <- KeyPress{Key: match, Bytes: prefixBytes}
						for j := 0; j < i; j++ {
							_, _, _ = prefix.ReadRune()
						}
						found = true
					}
				}
				if !found {
					r, _, err := prefix.ReadRune()
					if err == io.EOF {
						continue // don't reach here.
					}
					p.Out <- KeyPress{Key: match, Bytes: []byte(string(r))}
					_, _, _ = prefix.ReadRune()
				}
			}
		}
	}
}

func NewParser() *Parser {
	return &Parser{
		Input: make(chan rune, 1),
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
			p.Feed(c)
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
