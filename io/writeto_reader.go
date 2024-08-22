package io

import (
	"io"
)

// WriteToReader is a type that wraps an io.WriterTo implementation and implements io.Reader using io.Pipe.
type WriteToReader struct {
	io.WriterTo
	*io.PipeReader
	pipeW *io.PipeWriter
}

// NewWriteToReader creates a new WriteToReader.
func NewWriteToReader(writerTo io.WriterTo) *WriteToReader {
	pipeR, pipeW := io.Pipe()
	wtr := &WriteToReader{
		WriterTo:   writerTo,
		PipeReader: pipeR,
		pipeW:      pipeW,
	}
	go wtr.writeToPipe()
	return wtr
}

// writeToPipe writes data from the writerTo to the pipe.
func (wtr *WriteToReader) writeToPipe() {
	_, err := wtr.WriteTo(wtr.pipeW)
	wtr.pipeW.CloseWithError(err)
}
