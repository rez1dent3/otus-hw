package uuid

import (
	"encoding/hex"
)

type UUID [16]byte

func FromBytes(input []byte) UUID {
	var result UUID
	copy(result[:], input)

	return result
}

func FromString(input string) UUID {
	inputBytes, _ := hex.DecodeString(input)

	return FromBytes(inputBytes)
}

func (u *UUID) ToBytes() []byte {
	return u[:]
}

func (u *UUID) ToString() string {
	return string(u.ToBytes())
}
