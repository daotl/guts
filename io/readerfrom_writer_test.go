package io

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// TestReaderFromWriter tests the ReaderFromWriter implementation.
func TestReaderFromWriter(t *testing.T) {
	req := require.New(t)
	assr := assert.New(t)

	// Tests the ReaderFromWriter implementation with single write.
	t.Run("Single write", func(t *testing.T) {
		// Create an instance of ExampleReaderFrom.
		erf := &ExampleReaderFrom{}
		rfw := NewReaderFromWriter(erf)

		// Write data to the ReaderFromWriter.
		n, err := rfw.Write(TestBin)
		req.NoError(err)
		assr.Equal(len(TestBin), n)

		// Close the pipe writer to signal EOF.
		req.NoError(rfw.Close())

		// Check the data read by ExampleReaderFrom.
		assr.Equal(TestStr, string(erf.data))
	})

	// Tests the ReaderFromWriter implementation with multiple writes.
	t.Run("Multiple writes", func(t *testing.T) {
		// Create an instance of ExampleReaderFrom.
		erf := &ExampleReaderFrom{}
		rfw := NewReaderFromWriter(erf)

		// Write data to the ReaderFromWriter in multiple chunks.
		for _, chunk := range TestStrChunks {
			n, err := rfw.Write([]byte(chunk))
			req.NoError(err)
			assr.Equal(len(chunk), n)
		}

		// Close the pipe writer to signal EOF.
		req.NoError(rfw.Close())

		// Check the data read by ExampleReaderFrom.
		assr.Equal(TestStr, string(erf.data))
	})

	// // Tests the ReaderFromWriter implementation with an empty writer.
	t.Run("", func(t *testing.T) {
		// Create an instance of ExampleReaderFrom.
		erf := &ExampleReaderFrom{}
		rfw := NewReaderFromWriter(erf)

		// Close the pipe writer to signal EOF without writing any data.
		req.NoError(rfw.Close())

		// Check the data read by ExampleReaderFrom.
		assr.Empty(erf.data)
	})
}
