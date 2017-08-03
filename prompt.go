package prompt

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	logfile      = "/tmp/go-prompt-debug.log"
	envEnableLog = "GO_PROMPT_ENABLE_LOG"
)

type Executor func(context.Context, string) string
type Completer func(string) []Completion
type Completion struct {
	Text        string
	Description string
}

type Prompt struct {
	in             ConsoleParser
	buf            *Buffer
	renderer       *Render
	executor       Executor
	completer      Completer
	maxCompletions uint16
	selected       int // -1 means nothing one is selected.
	history        []string
}

func (p *Prompt) Run() {
	p.setUp()
	defer p.tearDown()

	if os.Getenv(envEnableLog) != "true" {
		log.SetOutput(ioutil.Discard)
	} else if f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.SetOutput(ioutil.Discard)
	} else {
		defer f.Close()
		log.SetOutput(f)
		log.Println("Logging is enabled.")
	}

	p.renderer.Render(p.buf, p.completer(p.buf.Text()), p.maxCompletions, p.selected)

	bufCh := make(chan []byte, 128)
	go readBuffer(bufCh)

	exitCh := make(chan struct{})
	winSizeCh := make(chan *WinSize)
	go handleSignals(p.in, exitCh, winSizeCh)

	for {
		select {
		case b := <-bufCh:
			if shouldExecute, shouldExit, input := p.feed(b); shouldExit {
				return
			} else if shouldExecute {
				ctx, _ := context.WithCancel(context.Background())
				p.renderer.RenderResult(p.executor(ctx, input))

				completions := p.completer(p.buf.Text())
				p.updateSelectedCompletion(completions)
				p.renderer.Render(p.buf, completions, p.maxCompletions, p.selected)
			} else {
				completions := p.completer(p.buf.Text())
				p.updateSelectedCompletion(completions)
				p.renderer.Render(p.buf, completions, p.maxCompletions, p.selected)
			}
		case w := <-winSizeCh:
			p.renderer.UpdateWinSize(w)
			completions := p.completer(p.buf.Text())
			p.renderer.Render(p.buf, completions, p.maxCompletions, p.selected)
		case <-exitCh:
			return
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (p *Prompt) feed(b []byte) (shouldExecute, shouldExit bool, input string) {
	shouldExecute = false
	key := p.in.GetKey(b)

	switch key {
	case ControlJ, Enter:
		if p.selected != -1 {
			c := p.completer(p.buf.Text())[p.selected]
			w := p.buf.Document().GetWordBeforeCursor()
			if w != "" {
				p.buf.DeleteBeforeCursor(len([]rune(w)))
			}
			p.buf.InsertText(c.Text, false, true)
		}
		p.renderer.BreakLine(p.buf)

		shouldExecute = true
		input = p.buf.Text()
		p.history = append(p.history, input)
		log.Printf("History: %s", input)
		p.buf = NewBuffer()
		p.selected = -1
	case ControlC:
		p.renderer.BreakLine(p.buf)
		p.buf = NewBuffer()
		p.selected = -1
	case ControlD:
		shouldExit = true
	case Up:
		fallthrough
	case BackTab:
		p.selected -= 1
	case ControlI:
		fallthrough
	case Down:
		fallthrough
	case Tab:
		p.selected += 1
	case Left:
		p.buf.CursorLeft(1)
	case Right:
		p.buf.CursorRight(1)
	case Backspace:
		p.buf.DeleteBeforeCursor(1)
	case NotDefined:
		if p.selected != -1 {
			c := p.completer(p.buf.Text())[p.selected]
			w := p.buf.Document().GetWordBeforeCursor()
			if w != "" {
				p.buf.DeleteBeforeCursor(len([]rune(w)))
			}
			p.buf.InsertText(c.Text, false, true)
		}
		p.selected = -1
		p.buf.InsertText(string(b), false, true)
	default:
		p.selected = -1
	}
	return
}

func (p *Prompt) updateSelectedCompletion(completions []Completion) {
	max := int(p.maxCompletions)
	if len(completions) < max {
		max = len(completions)
	}
	if p.selected >= max {
		p.selected = -1
	} else if p.selected < -1 {
		p.selected = max - 1
	}
}

func (p *Prompt) setUp() {
	p.in.Setup()
	p.renderer.Setup()
	p.renderer.UpdateWinSize(p.in.GetWinSize())
	p.selected = -1 // -1 means nothing one is selected.
}

func (p *Prompt) tearDown() {
	p.in.TearDown()
	p.renderer.TearDown()
}

func readBuffer(bufCh chan []byte) {
	buf := make([]byte, 1024)

	for {
		if n, err := syscall.Read(syscall.Stdin, buf); err == nil {
			bufCh <- buf[:n]
		}
	}
}

func handleSignals(in ConsoleParser, exitCh chan struct{}, winSizeCh chan *WinSize) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGWINCH,
	)

	for {
		s := <-sigCh
		switch s {
		// kill -SIGHUP XXXX
		case syscall.SIGHUP:
			exitCh <- struct{}{}

			// kill -SIGINT XXXX or Ctrl+c
		case syscall.SIGINT:
			exitCh <- struct{}{}

			// kill -SIGTERM XXXX
		case syscall.SIGTERM:
			exitCh <- struct{}{}

			// kill -SIGQUIT XXXX
		case syscall.SIGQUIT:
			exitCh <- struct{}{}

		case syscall.SIGWINCH:
			winSizeCh <- in.GetWinSize()
		default:
		}
	}
}
