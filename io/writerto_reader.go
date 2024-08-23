package io

import (
	"io"
)

// WriterToReader is a type that wraps an io.WriterTo implementation and implements io.Reader using io.Pipe.
type WriterToReader struct {
	io.WriterTo
	*io.PipeReader
	pipeW *io.PipeWriter
}

// NewWriterToReader creates a new WriterToReader.
func NewWriterToReader(writerTo io.WriterTo) *WriterToReader {
	pipeR, pipeW := io.Pipe()
	wtr := &WriterToReader{
		WriterTo:   writerTo,
		PipeReader: pipeR,
		pipeW:      pipeW,
	}
	go wtr.writeToPipe()
	return wtr
}

// writeToPipe writes data from the writerTo to the pipe.
func (wtr *WriterToReader) writeToPipe() {
	_, err := wtr.WriteTo(wtr.pipeW)
	wtr.pipeW.CloseWithError(err)
}
