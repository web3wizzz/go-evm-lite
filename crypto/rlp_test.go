package crypto

import (
	"encoding/hex"
	"testing"
)

func TestRLPEncoding(t *testing.T) {
	tests := []struct {
		name     string
		got      []byte
		expected string
	}{
		{
			name:     "Empty String / Zero Byte",
			got:      EncodeBytes([]byte{}),
			expected: "80",
		},
		{
			name:     "Single byte < 0x80",
			got:      EncodeBytes([]byte{0x0f}),
			expected: "0f",
		},
		{
			name:     "Short String 'dog'",
			got:      EncodeBytes([]byte("dog")),
			expected: "83646f67",
		},
		{
			name:     "Empty List",
			got:      EncodeList([][]byte{}),
			expected: "c0",
		},
		{
			name:     "List of ['dog', 'god']",
			got:      EncodeList([][]byte{EncodeBytes([]byte("dog")), EncodeBytes([]byte("god"))}),
			expected: "c883646f6783676f64",
		},
		{
			name:     "Uint 0",
			got:      EncodeUint(0),
			expected: "80",
		},
		{
			name:     "Uint 15 (0x0f)",
			got:      EncodeUint(15),
			expected: "0f",
		},
		{
			name:     "Uint 1024 (0x0400)",
			got:      EncodeUint(1024),
			expected: "820400",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHex := hex.EncodeToString(tt.got)
			if gotHex != tt.expected {
				t.Errorf("%s failed: got %s, want %s", tt.name, gotHex, tt.expected)
			}
		})
	}
}
