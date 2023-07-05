package prompt

import (
	"io"
	"strings"

	"github.com/mattn/go-runewidth"
)

// Position stores the coordinates
// of a p.
//
// (0, 0) represents the top-left corner of the prompt,
// while (n, n) the bottom-right corner.
type Position struct {
	X, Y int
}

// Join two positions and return a new position.
func (p Position) Join(other Position) Position {
	if other.Y == 0 {
		p.X += other.X
	} else {
		p.X = other.X
		p.Y += other.Y
	}
	return p
}

// Add two positions and return a new position.
func (p Position) Add(other Position) Position {
	return Position{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

// Subtract two positions and return a new position.
func (p Position) Subtract(other Position) Position {
	return Position{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
}

// positionAtEndOfString calculates the position of the
// p at the end of the given string.
func positionAtEndOfString(str string, columns int) Position {
	// fmt.Printf("%q\n", str)
	pos := positionAtEndOfReader(strings.NewReader(str), columns)
	return pos
}

// positionAtEndOfReader calculates the position of the
// p at the end of the given io.Reader.
func positionAtEndOfReader(reader io.RuneReader, columns int) Position {
	var down, right int

charLoop:
	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			break charLoop
		}

		switch char {
		case '\r':
			char, _, err := reader.ReadRune()
			if err != nil {
				break charLoop
			}

			if char == '\n' {
				down++
				right = 0
			}
		case '\n':
			down++
			right = 0
		default:
			right += runewidth.RuneWidth(char)
			if right == columns {
				right = 0
				down++
			}
		}
	}

	return Position{
		X: right,
		Y: down,
	}
}
