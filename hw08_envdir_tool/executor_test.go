package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("printenv $?=0", func(t *testing.T) {
		var buffer strings.Builder
		exitCode := runCmd(&buffer, []string{"/bin/bash", "testdata/printenv.sh"}, Environment{
			"MSG": {Value: "\"hello world\""},
		})

		require.Equal(t, 0, exitCode)
		require.Equal(t, "\"hello world\"\n", buffer.String())
	})

	t.Run("printenv $?=1", func(t *testing.T) {
		var buffer strings.Builder
		exitCode := runCmd(&buffer, []string{"/bin/bash", "testdata/printenv_foo.sh"}, nil)

		require.Equal(t, 1, exitCode)
		require.Equal(t, "", buffer.String())
	})

	t.Run("exit 1; $?=1", func(t *testing.T) {
		exitCode := runCmd(nil, []string{"/bin/bash", "testdata/exit1.sh"}, nil)

		require.Equal(t, 1, exitCode)
	})
}
