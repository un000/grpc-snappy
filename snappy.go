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
	cmpMu   sync.Mutex
	// cmpPool stores writers
	cmpPool sync.Pool

	dcmpMu   sync.Mutex
	// dcmpPool stores readers
	dcmpPool sync.Pool
)

type compressor struct {
}

func (c *compressor) Name() string {
	return Name
}

func (c *compressor) Compress(w io.Writer) (io.WriteCloser, error) {
	cmpMu.Lock()
	wr, inPool := cmpPool.Get().(*writeCloser)
	cmpMu.Unlock()
	if !inPool {
		return &writeCloser{Writer: snappy.NewBufferedWriter(w)}, nil
	} else {
		wr.Reset(w)
	}

	return wr, nil
}

func (c *compressor) Decompress(r io.Reader) (io.Reader, error) {
	dcmpMu.Lock()
	dr, inPool := dcmpPool.Get().(*reader)
	dcmpMu.Unlock()
	if !inPool {
		return &reader{Reader: snappy.NewReader(r)}, nil
	} else {
		dr.Reset(r)
	}

	return dr, nil
}

type writeCloser struct {
	*snappy.Writer
}

func (w *writeCloser) Close() error {
	defer func() {
		cmpMu.Lock()
		cmpPool.Put(w)
		cmpMu.Unlock()
	}()
	return w.Writer.Close()
}

type reader struct {
	*snappy.Reader
}

func (r *reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	if err == io.EOF {
		dcmpMu.Lock()
		dcmpPool.Put(r)
		dcmpMu.Unlock()
	}
	return n, err
}
