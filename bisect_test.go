package prompt

import (
	"testing"
)

// Thanks!! https://play.golang.org/p/y9NRj_XVIW

func TestBisectRight(t *testing.T) {
	in := []int{1, 2, 3, 3, 3, 6, 7}

	r := BisectRight(in, 0)
	if r != 0 {
		t.Error("number 0 should inserted at 0 position, but got %d", r)
	}

	r = BisectRight(in, 4)
	if r != 5 {
		t.Error("number 4 should inserted at 5 position, but got %d", r)
	}
}
