// Package bufreaderat wraps io.ReaderAt, by creating a wrapper object
// that also implement io.Reader but provide buffering.
package bufreaderat

import (
	"io"
)

const (
	DEFAULT_BUFFER_SIZE = 1024
)

type BufReaderAt struct {
	readerAt io.ReaderAt
	buf      []byte
	offset   int64
	len      int64
	err      error
}

func Default(r io.ReaderAt) *BufReaderAt {
	return &BufReaderAt{
		readerAt: r,
		buf:      make([]byte, DEFAULT_BUFFER_SIZE),
		offset:   0,
		len:      0,
	}
}

func New(r io.ReaderAt, size int) *BufReaderAt {
	return &BufReaderAt{
		readerAt: r,
		buf:      make([]byte, size),
		offset:   0,
		len:      0,
	}
}

func (r *BufReaderAt) bufEnd() int64 {
	return r.offset + r.len
}


// bufOffset returns the offset relative to r.offset (offset of the buffer from the start of the file) 
func (r *BufReaderAt) bufOffset(offset int64) int64 {
	return offset - r.offset
}

func (r *BufReaderAt) bufCap() int64 {
	return int64(cap(r.buf))

}


func (r *BufReaderAt) ReadAt(offset int64, p []byte) (n int, er error) {
	pn := int64(len(p))
	// read from buffer
	if offset >= r.offset && pn <= r.len {
		n = copy(p, r.buf[r.bufOffset(offset):r.len])
		return n, r.err
	}

	// read directly into p
	if pn > r.bufCap() {
		n, r.err = r.readerAt.ReadAt(p, offset)
		return n, r.err
	}

	n, r.err = r.readerAt.ReadAt(r.buf, offset)
		var read int64
	if n > 0 {
		r.offset = offset
		r.len = int64(n)
		
		if pn > r.len{
			read = r.len - r.bufOffset(offset) 
		} else {
			read = pn
		}
		
		copy(p, r.buf[r.bufOffset(offset):r.len])
	}

	if r.err == io.EOF && offset+pn < r.bufEnd() {
		r.err = nil
	}

	return int(read), r.err
}
