// Copyright (c) Tendermint. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that can be
// found in the LICENSE_APACHE_2.0 file.

package os_test

// From: https://github.com/tendermint/tendermint/blob/f28d629e280ddcdc0dd644ccf1586d73dddfb7a1/libs/os/os_test.go

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	gos "github.com/daotl/guts/os"
)

func TestCopyFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	content := []byte("hello world")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}

	copyfile := fmt.Sprintf("%s.copy", tmpfile.Name())
	if err := gos.CopyFile(tmpfile.Name(), copyfile); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(copyfile); os.IsNotExist(err) {
		t.Fatal("copy should exist")
	}
	data, err := os.ReadFile(copyfile)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, content) {
		t.Fatalf("copy file content differs: expected %v, got %v", content, data)
	}
	os.Remove(copyfile)
}

func TestTrapSignal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip()
	}

	if os.Getenv("TM_TRAP_SIGNAL_TEST") == "1" {
		t.Log("inside test process")
		killer()
		return
	}

	cmd, _, mockStderr := newTestProgram(t, "TM_TRAP_SIGNAL_TEST")

	err := cmd.Run()
	if err == nil {
		wantStderr := "exiting"
		if mockStderr.String() != wantStderr {
			t.Fatalf("stderr: want %q, got %q", wantStderr, mockStderr.String())
		}

		return
	}

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		t.Fatalf("wrong exit code, want 0, got %d", e.ExitCode())
	}

	t.Fatal("this error should not be triggered")
}

func TestEnsureDir(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ensure-dir")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	// Should be possible to create a new directory.
	err = gos.EnsureDir(filepath.Join(tmp, "dir"), 0755)
	require.NoError(t, err)
	require.DirExists(t, filepath.Join(tmp, "dir"))

	// Should succeed on existing directory.
	err = gos.EnsureDir(filepath.Join(tmp, "dir"), 0755)
	require.NoError(t, err)

	// Should fail on file.
	err = os.WriteFile(filepath.Join(tmp, "file"), []byte{}, 0644)
	require.NoError(t, err)
	err = gos.EnsureDir(filepath.Join(tmp, "file"), 0755)
	require.Error(t, err)

	// Should allow symlink to dir.
	err = os.Symlink(filepath.Join(tmp, "dir"), filepath.Join(tmp, "linkdir"))
	require.NoError(t, err)
	err = gos.EnsureDir(filepath.Join(tmp, "linkdir"), 0755)
	require.NoError(t, err)

	// Should error on symlink to file.
	err = os.Symlink(filepath.Join(tmp, "file"), filepath.Join(tmp, "linkfile"))
	require.NoError(t, err)
	err = gos.EnsureDir(filepath.Join(tmp, "linkfile"), 0755)
	require.Error(t, err)
}

type mockLogger struct{}

func (ml mockLogger) Info(args ...any) {}

func killer() {
	logger := mockLogger{}

	gos.TrapSignal(logger, func() { _, _ = fmt.Fprintf(os.Stderr, "exiting") })
	time.Sleep(1 * time.Second)

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		panic(err)
	}

	if err := p.Signal(syscall.SIGTERM); err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)
}

func newTestProgram(
	t *testing.T,
	environVar string,
) (cmd *exec.Cmd, stdout *bytes.Buffer, stderr *bytes.Buffer) {
	t.Helper()

	cmd = exec.Command(os.Args[0], "-test.run="+t.Name())
	stdout, stderr = bytes.NewBufferString(""), bytes.NewBufferString("")
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=1", environVar))
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return
}

// Ensure that using CopyFile does not truncate the destination file before
// the origin is positively a non-directory and that it is ready for copying.
// See https://github.com/tendermint/tendermint/issues/6427
func TestTrickedTruncation(t *testing.T) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "pwn_truncate")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpDir)

	originalWALPath := filepath.Join(tmpDir, "wal")
	originalWALContent := []byte("I AM BECOME DEATH, DESTROYER OF ALL WORLDS!")
	if err := os.WriteFile(originalWALPath, originalWALContent, 0755); err != nil {
		t.Fatal(err)
	}

	// 1. Sanity check.
	readWAL, err := os.ReadFile(originalWALPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(readWAL, originalWALContent) {
		t.Fatalf(
			"Cannot proceed as the content does not match\nGot:  %q\nWant: %q",
			readWAL,
			originalWALContent,
		)
	}

	// 2. Now cause the truncation of the original file.
	// It is absolutely legal to invoke os.Open on a directory.
	if err := gos.CopyFile(tmpDir, originalWALPath); err == nil {
		t.Fatal("Expected an error")
	}

	// 3. Check the WAL's content
	reReadWAL, err := os.ReadFile(originalWALPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(reReadWAL, originalWALContent) {
		t.Fatalf(
			"Oops, the WAL's content was changed :(\nGot:  %q\nWant: %q",
			reReadWAL,
			originalWALContent,
		)
	}
}
