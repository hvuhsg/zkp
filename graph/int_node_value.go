package graph

import (
	"encoding/binary"
	"errors"
)

var ErrDataTooShort = errors.New("data too short for uint16 value")

type IntNodeValue uint16

func (v IntNodeValue) Serialize() []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(v))
	return bytes
}

func DeserializeIntNodeValue(data []byte) (IntNodeValue, error) {
	if len(data) < 2 {
		return 0, ErrDataTooShort
	}
	return IntNodeValue(binary.BigEndian.Uint16(data[:2])), nil
}
