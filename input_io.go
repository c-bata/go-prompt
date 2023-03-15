package prompt

import (
	"bytes"
	"io"

	"github.com/mxk/go-flowrate/flowrate"
)

type ioParser struct {
	rows, cols uint16
	r          *flowrate.Reader
}

// Setup should be called before starting input
func (p *ioParser) Setup() error {
	p.r.SetBlocking(true)
	return nil
}

// TearDown should be called after stopping input
func (p *ioParser) TearDown() error {
	p.r.SetBlocking(false)
	return nil
}

// GetKey returns Key correspond to input byte codes.
func (p *ioParser) GetKey(b []byte) Key {
	for _, k := range ASCIISequences {
		if bytes.Equal(k.ASCIICode, b) {
			return k.Key
		}
	}
	return NotDefined
}

// GetWinSize returns WinSize object to represent width and height of terminal.
func (p *ioParser) GetWinSize() *WinSize {
	return &WinSize{Col: p.cols, Row: p.rows}
}

// Read returns byte array.
func (p *ioParser) Read() (b []byte, err error) {
	b = make([]byte, 1024)
	n, err := p.r.Read(b)
	if err != nil {
		return nil, err
	}
	return b[0:n], nil

}

var _ ConsoleParser = (*ioParser)(nil)

// NewIOParser returns a console parser backed by io.Reader
func NewIOParser(rows, cols uint16, r io.Reader) ConsoleParser {
	frr := flowrate.NewReader(r, 1<<20)
	frr.SetBlocking(true)

	return &ioParser{rows: rows, cols: cols, r: frr}
}
