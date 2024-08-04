package field_test

import "io"

type MockFile struct {
	content []byte
	offset  int64
}

func NewMockFile(content []byte) *MockFile {
	return &MockFile{content: content, offset: 0}
}

// Read reads up to len(p) bytes into p from the file
func (file *MockFile) Read(p []byte) (int, error) {
	if file.offset >= int64(len(file.content)) {
		return 0, io.EOF
	}
	n := copy(p, file.content[file.offset:])
	file.offset += int64(n)
	return n, nil
}

// ReadAt reads len(p) bytes into p from the file at offset off
func (file *MockFile) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(len(file.content)) {
		return 0, io.EOF
	}
	n := copy(p, file.content[off:])
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

// WriteAt writes len(p) bytes from p to the file at offset off
func (m *MockFile) WriteAt(p []byte, off int64) (int, error) {
	end := off + int64(len(p))
	if end > int64(len(m.content)) {
		newContent := make([]byte, end)
		copy(newContent, m.content)
		m.content = newContent
	}
	n := copy(m.content[off:], p)
	return n, nil
}
