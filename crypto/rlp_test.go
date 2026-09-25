package crypto

import (
	"bytes"
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

func TestRLPDecodeBytes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []byte
		rest  string
	}{
		{name: "Empty string", input: "80", want: []byte{}, rest: ""},
		{name: "Single byte", input: "0f", want: []byte{0x0f}, rest: ""},
		{name: "Short string", input: "83646f67", want: []byte("dog"), rest: ""},
		{name: "Trailing bytes", input: "83646f67ff", want: []byte("dog"), rest: "ff"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := hex.DecodeString(tt.input)
			if err != nil {
				t.Fatalf("hex decode failed: %v", err)
			}

			restHex, _ := hex.DecodeString(tt.rest)
			got, remaining, err := DecodeBytes(encoded)
			if err != nil {
				t.Fatalf("DecodeBytes returned error: %v", err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Fatalf("DecodeBytes() = %x, want %x", got, tt.want)
			}
			if !bytes.Equal(remaining, restHex) {
				t.Fatalf("DecodeBytes() remaining = %x, want %x", remaining, restHex)
			}
		})
	}
}

func TestRLPDecodeList(t *testing.T) {
	encoded, err := hex.DecodeString("c883646f6783676f64")
	if err != nil {
		t.Fatalf("hex decode failed: %v", err)
	}

	got, remaining, err := DecodeList(encoded)
	if err != nil {
		t.Fatalf("DecodeList returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("DecodeList() length = %d, want 2", len(got))
	}
	if !bytes.Equal(got[0], []byte("dog")) || !bytes.Equal(got[1], []byte("god")) {
		t.Fatalf("DecodeList() = %x, %x; want dog, god", got[0], got[1])
	}
	if len(remaining) != 0 {
		t.Fatalf("DecodeList() remaining = %x, want empty", remaining)
	}
}

func TestRLPRoundTrip(t *testing.T) {
	values := [][]byte{
		[]byte{},
		[]byte{0x0f},
		[]byte("dog"),
		[]byte("hello world"),
	}

	for _, v := range values {
		encoded := EncodeBytes(v)
		decoded, remaining, err := DecodeBytes(encoded)
		if err != nil {
			t.Fatalf("DecodeBytes(%x) error: %v", v, err)
		}
		if !bytes.Equal(decoded, v) {
			t.Fatalf("roundtrip failed for %x: got %x want %x", v, decoded, v)
		}
		if len(remaining) != 0 {
			t.Fatalf("roundtrip remaining bytes for %x: %x", v, remaining)
		}
	}
}
