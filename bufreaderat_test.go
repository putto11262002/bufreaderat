package bufreaderat

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
)

type ReadCounter struct {
	r       io.ReaderAt
	counter int
}

func (r *ReadCounter) ReadAt(p []byte, offset int64) (int, error) {
	r.counter++
	return r.r.ReadAt(p, offset)
}

func (r *ReadCounter) Count() int {
	return r.counter
}

func NewReadCounter(r io.ReaderAt) *ReadCounter {
	return &ReadCounter{
		r: r,
	}
}

func Setup(bufferSize int, n int) ([]byte, *ReadCounter, *BufReaderAt) {
	data := make([]byte, n)
	rand.Read(data)
	underlyingReader := bytes.NewReader(data)
	readCounter := NewReadCounter(underlyingReader)
	bufReader := New(readCounter, bufferSize)
	return data, readCounter, bufReader
}

func TestReadAt(t *testing.T) {
	// case: hit buffer
	{
		data, readCounter, bufReader := Setup(20, 20)
		p := make([]byte, 10)
		n, err := bufReader.ReadAt(p, 0)
		if err != nil {
			t.Fatal(err)
		}
		if n != len(p) {
			t.Fatalf("expected: %d, got: %d", len(p), n)
		}
		if !bytes.Equal(p, data[:10]) {
			t.Fatalf("expected: %s, got: %s", data[:2], p)
		}
		// Second read should hit buffer
		if n, err = bufReader.ReadAt(p, 10); err != nil {
			t.Fatal(err)
		}
		if n != len(p) {
			t.Fatalf("expected: %d, got: %d", len(p), n)
		}
		if !bytes.Equal(p, data[10:20]) {
			t.Fatalf("expected: %s, got: %s", data[2:4], p)
		}
		if readCounter.Count() != 1 {
			t.Fatalf("expected: %d, got: %d", 1, readCounter.Count())
		}
	}

	// case: underlying reader returns EOF but buffer is not exhausted
	{

		data, readCounter, bufReader := Setup(20, 20)
		p := make([]byte, 5)
		n, err := bufReader.ReadAt(p, 10)
		if err != nil {
			t.Fatal(err)
		}
		if n != len(p) {
			t.Fatalf("expected: %d, got: %d", len(p), n)
		}
		if !bytes.Equal(p, data[10:15]) {
			t.Fatalf("expected: %s, got: %s", data[10:15], p)
		}
		if readCounter.Count() != 1 {
			t.Fatalf("expected: %d, got: %d", 1, readCounter.Count())
		}
	}

	// case: underlying reader returns EOF and requested data is larger than the buffer
	{
		data, readCounter, bufReader := Setup(10, 20)
		p := make([]byte, 10)
		n, err := bufReader.ReadAt(p, 15)
		if err != io.EOF {
			t.Fatalf("expected: %v, got: %v", io.EOF, err)
		}
		if n != 5 {
			t.Fatalf("expected: %d, got: %d", 5, n)
		}
		if !bytes.Equal(data[15:], p[:5]) {
			t.Fatalf("expected: %s, got: %s", data[15:], p[:5])
		}
		if readCounter.Count() != 1 {
			t.Fatalf("expected: %d, got: %d", 1, readCounter.Count())
		}
	}

	// case: read more than buffer size
	{
		data, readCounter, bufReader := Setup(20, 30)
		p := make([]byte, 30)
		n, err := bufReader.ReadAt(p, 0)
		if err != nil {
			t.Fatal(err)
		}
		if n != len(p) {
			t.Fatalf("expected: %d, got: %d", len(p), n)
		}
		if !bytes.Equal(p, data) {
			t.Fatalf("expected: %s, got: %s", data, p)
		}

		if readCounter.Count() != 1 {
			t.Fatalf("expected: %d, got: %d", 1, readCounter.Count())
		}

	}
	// case: read more than buffer size and more than underlying data
	{
		data, readCounter, bufReader := Setup(5, 10)
		p := make([]byte, 10)
		n, err := bufReader.ReadAt(p, 5)
		if err != io.EOF {
			t.Fatal(err)
		}
		if n != 5 {
			t.Fatalf("expected: %d, got %d", 8, n)
		}
		if !bytes.Equal(data[5:], p[:5]) {
			t.Fatalf("expected: %s, got: %s", data[5:], p[:5])
		}
		if readCounter.Count() != 1 {
			t.Fatalf("expected: %d, got: %d", 1, readCounter.Count())
		}
	}

}
