package state

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestHexPrefixEncoding(t *testing.T) {
	tests := []struct {
		name     string
		nibbles  []byte
		isLeaf   bool
		expected string
	}{
		{
			name:     "Even length Extension node",
			nibbles:  []byte{1, 2, 3, 4},
			isLeaf:   false,
			expected: "001234",
		},
		{
			name:     "Odd length Extension node",
			nibbles:  []byte{1, 2, 3, 4, 5},
			isLeaf:   false,
			expected: "112345",
		},
		{
			name:     "Even length Leaf node",
			nibbles:  []byte{0, 1, 2, 3, 4, 5},
			isLeaf:   true,
			expected: "20012345",
		},
		{
			name:     "Odd length Leaf node",
			nibbles:  []byte{1, 2, 3, 4, 5},
			isLeaf:   true,
			expected: "312345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := HexPrefixEncode(tt.nibbles, tt.isLeaf)
			encodedHex := hex.EncodeToString(encoded)
			if encodedHex != tt.expected {
				t.Errorf("%s failed: got %s, want %s", tt.name, encodedHex, tt.expected)
			}

			// Roundtrip decode check
			decodedNibbles, decodedLeaf, err := HexPrefixDecode(encoded)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			if decodedLeaf != tt.isLeaf {
				t.Errorf("Leaf flag mismatch: got %v, want %v", decodedLeaf, tt.isLeaf)
			}

			if !bytes.Equal(decodedNibbles, tt.nibbles) {
				t.Errorf("Nibbles mismatch: got %v, want %v", decodedNibbles, tt.nibbles)
			}
		})
	}
}

func TestBytesToNibbles(t *testing.T) {
	input := []byte{0xab, 0xcd}
	expected := []byte{10, 11, 12, 13}

	got := BytesToNibbles(input)
	if !bytes.Equal(got, expected) {
		t.Errorf("BytesToNibbles failed: got %v, want %v", got, expected)
	}
}
