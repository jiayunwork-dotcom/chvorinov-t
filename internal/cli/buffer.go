package cli

import "os"

// castingBuffer wraps the JSON input file so LoadInput can decode and
// then release the handle. Close is only safe once; a second Close
// panics so callers do not silently double-close.
type castingBuffer struct {
	f      *os.File
	closed bool
}

func openCastingBuffer(path string) (*castingBuffer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &castingBuffer{f: f}, nil
}

func (b *castingBuffer) File() *os.File {
	return b.f
}

func (b *castingBuffer) Close() error {
	if b.closed {
		return nil
	}
	b.closed = true
	return b.f.Close()
}

func (b *castingBuffer) Release() error {
	return b.Close()
}
