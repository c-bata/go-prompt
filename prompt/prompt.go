package prompt

import (
	"syscall"
	"os"
	"os/signal"
	"time"
)

type Executor func(*Buffer) string
type Completer func(*Buffer) []string

type Prompt struct {
	in        *VT100Parser
	out       *VT100Writer
	buf       *Buffer
	renderer  *Render
	title     string
	executor  Executor
	completer Completer
}

func (p *Prompt) Run() {
	p.setUp()
	defer p.tearDown()


	bufCh := make(chan []byte, 128)
	go readBuffer(bufCh)

	exitCh := make(chan bool, 16)
	winSizeCh := make(chan *WinSize, 128)
	go handleSignals(p.in, exitCh, winSizeCh)

	for {
		select {
		case b := <- bufCh:
			ac := p.in.GetASCIICode(b)
			if ac == nil {
				p.buf.InsertText(string(b), false, true)
				p.out.EraseDown()
				p.out.WriteRaw(b)
				after := p.buf.Document().TextAfterCursor()
				p.out.WriteStr(after)
			} else if ac.Key == ControlJ || ac.Key == Enter {
				p.out.EraseDown()
				p.out.WriteStr(p.buf.Document().TextAfterCursor())
				res := p.executor(p.buf)
				p.out.WriteStr(res)
				p.buf = NewBuffer()
			} else if ac.Key == ControlC {
				return
			} else {
				InputHandler(ac, p.buf, p.out)
			}

			// Display completions
			if w := p.buf.Document().GetWordBeforeCursor(); w != "" {
				completions := p.completer(p.buf)
				p.renderer.RenderCompletion(completions)
			}

			completions := []string{"select", "insert", "update", "where"}
			p.renderer.Render(p.buf, completions)
			p.out.Flush()
		case w := <- winSizeCh:
			p.renderer.UpdateWinSize(w)
		case e := <- exitCh:
			if e {
				return
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (p *Prompt) setUp() {
	p.in.Setup()
	if p.title != "" {
		p.out.SetTitle(p.title)
	}
}

func (p *Prompt) tearDown() {
	p.in.TearDown()
	p.out.ClearTitle()
	p.out.EraseDown()
	p.out.Flush()
}

func readBuffer(bufCh chan []byte) {
	buf := make([]byte, 1024)

	for {
		if n, err := syscall.Read(syscall.Stdin, buf); err == nil {
			bufCh <- buf[:n]
		}
	}
}

func handleSignals(in *VT100Parser, exitCh chan bool, winSizeCh chan *WinSize) {
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
			exitCh <- true

			// kill -SIGINT XXXX or Ctrl+c
		case syscall.SIGINT:
			exitCh <- true

			// kill -SIGTERM XXXX
		case syscall.SIGTERM:
			exitCh <- true

			// kill -SIGQUIT XXXX
		case syscall.SIGQUIT:
			exitCh <- true

		case syscall.SIGWINCH:
			winSizeCh <- in.GetWinSize()
		default:
		}
	}
}

func NewPrompt(executor Executor, completer Completer) *Prompt {
	out := NewVT100Writer()
	return &Prompt{
		in: NewVT100Parser(),
		out: out,
		renderer: &Render{
			Prefix: ">>> ",
			out:    out,
		},
		title: "Hello! this is prompt toolkit",
		buf: NewBuffer(),
		executor: executor,
		completer: completer,
	}
}
