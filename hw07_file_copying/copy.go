package main

import (
	"bufio"
	"errors"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrCanNotCreateFile      = errors.New("can not create file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type bar struct {
	cur, max int64
}

func (b *bar) inc() {
	b.cur++
}

func (b *bar) display() int64 {
	return b.max
}

func reads(f *os.File, offset, limit int64) (<-chan byte, error) {
	if _, err := f.Seek(offset, 0); err != nil {
		return nil, ErrOffsetExceedsFileSize
	}

	ch := make(chan byte)
	go func() {
		defer close(ch)
		reader := bufio.NewReader(f)

		var i int64
		for ; i < limit || limit <= 0; i++ {
			chr, err := reader.ReadByte()
			if err != nil {
				return
			}

			ch <- chr
		}
	}()

	return ch, nil
}

func copyWithProgressBar(fromPath, toPath string, offset, limit int64, _ *bar) error {
	input, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(input)

	output, err := os.Create(toPath)
	if err != nil {
		return ErrCanNotCreateFile
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(output)

	ch, err := reads(input, offset, limit)
	if err != nil {
		return err
	}

	for char := range ch {
		_, err = output.Write([]byte{char})
		if err != nil {
			return err
		}
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	return copyWithProgressBar(fromPath, toPath, offset, limit, nil)
}
