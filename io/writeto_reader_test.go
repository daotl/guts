package io

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ExampleWriterTo is an example type that implements io.WriterTo.
type ExampleWriterTo struct {
	data []byte
}

// WriteTo writes data to the provided io.Writer.
func (ewt *ExampleWriterTo) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(ewt.data)
	return int64(n), err
}

// TestWriteToReader tests the WriteToReader implementation.
func TestWriteToReader(t *testing.T) {
	req := require.New(t)
	assr := assert.New(t)

	// Tests the WriteToReader implementation with single read.
	t.Run("Single read", func(t *testing.T) {
		// Create an instance of ExampleWriterTo with some data.
		ewt := &ExampleWriterTo{data: TestBin}

		// Create a WriteToReader wrapping the ExampleWriterTo.
		wtr := NewWriteToReader(ewt)

		// Read data from the WriteToReader and compare it to the expected output.
		buf := make([]byte, len(TestStr))
		n, err := wtr.Read(buf)
		req.NoError(err)
		assr.Equal(len(TestStr), n)
		assr.Equal(TestStr, string(buf))
	})

	// Tests the WriteToReader implementation with multiple reads.
	t.Run("Multiple reads", func(t *testing.T) {
		// Create an instance of ExampleWriterTo with some data.
		ewt := &ExampleWriterTo{data: TestBin}

		// Create a WriteToReader wrapping the ExampleWriterTo.
		wtr := NewWriteToReader(ewt)

		// Read data from the WriteToReader in chunks and compare it to the TestStr output.
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

	// Tests the WriteToReader implementation with an empty writer.
	t.Run("Empty writer", func(t *testing.T) {
		// Create an instance of ExampleWriterTo with no data.
		ewt := &ExampleWriterTo{data: []byte{}}

		// Create a WriteToReader wrapping the ExampleWriterTo.
		wtr := NewWriteToReader(ewt)

		// Read data from the WriteToReader and compare it to the expected output.
		buf := make([]byte, 5)
		n, err := wtr.Read(buf)
		req.NoError(err)
		req.Equal(0, n)
		n2, err2 := wtr.Read(buf)
		assr.Equal(io.EOF, err2)
		assr.Equal(0, n2)
	})
}
