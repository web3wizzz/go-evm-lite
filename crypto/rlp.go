package crypto

import (
	"encoding/binary"
	"fmt"
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

	// Rule 5: Long list (> 55 bytes payload)
	lenBytes := bigEndianLength(uint64(len(payload)))
	header := byte(0xf7 + len(lenBytes))
	result := append([]byte{header}, lenBytes...)
	return append(result, payload...)
}

// DecodeBytes decodes a single RLP-encoded string item and returns the value and remaining bytes.
func DecodeBytes(data []byte) ([]byte, []byte, error) {
	if len(data) == 0 {
		return nil, nil, fmt.Errorf("rlp: empty input")
	}

	prefix := data[0]

	// Single byte: 0x00-0x7f
	if prefix <= 0x7f {
		return []byte{prefix}, data[1:], nil
	}

	// String prefix: 0x80-0xb7
	if prefix <= 0xb7 {
		length := int(prefix - 0x80)
		if len(data) < 1+length {
			return nil, nil, fmt.Errorf("rlp: invalid short string")
		}
		value := append([]byte(nil), data[1:1+length]...)
		return value, data[1+length:], nil
	}

	// Long string prefix: 0xb8-0xbf
	if prefix <= 0xbf {
		lengthOfLength := int(prefix - 0xb7)
		if len(data) < 1+lengthOfLength {
			return nil, nil, fmt.Errorf("rlp: invalid long string length")
		}
		length := decodeLength(data[1 : 1+lengthOfLength])
		if len(data) < 1+lengthOfLength+length {
			return nil, nil, fmt.Errorf("rlp: truncated long string")
		}
		value := append([]byte(nil), data[1+lengthOfLength:1+lengthOfLength+length]...)
		return value, data[1+lengthOfLength+length:], nil
	}

	// Lists are not byte strings and must be decoded with DecodeList.
	return nil, nil, fmt.Errorf("rlp: list prefix is not a string")
}

// DecodeList decodes an RLP list into a list of byte slices.
func DecodeList(data []byte) ([][]byte, []byte, error) {
	if len(data) == 0 {
		return nil, nil, fmt.Errorf("rlp: empty input")
	}

	prefix := data[0]
	if prefix < 0xc0 {
		return nil, nil, fmt.Errorf("rlp: input is not a list")
	}

	payloadStart := 1
	payloadLen := int(prefix - 0xc0)
	if prefix >= 0xf8 {
		lengthOfLength := int(prefix - 0xf7)
		if len(data) < 1+lengthOfLength {
			return nil, nil, fmt.Errorf("rlp: invalid list length")
		}
		payloadLen = decodeLength(data[1 : 1+lengthOfLength])
		payloadStart = 1 + lengthOfLength
	}

	if len(data) < payloadStart+payloadLen {
		return nil, nil, fmt.Errorf("rlp: truncated list")
	}

	payload := data[payloadStart : payloadStart+payloadLen]
	items := make([][]byte, 0)
	for len(payload) > 0 {
		item, rest, err := DecodeBytes(payload)
		if err != nil {
			return nil, nil, err
		}
		items = append(items, item)
		payload = rest
	}

	return items, data[payloadStart+payloadLen:], nil
}

// Helper: Returns minimal big-endian representation of a length/uint
func bigEndianLength(val uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	for i, b := range buf {
		if b != 0 {
			return buf[i:]
		}
	}
	return []byte{}
}

func decodeLength(data []byte) int {
	if len(data) == 0 {
		return 0
	}

	if len(data) > 8 {
		lenData := make([]byte, 8)
		copy(lenData[8-len(data):], data)
		return int(binary.BigEndian.Uint64(lenData))
	}

	buf := make([]byte, 8)
	copy(buf[8-len(data):], data)
	return int(binary.BigEndian.Uint64(buf))
}
