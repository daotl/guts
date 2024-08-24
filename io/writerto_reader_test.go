package io_test

import (
	"bytes"
	"io"
	"testing"

	gio "github.com/daotl/guts/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testWriterToReader(
	t *testing.T,
	initWtr func(wt io.WriterTo) gio.WriteToReadCloser,
) {
	req := require.New(t)
	assr := assert.New(t)

	// Tests the WriterToReader implementation with single read.
	t.Run("Single read", func(t *testing.T) {
		// Create an instance of ExampleWriterTo with some data.
		ewt := &ExampleWriterTo{data: TestBin}
		// Create a WriterToReader wrapping the ExampleWriterTo.
		wtr := initWtr(ewt)

		// Read data from the WriterToReader and compare it to the expected output.
		buf := make([]byte, len(TestStr))
		n, err := wtr.Read(buf)
		req.NoError(err)
		assr.Equal(len(TestStr), n)
		assr.Equal(TestStr, string(buf))
	})

	// Tests the WriterToReader implementation with multiple reads.
	t.Run("Multiple reads", func(t *testing.T) {
		ewt := &ExampleWriterTo{data: TestBin}
		// Create a WriterToReader wrapping the ExampleWriterTo.
		wtr := initWtr(ewt)

		// Read data from the WriterToReader in chunks and compare it to the TestStr output.
		var result bytes.Buffer
		buf := make([]byte, 5)
		for {
			n, err := wtr.Read(buf)
			if err != nil && err != io.EOF {
				req.NoError(err)
			}
			if n == 0 {
				break
			}
			result.Write(buf[:n])
		}

		assr.Equal(TestStr, result.String())
	})

	// Tests the WriterToReader implementation with an empty WriterTo.
	t.Run("Empty WriterTo", func(t *testing.T) {
		// Create an instance of ExampleWriterTo with no data.
		ewt := &ExampleWriterTo{data: []byte{}}
		wtr := initWtr(ewt)

		// Read data from the WriterToReader and compare it to the expected output.
		buf := make([]byte, 5)
		n, err := wtr.Read(buf)
		req.NoError(err)
		req.Equal(0, n)
		n2, err2 := wtr.Read(buf)
		assr.Equal(io.EOF, err2)
		assr.Equal(0, n2)
	})
}

// TestWriterToReader tests the WriterToReader implementation.
func TestWriterToReader(t *testing.T) {
	testWriterToReader(
		t,
		func(wt io.WriterTo) gio.WriteToReadCloser { return gio.NewWriterToReader(wt) },
	)
}
