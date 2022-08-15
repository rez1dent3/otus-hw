package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	input, err := os.Open(from)
	if err != nil {
		log.Fatalln(ErrUnsupportedFile)
	}

	output, err := os.Create(to)
	if err != nil {
		_ = input.Close()
		log.Fatalln(ErrCanNotCreateFile)
	}

	max := limit
	if stat, _ := input.Stat(); err == nil {
		max = stat.Size()
	}

	progress := &progressBar{cur: 0, limit: limit, max: max - offset}
	writer := io.MultiWriter(output, progress)

	err = ioCopy(input, writer, offset, limit)
	_ = input.Close()
	_ = output.Close()

	fmt.Println()

	if err != nil {
		log.Fatalln(err)
	}
}
