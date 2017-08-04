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
type Completer func(string) []Suggest

type Prompt struct {
	in             ConsoleParser
	buf            *Buffer
	renderer       *Render
	executor       Executor
	completer      Completer
	history        *History
	completion     *CompletionManager
}

type Exec struct {
	input string
	ctx   context.Context
}

func (e *Exec) Context() context.Context {
	if e.ctx == nil {
		e.ctx = context.Background()
	}
	return e.ctx
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
		log.Println("[INFO] Logging is enabled.")
	}

	p.renderer.Render(p.buf, p.completer(p.buf.Text()), p.completion.Max, p.completion.selected)

	bufCh := make(chan []byte, 128)
	go readBuffer(bufCh)

	exitCh := make(chan int)
	winSizeCh := make(chan *WinSize)
	go handleSignals(p.in, exitCh, winSizeCh)

	for {
		select {
		case b := <-bufCh:
			if shouldExit, exec := p.feed(b); shouldExit {
				return
			} else if exec != nil {
				p.runExecutor(exec, bufCh)

				completions := p.completer(p.buf.Text())
				p.completion.update(completions)
				p.renderer.Render(p.buf, completions, p.completion.Max, p.completion.selected)
			} else {
				completions := p.completer(p.buf.Text())
				p.completion.update(completions)
				p.renderer.Render(p.buf, completions, p.completion.Max, p.completion.selected)
			}
		case w := <-winSizeCh:
			p.renderer.UpdateWinSize(w)
			completions := p.completer(p.buf.Text())
			p.renderer.Render(p.buf, completions, p.completion.Max, p.completion.selected)
		case code := <-exitCh:
			p.tearDown()
			os.Exit(code)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (p *Prompt) runExecutor(exec *Exec, bufCh chan []byte) {
	resCh := make(chan string, 1)
	ctx, cancel := context.WithCancel(exec.Context())
	go func() {
		resCh <- p.executor(ctx, exec.input)
	}()

	for {
		select {
		case r := <-resCh:
			p.renderer.RenderResult(r)
			return
		case b := <-bufCh:
			if p.in.GetKey(b) == ControlC {
				log.Println("[INFO] Executor is canceled.")
				cancel()
			}
		}
	}
	return
}

func (p *Prompt) feed(b []byte) (shouldExit bool, exec *Exec) {
	key := p.in.GetKey(b)

	switch key {
	case ControlJ, Enter:
		if p.completion.Completing() {
			c := p.completer(p.buf.Text())[p.completion.selected]
			w := p.buf.Document().GetWordBeforeCursor()
			if w != "" {
				p.buf.DeleteBeforeCursor(len([]rune(w)))
			}
			p.buf.InsertText(c.Text, false, true)
		}
		p.renderer.BreakLine(p.buf)

		exec = &Exec{input: p.buf.Text()}
		log.Printf("[History] %s", p.buf.Text())
		p.buf = NewBuffer()
		p.completion.Reset()
		if exec.input != "" {
			p.history.Add(exec.input)
		}
	case ControlC:
		p.renderer.BreakLine(p.buf)
		p.buf = NewBuffer()
		p.completion.Reset()
		p.history.Clear()
	case ControlD:
		shouldExit = true
	case Up:
		if !p.completion.Completing() {
			if newBuf, changed := p.history.Older(p.buf); changed {
				p.buf = newBuf
			}
			return
		}
		fallthrough
	case BackTab:
		p.completion.Previous()
	case Down:
		if !p.completion.Completing() {
			if newBuf, changed := p.history.Newer(p.buf); changed {
				p.buf = newBuf
			}
			return
		}
		fallthrough
	case Tab, ControlI:
		p.completion.Next()
	case Left:
		p.buf.CursorLeft(1)
	case Right:
		p.buf.CursorRight(1)
	case Backspace:
		if p.completion.Completing() {
			c := p.completer(p.buf.Text())[p.completion.selected]
			w := p.buf.Document().GetWordBeforeCursor()
			if w != "" {
				p.buf.DeleteBeforeCursor(len([]rune(w)))
			}
			p.buf.InsertText(c.Text, false, true)
			p.completion.Reset()
		}
		p.buf.DeleteBeforeCursor(1)
	case NotDefined:
		if p.completion.Completing() {
			c := p.completer(p.buf.Text())[p.completion.selected]
			w := p.buf.Document().GetWordBeforeCursor()
			if w != "" {
				p.buf.DeleteBeforeCursor(len([]rune(w)))
			}
			p.buf.InsertText(c.Text, false, true)
		}
		p.completion.Reset()
		p.buf.InsertText(string(b), false, true)
	default:
		p.completion.Reset()
	}
	return
}

func (p *Prompt) setUp() {
	p.in.Setup()
	p.renderer.Setup()
	p.renderer.UpdateWinSize(p.in.GetWinSize())
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

func handleSignals(in ConsoleParser, exitCh chan int, winSizeCh chan *WinSize) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGWINCH,
	)

	for {
		s := <-sigCh
		switch s {
		case syscall.SIGINT:  // kill -SIGINT XXXX or Ctrl+c
			log.Println("[SIGNAL] Catch SIGINT")
			exitCh <- 0

		case syscall.SIGTERM:  // kill -SIGTERM XXXX
			log.Println("[SIGNAL] Catch SIGTERM")
			exitCh <- 1

		case syscall.SIGQUIT:  // kill -SIGQUIT XXXX
			log.Println("[SIGNAL] Catch SIGQUIT")
			exitCh <- 0

		case syscall.SIGWINCH:
			log.Println("[SIGNAL] Catch SIGWINCH")
			winSizeCh <- in.GetWinSize()

		// TODO: SIGUSR1 -> Reopen log file.
		default:
		}
	}
}
