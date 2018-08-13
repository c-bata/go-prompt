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

func TestHistoryGetLines(t *testing.T) {
	h := NewHistory()
	h.Add("echo 1")
	h.Add("echo 2")
	h.Add("echo 3")
	expectedGetLines2 := []string{
		"echo 2",
		"echo 3",
	}
	returnGet := h.GetLines(2)
	if !reflect.DeepEqual(expectedGetLines2, returnGet) {
		t.Errorf("History.GetLines(2) returned %s, expected %s", returnGet, expectedGetLines2)
	}

	expectedGetLines3 := []string{
		"echo 1",
		"echo 2",
		"echo 3",
	}
	returnGet = h.GetLines(3)
	if !reflect.DeepEqual(expectedGetLines3, returnGet) {
		t.Errorf("History.GetLines(3) returned %s, expected %s", returnGet, expectedGetLines3)
	}

	// make sure requesting more lines than
	// there are history entries does not fail
	returnGet = h.GetLines(5)
	if !reflect.DeepEqual(expectedGetLines3, returnGet) {
		t.Errorf("History.GetLines(5) returned %s, expected %s", returnGet, expectedGetLines3)
	}

	expectedGetLines4 := []string{"history1"}
	h = &History{    
	    histories: []string{"history1", "history2", "history3", "history4"},
	    selected: 1,
	}
	returnGet = h.GetLines(4)
	if !reflect.DeepEqual(expectedGetLines4, returnGet) {
		t.Errorf("History.GetLines(4) returned %s, expected %s", returnGet, expectedGetLines4)
	}
}

func TestHistoryGetLast(t *testing.T) {
	h := NewHistory()
	h.Add("echo 1")
	h.Add("echo 2")
	h.Add("echo 3")
	expectedGetLast := "echo 3"
	returnGet := h.GetLast()
	if !reflect.DeepEqual(expectedGetLast, returnGet) {
		t.Errorf("History.GetLast() returned %s, expected %s", returnGet, expectedGetLast)
	}
}

func TestHistoryGetLine(t *testing.T) {
	h := NewHistory()
	h.Add("echo 1")
	h.Add("echo 2")
	h.Add("echo 3")
	expectedGetLine1 := "echo 2"
	returnGet := h.GetLine(1)
	if !reflect.DeepEqual(expectedGetLine1, returnGet) {
		t.Errorf("History.GetLine(1) returned %s, expected %s", returnGet, expectedGetLine1)
	}

	expectedGetLine0 := "echo 1"
	returnGet = h.GetLine(0)
	if !reflect.DeepEqual(expectedGetLine0, returnGet) {
		t.Errorf("History.GetLine(0) returned %s, expected %s", returnGet, expectedGetLine0)
	}
}
