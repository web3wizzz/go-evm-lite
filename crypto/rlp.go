package crypto

import (
	"encoding/binary"
)

// EncodeBytes encodes a byte slice according to RLP rules.

func EncodeBytes(b []byte) []byte {
	// Rule 1: Single byte < 0x80
	if len(b) == 1 && b[0] < 0x80 {
		return b
	}

	// Rule 2: Short byte array (0-55 bytes)
	if len(b) <= 55 {
		header := byte(0x80 + len(b))
		return append([]byte{header}, b...)
	}

	// Rule 3: Long byte array (> 55 bytes)
	lenBytes := bigEndianLength(uint64(len(b)))
	header := byte(0xb7 + len(lenBytes))
	result := append([]byte{header}, lenBytes...)
	return append(result, b...)

}

// EncodeUint encodes a uint64 into RLP format.
func EncodeUint(val uint64) []byte {
	if val == 0 {
		return EncodeBytes([]byte{})
	}
	return EncodeBytes(bigEndianLength(val))
}

// EncodeList encodes a slice of already-encoded RLP items into a list.
func EncodeList(items [][]byte) []byte {
	var payload []byte
	for _, item := range items {
		payload = append(payload, item...)
	}

	// Rule 4: Short list (0-55 bytes payload)
	if len(payload) <= 55 {
		header := byte(0xc0 + len(payload))
		return append([]byte{header}, payload...)
	}

	// Rule 5: Long list (&gt; 55 bytes payload)
	lenBytes := bigEndianLength(uint64(len(payload)))
	header := byte(0xf7 + len(lenBytes))
	result := append([]byte{header}, lenBytes...)
	return append(result, payload...)
}

// Helper: Returns minimal big-endian representation of a length/uint
func bigEndianLength(val uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	// Strip leading zeros
	for i, b := range buf {
		if b != 0 {
			return buf[i:]
		}
	}
	return []byte{}
}
