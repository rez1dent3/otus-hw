package main

import (
	"bufio"
	"errors"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrCanNotCreateFile      = errors.New("Can not create file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func reads(f *os.File, limit int64) <-chan byte {
	ch := make(chan byte)
	go func() {
		defer close(ch)
		reader := bufio.NewReader(f)

		var i int64 = 0
		for ; i < limit || limit == 0; i++ {
			chr, err := reader.ReadByte()
			if err != nil {
				return
			}

			ch <- chr
		}
	}()

	return ch
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	input, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(input)

	_, err = input.Seek(offset, 0)
	if err != nil {
		return ErrOffsetExceedsFileSize
	}

	output, err := os.Create(toPath)
	if err != nil {
		return ErrCanNotCreateFile
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(output)

	for char := range reads(input, limit) {
		_, err = output.Write([]byte{char})
		if err != nil {
			return err
		}
	}

	return nil
}
