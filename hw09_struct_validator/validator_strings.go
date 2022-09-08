package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var ErrUnableBuildRegExp = errors.New("unable to build regular expression")

func stringLength(v interface{}, args string) error {
	value, err := strconv.ParseInt(args, 0, 64)
	if err != nil {
		return ErrFailedReadArgument
	}

	if inputData := stringNullable(v); inputData != nil && value != int64(len(*inputData)) {
		return ErrIncorrectLength
	}

	return nil
}

func stringIn(v interface{}, args string) error {
	if inputData := stringNullable(v); inputData != nil {
		values := strings.Split(args, ",")
		for _, value := range values {
			if value == *inputData {
				return nil
			}
		}

		return ErrStringIsNotIncludedInSet
	}

	return nil
}

func stringRegexp(v interface{}, expr string) error {
	compile, err := regexp.Compile(expr)
	if err != nil {
		return ErrUnableBuildRegExp
	}

	if inputData := stringNullable(v); inputData != nil && !compile.Match([]byte(*inputData)) {
		return ErrStringNotMatchRegExp
	}

	return nil
}

func stringNullable(input interface{}) *string {
	if reflect.TypeOf(input).Kind() == reflect.String {
		result := reflect.ValueOf(input).String()
		return &result
	}
	return nil
}
