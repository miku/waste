package waste

import "io"

// Counter reads the number of bytes read from an underlying reader. Not goroutine
// save.
type Counter struct {
	r io.Reader
	n int
}

func (r *Counter) Read(p []byte) (n int, err error) {
	if r.r == nil {
		return 0, io.EOF
	}
	n, err = r.r.Read(p)
	r.n += n
	return
}

// N returns the number of bytes read.
func (r *Counter) N() int {
	return r.n
}
