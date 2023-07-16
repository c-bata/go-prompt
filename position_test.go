//go:build !windows
// +build !windows

package prompt

import (
	"testing"

	istrings "github.com/elk-language/go-prompt/strings"
	"github.com/google/go-cmp/cmp"
)

func TestPositionAtEndOfString(t *testing.T) {
	tests := map[string]struct {
		input   string
		columns istrings.Width
		want    Position
	}{
		"empty": {
			input:   "",
			columns: 20,
			want: Position{
				X: 0,
				Y: 0,
			},
		},
		"one letter": {
			input:   "f",
			columns: 20,
			want: Position{
				X: 1,
				Y: 0,
			},
		},
		"one word": {
			input:   "foo",
			columns: 20,
			want: Position{
				X: 3,
				Y: 0,
			},
		},
		"one-line fits in columns": {
			input:   "foo bar",
			columns: 20,
			want: Position{
				X: 7,
				Y: 0,
			},
		},
		"multiline": {
			input:   "foo\nbar\n",
			columns: 20,
			want: Position{
				X: 0,
				Y: 2,
			},
		},
		"one-line wrapping": {
			input:   "foobar",
			columns: 3,
			want: Position{
				X: 0,
				Y: 2,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := positionAtEndOfString(tc.input, tc.columns)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestPositionAdd(t *testing.T) {
	tests := map[string]struct {
		left  Position
		right Position
		want  Position
	}{
		"empty": {
			left:  Position{},
			right: Position{},
			want:  Position{},
		},
		"only X": {
			left:  Position{X: 1},
			right: Position{X: 2},
			want:  Position{X: 3},
		},
		"only Y": {
			left:  Position{Y: 1},
			right: Position{Y: 2},
			want:  Position{Y: 3},
		},
		"different coordinates": {
			left:  Position{X: 1},
			right: Position{Y: 2},
			want:  Position{X: 1, Y: 2},
		},
		"both X and Y": {
			left:  Position{X: 1, Y: 5},
			right: Position{X: 10, Y: 2},
			want:  Position{X: 11, Y: 7},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.left.Add(tc.right)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestPositionSubtract(t *testing.T) {
	tests := map[string]struct {
		left  Position
		right Position
		want  Position
	}{
		"empty": {
			left:  Position{},
			right: Position{},
			want:  Position{},
		},
		"only X": {
			left:  Position{X: 1},
			right: Position{X: 2},
			want:  Position{X: -1},
		},
		"only Y": {
			left:  Position{Y: 5},
			right: Position{Y: 2},
			want:  Position{Y: 3},
		},
		"different coordinates": {
			left:  Position{X: 1},
			right: Position{Y: 2},
			want:  Position{X: 1, Y: -2},
		},
		"both X and Y": {
			left:  Position{X: 1, Y: 5},
			right: Position{X: 10, Y: 2},
			want:  Position{X: -9, Y: 3},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.left.Subtract(tc.right)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestPositionJoin(t *testing.T) {
	tests := map[string]struct {
		left  Position
		right Position
		want  Position
	}{
		"empty": {
			left:  Position{},
			right: Position{},
			want:  Position{},
		},
		"only X": {
			left:  Position{X: 1},
			right: Position{X: 2},
			want:  Position{X: 3},
		},
		"only Y": {
			left:  Position{Y: 1},
			right: Position{Y: 2},
			want:  Position{Y: 3},
		},
		"different coordinates": {
			left:  Position{X: 5},
			right: Position{Y: 2},
			want:  Position{X: 0, Y: 2},
		},
		"both X and Y": {
			left:  Position{X: 1, Y: 5},
			right: Position{X: 10, Y: 2},
			want:  Position{X: 10, Y: 7},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.left.Join(tc.right)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
