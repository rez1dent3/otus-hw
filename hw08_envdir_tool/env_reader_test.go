package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("testdata", func(t *testing.T) {
		env, err := ReadDir("testdata/env")

		results := Environment{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: true},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: "\"hello\"", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		}

		require.Nil(t, err)
		require.Truef(t, reflect.DeepEqual(env, results), "actual - %#v\nexcepted - %#v", env, results)
	})
}
