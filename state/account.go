package state

import (
	"math/big"
)

// EmptyRootHash is the Keccak-256 hash of an empty Merkle Patricia Trie
var EmptyRootHash = []byte{
	0x56, 0xe8, 0x1f, 0x17, 0x1b, 0xcc, 0x55, 0xa6,
	0xff, 0x83, 0x45, 0xe6, 0x92, 0xc0, 0xf8, 0x6e,
	0x5b, 0x48, 0xe0, 0x1b, 0x99, 0x6c, 0xad, 0xc0,
	0x01, 0x62, 0x2f, 0xb5, 0xe3, 0x63, 0xb4, 0x21,
}

// EmptyCodeHash is the Keccak-256 hash of empty bytecode (e.g. for EOAs)
var EmptyCodeHash = []byte{
	0xc5, 0xd2, 0x46, 0x01, 0x86, 0xf7, 0x23, 0x3c,
	0x92, 0x7e, 0x7d, 0xb2, 0xdc, 0xc7, 0x03, 0xc0,
	0xe5, 0x00, 0xb6, 0x53, 0xca, 0x82, 0x27, 0x3b,
	0x7b, 0xad, 0xd6, 0xec, 0x15, 0xe7, 0x37, 0x47,
}

// Account represents the Ethereum account state tuple: [Nonce, Balance, StorageRoot, CodeHash]
type Account struct {
	Nonce       uint64
	Balance     *big.Int
	StorageRoot []byte
	CodeHash    []byte
}

// NewAccount creates an account initialized with default storage root and code hash
func NewAccount(nonce uint64, balance *big.Int) *Account {
	if balance == nil {
		balance = new(big.Int)
	}
	return &Account{
		Nonce:       nonce,
		Balance:     balance,
		StorageRoot: append([]byte(nil), EmptyRootHash...),
		CodeHash:    append([]byte(nil), EmptyCodeHash...),
	}
}

// IsContract returns true if the account contains deployed contract bytecode
func (a *Account) IsContract() bool {
	if len(a.CodeHash) != 32 {
		return false
	}
	for i, b := range a.CodeHash {
		if b != EmptyCodeHash[i] {
			return true
		}
	}
	return false
}
