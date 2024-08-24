package io_test

import "io"

const TestStr = "Hello, World!"

var (
	TestStrChunks = []string{"Hello, ", "World", "!"}
	TestBin       = []byte(TestStr)
)

// ExampleReaderFrom is an example type that implements io.ReaderFrom.
type ExampleReaderFrom struct {
	data []byte
}

// ReadFrom reads data from the provided io.Reader.
func (erf *ExampleReaderFrom) ReadFrom(r io.Reader) (int64, error) {
	buf := make([]byte, 1024)
	var totalBytes int64
	for {
		n, err := r.Read(buf)
		if n > 0 {
			erf.data = append(erf.data, buf[:n]...)
			totalBytes += int64(n)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return totalBytes, err
		}
	}
	return totalBytes, nil
}

// ExampleWriterTo is an example type that implements io.WriterTo.
type ExampleWriterTo struct {
	data []byte
}

// WriteTo writes data to the provided io.Writer.
func (ewt *ExampleWriterTo) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(ewt.data)
	return int64(n), err
}
