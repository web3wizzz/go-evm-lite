package state

import (
	"bytes"
	"math/big"
	"testing"
)

func TestAccountCreation(t *testing.T) {
	balance := big.NewInt(1000000000000000000) // 1 ETH in wei
	acc := NewAccount(1, balance)

	if acc.Nonce != 1 {
		t.Errorf("expected nonce 1, got %d", acc.Nonce)
	}

	if acc.Balance.Cmp(balance) != 0 {
		t.Errorf("expected balance %s, got %s", balance.String(), acc.Balance.String())
	}

	if !bytes.Equal(acc.StorageRoot, EmptyRootHash) {
		t.Errorf("expected empty storage root, got %x", acc.StorageRoot)
	}

	if !bytes.Equal(acc.CodeHash, EmptyCodeHash) {
		t.Errorf("expected empty code hash, got %x", acc.CodeHash)
	}

	if acc.IsContract() {
		t.Errorf("expected IsContract() to be false for EOA")
	}
}

func TestAccountIsContract(t *testing.T) {
	acc := NewAccount(0, big.NewInt(0))

	// Set a custom non-empty contract code hash
	customCodeHash := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
	}
	acc.CodeHash = customCodeHash

	if !acc.IsContract() {
		t.Errorf("expected IsContract() to be true for contract account")
	}
}
