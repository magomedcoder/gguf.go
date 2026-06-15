package format

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

// MappedReader - GGUF reader с memory-mapped файлом
type MappedReader struct {
	*Reader
	file *os.File
	data []byte
}

// OpenFileMapped открывает GGUF через mmap для zero-copy доступа к весам
func OpenFileMapped(path string) (*MappedReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, int(info.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("gguf: mmap: %w", err)
	}

	src := &mmapSource{data: data}
	r, err := Open(src)
	if err != nil {
		syscall.Munmap(data)
		f.Close()
		return nil, err
	}

	return &MappedReader{
		Reader: r,
		file:   f,
		data:   data,
	}, nil
}

// Close снимает mmap и закрывает файл
func (m *MappedReader) Close() error {
	if m.data != nil {
		_ = syscall.Munmap(m.data)
		m.data = nil
	}

	if m.file != nil {
		err := m.file.Close()
		m.file = nil
		return err
	}

	return nil
}

// Data возвращает mmap-срез всего файла
func (m *MappedReader) Data() []byte {
	return m.data
}

type mmapSource struct {
	data []byte
	pos  int64
}

func (m *mmapSource) Read(p []byte) (int, error) {
	if m.pos >= int64(len(m.data)) {
		return 0, io.EOF
	}

	n := copy(p, m.data[m.pos:])
	m.pos += int64(n)

	return n, nil
}

func (m *mmapSource) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case 0:
		abs = offset
	case 1:
		abs = m.pos + offset
	case 2:
		abs = int64(len(m.data)) + offset
	default:
		return 0, fmt.Errorf("invalid whence")
	}

	if abs < 0 || abs > int64(len(m.data)) {
		return 0, fmt.Errorf("seek out of range")
	}

	m.pos = abs

	return abs, nil
}

func (m *mmapSource) ReadAt(p []byte, off int64) (int, error) {
	if off < 0 || off > int64(len(m.data)) {
		return 0, fmt.Errorf("readat out of range")
	}

	n := copy(p, m.data[off:])
	if n < len(p) {
		return n, io.EOF
	}

	return n, nil
}

func (m *mmapSource) Slice(off, size int64) []byte {
	if off < 0 || off+size > int64(len(m.data)) {
		return nil
	}

	return m.data[off : off+size]
}
