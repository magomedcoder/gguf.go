package format

import "io"

// TensorInfo представляет тензор в файле GGUF
type TensorInfo struct {
	reader *Reader

	Name       string
	Dimensions []uint64
	Type       GGML
	Offset     uint64
}

// RawView возвращает zero-copy срез данных тензора, если источник - mmap или []byte
func (t *TensorInfo) RawView() ([]byte, bool) {
	start := t.reader.tensorOffset + int64(t.Offset)
	size := t.Size()

	if dr, ok := t.reader.r.(dataReader); ok {
		return dr.Slice(start, size), true
	}

	return nil, false
}

type dataReader interface {
	Slice(off, size int64) []byte
}

// Reader возвращает io.Reader для чтения данных тензора
// Читатель ограничен размером данных тензора и не меняет позицию исходного файла
func (t *TensorInfo) Reader() (io.Reader, error) {
	start := t.reader.tensorOffset + int64(t.Offset)
	size := t.Size()

	ra, ok := t.reader.r.(io.ReaderAt)
	if !ok {
		return nil, errReaderAtRequired
	}

	return io.NewSectionReader(ra, start, size), nil
}

// Size возвращает размер данных тензора в байтах
func (t *TensorInfo) Size() int64 {
	return t.Type.dataSize(t.Dimensions)
}

// ValuesCount возвращает число элементов тензора
func (t *TensorInfo) ValuesCount() int64 {
	n := uint64(1)
	for _, d := range t.Dimensions {
		n *= d
	}

	return int64(n)
}
