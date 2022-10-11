package uuid

type UUID [16]byte

func FromBytes(bytes []byte) UUID {
	var result UUID
	for key := range result {
		result[key] = bytes[key]
	}

	return result
}

func (u *UUID) ToBytes() []byte {
	return u[:]
}

func (u *UUID) ToString() string {
	return string(u.ToBytes())
}
