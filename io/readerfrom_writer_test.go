package io_test

import (
	"io"
	"testing"

	gio "github.com/daotl/guts/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testReaderFromWriter(
	t *testing.T,
	initRfw func(rf io.ReaderFrom) gio.ReadFromWriteCloser,
) {
	req := require.New(t)
	assr := assert.New(t)

	// Tests the ReaderFromWriter implementation with single write.
	t.Run("Single write", func(t *testing.T) {
		// Create an instance of ExampleReaderFrom and an instance of ReaderFromWriter.
		erf := &ExampleReaderFrom{}
		rfw := initRfw(erf)

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
		erf := &ExampleReaderFrom{}
		rfw := initRfw(erf)

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

	// Tests the ReaderFromWriter implementation with no write.
	t.Run("No write", func(t *testing.T) {
		erf := &ExampleReaderFrom{}
		rfw := initRfw(erf)

		// Close the pipe writer to signal EOF without writing any data.
		req.NoError(rfw.Close())

		// Check the data read by ExampleReaderFrom.
		assr.Empty(erf.data)
	})
}

// TestReaderFromWriter tests the ReaderFromWriter implementation.
func TestReaderFromWriter(t *testing.T) {
	testReaderFromWriter(
		t,
		func(rf io.ReaderFrom) gio.ReadFromWriteCloser { return gio.NewReaderFromWriter(rf) },
	)
}
