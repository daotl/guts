package io

import (
	"errors"
	"io"
)

// ReaderFromWriteToReadWriteCloser is a type that wraps an io.ReaderFrom and a io.WriterTo, and
// implements io.ReadWriteCloser, ReadFromWriteCloser and WriteToReadCloser using io.Pipe.
// Due to the asynchronous nature of io.Pipe, Write() will only be guaranteed to be visible after
// a call to Sync() or Close().
type ReaderFromWriteToReadWriteCloser struct {
	*ReaderFromWriter
	*WriterToReader
}

var _ ReadFromWriteToReadWriteCloser = (*ReaderFromWriteToReadWriteCloser)(nil)

// NewReadFromWriteToReadWriterCloser creates a new ReadFromWriteToReadWriterCloser.
func NewReadFromWriteToReadWriterCloser(
	readerFrom io.ReaderFrom,
	writerTo io.WriterTo,
) *ReaderFromWriteToReadWriteCloser {
	return &ReaderFromWriteToReadWriteCloser{
		ReaderFromWriter: NewReaderFromWriter(readerFrom),
		WriterToReader:   NewWriterToReader(writerTo),
	}
}

// Close closes the pipe writer and waits for the readFromPipe goroutine to finish.
func (rwc *ReaderFromWriteToReadWriteCloser) Close() error {
	err1 := rwc.ReaderFromWriter.Close()
	err2 := rwc.WriterToReader.Close()
	return errors.Join(err1, err2)
}
