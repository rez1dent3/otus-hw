package hw09structvalidator

import (
	"strconv"
	"strings"
)

func intMin(v interface{}, args string) error {
	i, ok := v.(int)
	if !ok {
		return nil
	}

	min, err := strconv.Atoi(args)
	if err != nil {
		return ErrFailedReadArgument
	}

	if i < min {
		return ErrNumberIsLessThanMinimumAllowed
	}

	return nil
}

func intMax(v interface{}, args string) error {
	i, ok := v.(int)
	if !ok {
		return nil
	}

	max, err := strconv.Atoi(args)
	if err != nil {
		return ErrFailedReadArgument
	}

	if i > max {
		return ErrNumberIsGreaterThanMaximumAllowed
	}

	return nil
}

func intIn(v interface{}, args string) error {
	i, ok := v.(int)
	if !ok {
		return nil
	}

	values := strings.Split(args, ",")
	for _, value := range values {
		if value == strconv.Itoa(i) {
			return nil
		}
	}

	return ErrNumberIsNotIncludedInSet
}
