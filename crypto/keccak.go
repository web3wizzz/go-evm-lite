package crypto

import (
	"golang.org/x/crypto/sha3"
)

// Keccak256 calculates and returns the 32-byte Keccak-256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}
