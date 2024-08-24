package io

import "io"

// ReadFromWriter is the interface that groups the basic ReadFrom and Write methods.
type ReadFromWriter = interface {
	io.ReaderFrom
	io.Writer
}

// ReadFromWriteCloser is the interface that groups the basic ReadFrom, Write and Close methods.
type ReadFromWriteCloser = interface {
	io.ReaderFrom
	io.WriteCloser
}

// WriteToReader is the interface that groups the basic WriteTo and Read methods.
type WriteToReader = interface {
	io.WriterTo
	io.Reader
}

// WriteToReadCloser is the interface that groups the basic WriteTo, Read and Close methods.
type WriteToReadCloser = interface {
	io.WriterTo
	io.ReadCloser
}

// ReadFromWriteToReadWriteCloser is the interface that groups the basic WriteTo, ReadFrom, Read, Write and Close methods.
type ReadFromWriteToReadWriteCloser = interface {
	io.ReaderFrom
	io.WriterTo
	io.ReadWriteCloser
}
