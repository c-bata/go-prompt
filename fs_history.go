package prompt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FSHistory stores entered text in memory and in the file at FilePath.
type FSHistory struct {
	History
	// FilePath is the history file.
	FilePath string
}

var _ HistoryInterface = &FSHistory{}

const fsHistoryFileFormatVersion = 1

// NewFSHistory loads history from the file at path.
// A missing file starts an empty history.
func NewFSHistory(path string) (*FSHistory, error) {
	if path == "" {
		return nil, fmt.Errorf("history path is empty")
	}
	path = filepath.Clean(path)
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return nil, fmt.Errorf("history path %q is a directory", path)
	} else if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	h := &FSHistory{
		History: History{
			histories: []string{},
			tmp:       []string{""},
			selected:  0,
			size:      defaultHistorySize,
		},
		FilePath: path,
	}
	if err := h.load(); err != nil {
		return nil, err
	}
	return h, nil
}

// Add stores input in memory and rewrites the history file.
func (h *FSHistory) Add(input string) {
	atMaxSize := h.AtMaxSize()

	h.History.Add(input)

	if atMaxSize {
		err := h.save()
		if err != nil {
			panic(fmt.Sprintf("error on prompt history Add: save: %s", err))
		}
		return
	}

	serialized, err := h.serializeEntry(input)
	if err != nil {
		panic(fmt.Sprintf("error on prompt history Add: serialization: %s", err))
	}
	if err := h.append(serialized); err != nil {
		panic(fmt.Sprintf("error on prompt history Add: file append: %s", err))
	}
}

// DeleteAll removes every entry from memory and from the history file.
func (h *FSHistory) DeleteAll() {
	h.History.DeleteAll()
	err := h.save()
	if err != nil {
		panic(fmt.Sprintf("error on prompt history DeleteAll: %s", err))
	}
}

const fsHistoryLengthBytes = 4
const fsHistoryHeaderBytes = fsHistoryLengthBytes + 1

func (h *FSHistory) load() error {
	file, err := os.OpenFile(h.FilePath, os.O_RDONLY|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() == 0 {
		h.save()
		h.Clear()
		return nil
	}

	formatVersionBytes := make([]byte, 1)
	_, err = file.Read(formatVersionBytes)
	if err != nil {
		return err
	}

	if formatVersionBytes[0] != fsHistoryFileFormatVersion {
		h.save()
		h.Clear()
		return nil
	}

	// skip the first two bytes with the file format version
	file.Seek(2, 0)

	entries := make([]string, 0)
	for {
		entry, err := h.readEntry(file)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		entries = append(entries, entry)
	}

	h.histories = entries
	h.Clear()
	return nil
}

func (h *FSHistory) readEntry(file *os.File) (string, error) {
	// move the cursor 1 byte forward for the newline before the length
	file.Seek(1, 1)

	entryLengthBytes := make([]byte, fsHistoryLengthBytes)
	_, err := file.Read(entryLengthBytes)
	if err != nil {
		return "", err
	}

	// move the cursor 1 byte forward for the newline after the length
	file.Seek(1, 1)

	entryLength := int(binary.BigEndian.Uint32(entryLengthBytes))
	entryBodyBytes := make([]byte, entryLength)
	_, err = file.Read(entryBodyBytes)
	if err != nil {
		return "", err
	}

	return string(entryBodyBytes), nil
}

func (h *FSHistory) serializeEntry(entry string) ([]byte, error) {
	var buf bytes.Buffer
	err := h.serializeEntryToBuffer(&buf, entry)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *FSHistory) serializeEntryToBuffer(buf *bytes.Buffer, entry string) error {
	buf.WriteByte('\n')
	lengthBytes := make([]byte, fsHistoryLengthBytes)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(entry)))
	_, err := buf.Write(lengthBytes)
	if err != nil {
		return err
	}
	buf.WriteByte('\n')
	buf.WriteString(entry)
	return nil
}

func (h *FSHistory) save() error {
	entries := h.histories
	if entries == nil {
		entries = []string{}
	}

	var buf bytes.Buffer

	buf.WriteByte(fsHistoryFileFormatVersion)
	buf.WriteByte('\n')

	for _, entry := range entries {
		h.serializeEntryToBuffer(&buf, entry)
	}

	return h.write(buf.Bytes())
}

func (h *FSHistory) write(data []byte) error {
	return os.WriteFile(h.FilePath, data, 0o600)
}

func (h *FSHistory) append(data []byte) error {
	file, err := os.OpenFile(h.FilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	_, writeErr := file.Write(data)
	closeErr := file.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}
