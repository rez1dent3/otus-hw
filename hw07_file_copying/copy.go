package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrCanNotCreateFile      = errors.New("can not create file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func ioCopy(reader io.ReadSeeker, writer io.Writer, offset, limit int64) error {
	if _, err := reader.Seek(offset, io.SeekStart); err != nil {
		return ErrOffsetExceedsFileSize
	}

	var ioReader io.Reader = reader
	if limit > 0 {
		ioReader = io.LimitReader(reader, limit)
	}

	if _, err := io.Copy(writer, ioReader); err != nil {
		return err
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	input, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(input)

	if stat, err := input.Stat(); err != nil {
		return err
	} else if stat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	output, err := os.Create(toPath)
	if err != nil {
		return ErrCanNotCreateFile
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(output)

	return ioCopy(input, output, offset, limit)
}
