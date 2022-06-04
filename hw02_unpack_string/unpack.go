package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	const ecs = '\\'
	response := ""
	escape := false
	var buffer *int32
	for _, ch := range input {
		ch := ch
		if buffer == nil && unicode.IsDigit(ch) {
			return "", ErrInvalidString
		}

		if escape && !unicode.IsDigit(ch) && ch != ecs {
			return "", ErrInvalidString
		}

		if !escape && unicode.IsDigit(ch) {
			// count, _ := strconv.Atoi(string(ch))
			// strings.Repeat(string(*buffer), count)
			response += strings.Repeat(string(*buffer), int(ch)-'0')
			buffer = nil
			continue
		}

		if !escape && buffer != nil {
			response += string(*buffer)
		}

		escape = !escape && ch == ecs
		buffer = &ch
	}

	if escape {
		return "", ErrInvalidString
	}

	if buffer != nil {
		response += string(*buffer)
	}

	return response, nil
}
