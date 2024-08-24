package io_test

import (
	"io"
	"testing"

	gio "github.com/daotl/guts/io"
)

// ReaderFromWriteToReadWriteCloser tests the ReaderFromWriteToReadWriteCloser implementation.
func TestReaderFromWriteToReadWriteCloser(t *testing.T) {
	testReaderFromWriter(
		t,
		func(rf io.ReaderFrom) gio.ReadFromWriteCloser {
			return gio.NewReadFromWriteToReadWriterCloser(rf, &ExampleWriterTo{})
		},
	)

	testWriterToReader(
		t,
		func(wt io.WriterTo) gio.WriteToReadCloser {
			return gio.NewReadFromWriteToReadWriterCloser(&ExampleReaderFrom{}, wt)
		},
	)
}
