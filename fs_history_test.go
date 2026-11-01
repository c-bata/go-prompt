package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewFSHistoryPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history")
	h, err := NewFSHistory(path)
	if err != nil {
		t.Fatal(err)
	}
	if h.FilePath != path {
		t.Fatalf("path %#v, want %#v", h.FilePath, path)
	}

	if _, err := NewFSHistory(""); err == nil {
		t.Fatal("empty path should be rejected")
	}
	if _, err := NewFSHistory(t.TempDir()); err == nil {
		t.Fatal("directory path should be rejected")
	}
}

func TestFSHistoryPersistsEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".repl-history")
	h, err := NewFSHistory(path)
	if err != nil {
		t.Fatal(err)
	}

	h.Add("echo 1")
	h.Add("line1\nline2")
	h.Add(`say "hi"`)

	reloaded, err := NewFSHistory(path)
	if err != nil {
		t.Fatal(err)
	}
	got := reloaded.Entries()
	want := []string{"echo 1", "line1\nline2", `say "hi"`}
	if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("entries %#v, want %#v", got, want)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.IsDir() {
		t.Fatal("history path is a directory")
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("file mode %o, want 600", info.Mode().Perm())
	}

	buf := NewBuffer()
	buf.InsertTextMoveCursor("draft", DefColCount, DefRowCount, false)
	older, changed := reloaded.Older(buf, DefColCount, DefRowCount)
	if !changed || older.Text() != `say "hi"` {
		t.Fatalf("older text %#v, changed %v", older.Text(), changed)
	}

	reloaded.DeleteAll()
	reloaded, err = NewFSHistory(path)
	if err != nil {
		t.Fatal(err)
	}
	if entries := reloaded.Entries(); len(entries) != 0 {
		t.Fatalf("entries after delete %#v", entries)
	}
}

func TestFSHistoryKeepsNewestEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history")
	h, err := NewFSHistory(path)
	if err != nil {
		t.Fatal(err)
	}
	h.size = 2
	h.Add("one")
	h.Add("two")
	h.Add("three")

	reloaded, err := NewFSHistory(path)
	if err != nil {
		t.Fatal(err)
	}
	reloaded.size = 2
	got := reloaded.Entries()
	want := []string{"two", "three"}
	if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("entries %#v, want %#v", got, want)
	}
}
