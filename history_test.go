package prompt

import (
	"reflect"
	"testing"
)

func TestHistoryClear(t *testing.T) {
	h := NewHistory()
	h.Add("foo")
	h.Clear()
	expected := &History{
		histories: []string{"foo"},
		tmp:       []string{"foo", ""},
		selected:  1,
	}
	if !reflect.DeepEqual(expected, h) {
		t.Errorf("Should be %#v, but got %#v", expected, h)
	}
}

func TestHistoryAdd(t *testing.T) {
	h := NewHistory()
	h.Add("echo 1")
	expected := &History{
		histories: []string{"echo 1"},
		tmp:       []string{"echo 1", ""},
		selected:  1,
	}
	if !reflect.DeepEqual(h, expected) {
		t.Errorf("Should be %v, but got %v", expected, h)
	}
}

func TestHistoryOlder(t *testing.T) {
	h := NewHistory()
	h.Add("echo 1")

	// Prepare buffer
	buf := NewBuffer()
	buf.InsertText("echo 2", false, true)

	// [1 time] Call Older function
	buf1, changed := h.Older(buf)
	if !changed {
		t.Error("Should be changed history but not changed.")
	}
	if buf1.Text() != "echo 1" {
		t.Errorf("Should be %#v, but got %#v", "echo 1", buf1.Text())
	}

	// [2 times] Call Older function
	buf = NewBuffer()
	buf.InsertText("echo 1", false, true)
	buf2, changed := h.Older(buf)
	if changed {
		t.Error("Should be not changed history but changed.")
	}
	if !reflect.DeepEqual("echo 1", buf2.Text()) {
		t.Errorf("Should be %#v, but got %#v", "echo 1", buf2.Text())
	}
}

func TestHistoryGet(t *testing.T) {
	h := NewHistory()
	h.Add("echo 1")
	h.Add("echo 2")
	h.Add("echo 3")
	expectedGet2 := []string{
		"echo 2",
		"echo 3",
	}
	returnGet := h.Get(2)
	if !reflect.DeepEqual(expectedGet2, returnGet) {
		t.Errorf("History.Get(2) return %s, expected %s", returnGet, expectedGet2)
	}

	expectedGet3 := []string{
		"echo 1",
		"echo 2",
		"echo 3",
	}
	returnGet = h.Get(3)
	if !reflect.DeepEqual(expectedGet3, returnGet) {
		t.Errorf("History.Get(3) returned %s, expected %s", returnGet, expectedGet3)
	}

	// make sure requesting more lines than
	// there are history entries does not fail
	returnGet = h.Get(5)
	if !reflect.DeepEqual(expectedGet3, returnGet) {
		t.Errorf("History.Get(5) returned %s, expected %s", returnGet, expectedGet3)
	}
	
}
