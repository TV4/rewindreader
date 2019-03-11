// Package rewindreader seamlessly buffers an io.Reader to allow rewinding.
package rewindreader

import (
	"bytes"
	"errors"
	"io"
)

// RewindReader buffers an io.Reader until the first call to its Rewind method,
// allowing the buffered data to be re-read one or more times. If the consumer
// reads past the size of the buffered data the buffer is deleted and the
// underlying data source is read from directly.
type RewindReader struct {
	r   io.Reader
	buf io.Reader
	src io.Reader
}

// New returns a new *RewindReader reading from the given io.Reader.
func New(r io.Reader) *RewindReader {
	var buf bytes.Buffer

	return &RewindReader{
		r:   io.TeeReader(r, &buf),
		buf: &buf,
		src: r,
	}
}

// Read implements io.Reader.
func (rwr *RewindReader) Read(p []byte) (n int, err error) {
	n, err = rwr.r.Read(p)
	if rwr.r == rwr.buf && n == 0 && err == io.EOF {
		rwr.r = rwr.src
		rwr.buf = nil
		return rwr.r.Read(p)
	}
	return n, err
}

// Rewind rewinds the stream to the start. If the stream has been consumed past
// the point at which a previous call to Rewind happened, an error will be
// returned.
func (rwr *RewindReader) Rewind() error {
	switch buf := rwr.buf.(type) {
	case *bytes.Buffer:
		rwr.buf = bytes.NewReader(buf.Bytes())
		rwr.r = rwr.buf
	case *bytes.Reader:
		buf.Seek(0, io.SeekStart)
	default:
		return errors.New("cannot rewind unbuffered stream")
	}

	return nil
}
