package uuid

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type UUID [16]byte

func Gen() UUID {
	input := make([]byte, 16)
	_, err := rand.Read(input)
	if err != nil {
		return [16]byte{}
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
	switch src.(type) {
	case string:
		*u = FromString(src.(string))
		return nil
	case []byte:
		*u = FromBytes(src.([]byte))
		return nil
	case UUID:
		*u = src.(UUID)
		return nil
	default:
		return errors.New("invalid type")
	}
}

func (u UUID) Value() (driver.Value, error) {
	return u.String(), nil
}
