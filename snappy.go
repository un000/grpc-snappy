package snappy

import (
	"io"
	"sync"

	"github.com/golang/snappy"
	"google.golang.org/grpc/encoding"
)

// Name is the name registered for the snappy compressor.
const Name = "snappy"

func init() {
	encoding.RegisterCompressor(&compressor{})
}

var (
	// cmpPool stores writers
	cmpPool sync.Pool
	// dcmpPool stores readers
	dcmpPool sync.Pool
)

type compressor struct {
}

func (c *compressor) Name() string {
	return Name
}

func (c *compressor) Compress(w io.Writer) (io.WriteCloser, error) {
	wr, inPool := cmpPool.Get().(*writeCloser)
	if !inPool {
		return &writeCloser{Writer: snappy.NewBufferedWriter(w)}, nil
	}
	wr.Reset(w)

	return wr, nil
}

func (c *compressor) Decompress(r io.Reader) (io.Reader, error) {
	dr, inPool := dcmpPool.Get().(*reader)
	if !inPool {
		return &reader{Reader: snappy.NewReader(r)}, nil
	}
	dr.Reset(r)

	return dr, nil
}

type writeCloser struct {
	*snappy.Writer
}

func (w *writeCloser) Close() error {
	defer func() {
		cmpPool.Put(w)
	}()

	return w.Writer.Close()
}

type reader struct {
	*snappy.Reader
}

func (r *reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	if err == io.EOF {
		dcmpPool.Put(r)
	}

	return n, err
}
