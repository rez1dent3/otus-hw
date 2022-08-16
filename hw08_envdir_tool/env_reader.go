package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	results := make(Environment)
	for _, file := range files {
		if file.Type().IsDir() {
			continue
		}

		if strings.Contains(file.Name(), "=") {
			continue
		}

		fileHandle, err := os.Open(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(fileHandle)
		scanner.Split(bufio.ScanLines)

		scanner.Scan()
		value := scanner.Text()
		value = strings.TrimRight(value, " \t")
		value = strings.ReplaceAll(value, string([]byte{0x00}), "\n")

		results[file.Name()] = EnvValue{Value: value, NeedRemove: value == ""}

		err = fileHandle.Close()
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}
