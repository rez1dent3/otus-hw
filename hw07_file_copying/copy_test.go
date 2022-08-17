package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// errorWriter needed for testing write errors.
type errorWriter struct {
	io.Writer
}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func TestIoCopy(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		temp, err := os.CreateTemp(os.TempDir(), "empty")
		require.Nil(t, err)

		result, err := os.CreateTemp(os.TempDir(), "result")
		require.Nil(t, err)

		err = ioCopy(temp, result, 0, 0)
		require.Nil(t, err)

		file, err := os.ReadFile(result.Name())
		require.Nil(t, err)
		require.Equal(t, []byte{}, file)

		require.Nil(t, os.Remove(temp.Name()))
		require.Nil(t, os.Remove(result.Name()))
	})

	t.Run("negative offset", func(t *testing.T) {
		result, err := os.CreateTemp(os.TempDir(), "result")
		require.Nil(t, err)

		err = ioCopy(result, result, -1, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)

		require.Nil(t, os.Remove(result.Name()))
	})

	t.Run("negative limit", func(t *testing.T) {
		result, err := os.CreateTemp(os.TempDir(), "result")
		require.Nil(t, err)

		err = ioCopy(result, result, 0, -1)
		require.Nil(t, err)

		require.Nil(t, os.Remove(result.Name()))
	})

	t.Run("error write", func(t *testing.T) {
		reader := bytes.NewReader(make([]byte, 100))

		err := ioCopy(reader, &errorWriter{}, 0, 100)
		require.ErrorIs(t, err, io.ErrShortWrite)
	})

	t.Run("double buffer", func(t *testing.T) {
		temp, err := os.CreateTemp(os.TempDir(), "short")
		require.Nil(t, err)

		result, err := os.CreateTemp(os.TempDir(), "result")
		require.Nil(t, err)

		byteBuf := make([]byte, 32768*2)
		write, err := temp.Write(byteBuf)
		require.Equal(t, len(byteBuf), write)
		require.Nil(t, err)

		writer := io.MultiWriter(result)
		err = ioCopy(temp, writer, 0, 0)
		require.Nil(t, err)

		file, err := os.ReadFile(result.Name())
		require.Nil(t, err)
		require.Equal(t, byteBuf, file)

		require.Nil(t, os.Remove(temp.Name()))
		require.Nil(t, os.Remove(result.Name()))
	})

	t.Run("limit greater than file size", func(t *testing.T) {
		temp, err := os.CreateTemp(os.TempDir(), "offset")
		require.Nil(t, err)

		result, err := os.CreateTemp(os.TempDir(), "result")
		require.Nil(t, err)

		write, err := temp.Write([]byte{42, 33, 42, 32})
		require.Equal(t, 4, write)
		require.Nil(t, err)

		err = ioCopy(temp, result, 0, 5)
		require.Nil(t, err)

		file, err := os.ReadFile(result.Name())
		require.Nil(t, err)
		require.Equal(t, []byte{42, 33, 42, 32}, file)

		require.Nil(t, os.Remove(temp.Name()))
		require.Nil(t, os.Remove(result.Name()))
	})

	t.Run("offset&limit position from file", func(t *testing.T) {
		temp, err := os.CreateTemp(os.TempDir(), "offset")
		require.Nil(t, err)

		result, err := os.CreateTemp(os.TempDir(), "result")
		require.Nil(t, err)

		write, err := temp.Write([]byte{42, 33, 42, 32})
		require.Equal(t, 4, write)
		require.Nil(t, err)

		err = ioCopy(temp, result, 1, 2)
		require.Nil(t, err)

		file, err := os.ReadFile(result.Name())
		require.Nil(t, err)
		require.Equal(t, []byte{33, 42}, file)

		require.Nil(t, os.Remove(temp.Name()))
		require.Nil(t, os.Remove(result.Name()))
	})

	t.Run("byte write", func(t *testing.T) {
		temp, err := os.CreateTemp(os.TempDir(), "empty")
		require.Nil(t, err)

		result, err := os.CreateTemp(os.TempDir(), "result")
		require.Nil(t, err)

		write, err := temp.Write([]byte{42})
		require.Equal(t, 1, write)
		require.Nil(t, err)

		err = ioCopy(temp, result, 0, 0)
		require.Nil(t, err)

		file, err := os.ReadFile(result.Name())
		require.Nil(t, err)
		require.Equal(t, []byte{42}, file)

		require.Nil(t, os.Remove(temp.Name()))
		require.Nil(t, os.Remove(result.Name()))
	})
}

func TestCopy(t *testing.T) {
	t.Run("offset exceeds file size", func(t *testing.T) {
		temp, err := os.CreateTemp(os.TempDir(), "empty")
		require.Nil(t, err)

		stat, err := temp.Stat()
		require.Nil(t, err)
		require.Equal(t, int64(0), stat.Size())

		err = Copy(temp.Name(), temp.Name(), 1, 0)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)
	})

	t.Run("unsupported file", func(t *testing.T) {
		result, err := os.CreateTemp(os.TempDir(), "result")
		require.Nil(t, err)

		require.NoFileExists(t, "/dev/tt0cp")
		err = Copy("/dev/tt0cp", result.Name(), 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)

		require.Nil(t, os.Remove(result.Name()))
	})

	t.Run("can not create file", func(t *testing.T) {
		// joke about Torvalds.
		// everything we write to /dev/null, Torvalds reads from /dev/all
		err := Copy("/dev/null", "/dev/all", 0, 0)
		require.ErrorIs(t, err, ErrCanNotCreateFile)
	})

	t.Run("offset position from file", func(t *testing.T) {
		temp, err := os.CreateTemp(os.TempDir(), "offset")
		require.Nil(t, err)

		result, err := os.CreateTemp(os.TempDir(), "result")
		require.Nil(t, err)

		write, err := temp.Write([]byte{42, 33, 42, 32})
		require.Equal(t, 4, write)
		require.Nil(t, err)

		err = Copy(temp.Name(), result.Name(), 1, 0)
		require.Nil(t, err)

		file, err := os.ReadFile(result.Name())
		require.Nil(t, err)
		require.Equal(t, []byte{33, 42, 32}, file)

		require.Nil(t, os.Remove(temp.Name()))
		require.Nil(t, os.Remove(result.Name()))
	})
}
