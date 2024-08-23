package io

import (
	"io"
	"sync"
)

// ReaderFromWriter is a type that wraps an io.ReaderFrom and implements io.Writer using io.Pipe.
// Due to the asynchronous nature of io.Pipe, Write() will only be guaranteed to be visible after
// a call to Sync() or Close().
type ReaderFromWriter struct {
	io.ReaderFrom
	*io.PipeWriter

	pipeR *io.PipeReader
	wg    sync.WaitGroup
}

// NewReaderFromWriter creates a new ReaderFromWriter.
func NewReaderFromWriter(readerFrom io.ReaderFrom) *ReaderFromWriter {
	pipeR, pipeW := io.Pipe()
	rfw := &ReaderFromWriter{
		ReaderFrom: readerFrom,
		pipeR:      pipeR,
		PipeWriter: pipeW,
	}
	rfw.wg.Add(1)
	go rfw.readFromPipe()
	return rfw
}

// readFromPipe reads data from the pipe and writes it to the readerFrom.
func (rfw *ReaderFromWriter) readFromPipe() {
	defer rfw.wg.Done()
	_, err := rfw.ReadFrom(rfw.pipeR)
	rfw.CloseWithError(err)
}

// Sync waits for the readFromPipe goroutine to finish.
func (rfw *ReaderFromWriter) Sync() {
	rfw.wg.Wait()
}

// Close closes the pipe writer and waits for the readFromPipe goroutine to finish.
func (rfw *ReaderFromWriter) Close() error {
	err := rfw.PipeWriter.Close()
	rfw.Sync()
	return err
}
