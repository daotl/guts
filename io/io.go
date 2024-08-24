package io

import "io"

// ReadFromWriter is the interface that groups the basic ReadFrom and Write methods.
type ReadFromWriter = interface {
	io.ReaderFrom
	io.Writer
}

// ReadFromWriter is the interface that groups the basic ReadFrom, Write and Close methods.
type ReadFromWriteCloser = interface {
	io.ReaderFrom
	io.Writer
	io.Closer
}

// WriteToReader is the interface that groups the basic WriteTo and Read methods.
type WriteToReader = interface {
	io.WriterTo
	io.Reader
}

// WriteToReader is the interface that groups the basic WriteTo, Read and Close methods.
type WriteToReadCloser = interface {
	io.WriterTo
	io.Reader
	io.Closer
}
