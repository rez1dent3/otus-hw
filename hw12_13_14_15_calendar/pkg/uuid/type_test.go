package uuid_test

import (
	"encoding/hex"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUUID_ToString(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		val := uuid.UUID{}
		require.Equal(t, "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00", val.ToString())
		require.Equal(t, "00000000000000000000000000000000", hex.EncodeToString(val.ToBytes()))
	})

	t.Run("create from string", func(t *testing.T) {
		str := "1f8c3858c5f948b4bc1186cd4ad4f1e4"
		val := uuid.FromString(str)

		require.Equal(t, "\x1f\x8c\x38\x58\xc5\xf9\x48\xb4\xbc\x11\x86\xcd\x4a\xd4\xf1\xe4", val.ToString())
		require.Equal(t, "1f8c3858c5f948b4bc1186cd4ad4f1e4", hex.EncodeToString(val.ToBytes()))
	})

	t.Run("uuid check", func(t *testing.T) {
		bytesInput := []byte("\x1f\x8c\x38\x58\xc5\xf9\x48\xb4\xbc\x11\x86\xcd\x4a\xd4\xf1\xe4")

		require.Len(t, bytesInput, 16)
		require.Equal(t, "1f8c3858c5f948b4bc1186cd4ad4f1e4", hex.EncodeToString(bytesInput))

		actual := uuid.FromBytes(bytesInput)

		require.Equal(t, "\x1f\x8c\x38\x58\xc5\xf9\x48\xb4\xbc\x11\x86\xcd\x4a\xd4\xf1\xe4", actual.ToString())
		require.Equal(t, "1f8c3858c5f948b4bc1186cd4ad4f1e4", hex.EncodeToString(actual.ToBytes()))
	})
}
