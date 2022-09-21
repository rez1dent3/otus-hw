package logger_test

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/logger"
)

func TestLogger(t *testing.T) {
	t.Run("check level", func(t *testing.T) {
		data := []struct {
			level    string
			expected string
		}{
			{
				level:    "off",
				expected: "",
			},
			{
				level:    "error",
				expected: "error\n",
			},
			{
				level:    "warning",
				expected: "error\nwarning\n",
			},
			{
				level:    "info",
				expected: "error\nwarning\ninfo\n",
			},
			{
				level:    "debug",
				expected: "error\nwarning\ninfo\ndebug\n",
			},
		}

		for _, datum := range data {
			buffer := &bytes.Buffer{}
			log := logger.New(datum.level, buffer)

			log.Error("error")
			log.Warning("warning")
			log.Info("info")
			log.Debug("debug")

			require.Equal(t, datum.expected, buffer.String())
		}
	})
}
