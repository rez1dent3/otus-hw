package uuid

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

type UUID [16]byte

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func Gen() UUID {
	input := make([]byte, 16)
	if _, err := rand.Read(input); err != nil {
		return UUID{}
	}

	return FromBytes(input)
}

func FromBytes(input []byte) UUID {
	var result UUID
	copy(result[:], input)

	return result
}

func FromString(input string) UUID {
	inputBytes, _ := hex.DecodeString(strings.ReplaceAll(input, "-", ""))

	return FromBytes(inputBytes)
}

func (u *UUID) ToBytes() []byte {
	return u[:]
}

func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
}

func (u *UUID) Scan(src interface{}) error {
	if val, ok := src.(string); ok {
		*u = FromString(val)
		return nil
	}

	if val, ok := src.([]byte); ok {
		*u = FromBytes(val)
		return nil
	}

	return errors.New("invalid type")
}

func (u UUID) Value() (driver.Value, error) {
	return u.String(), nil
}
