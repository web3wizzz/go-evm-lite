package state

import (
	"bytes"
	"math/big"
	"testing"
)

func TestStateDBAccountOperations(t *testing.T) {
	sdb := NewStateDB()
	addr := []byte("0x12345678901234567890")

	// Test 1: Non-existent account returns default zero values
	bal := sdb.GetBalance(addr)
	if bal.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("expected balance 0 for non-existent account, got %s", bal.String())
	}

	nonce := sdb.GetNonce(addr)
	if nonce != 0 {
		t.Errorf("expected nonce 0 for non-existent account, got %d", nonce)
	}

	// Test 2: SetBalance and GetBalance
	newBal := big.NewInt(500000000000000000) // 0.5 ETH in wei
	err := sdb.SetBalance(addr, newBal)
	if err != nil {
		t.Fatalf("SetBalance failed: %v", err)
	}

	gotBal := sdb.GetBalance(addr)
	if gotBal.Cmp(newBal) != 0 {
		t.Errorf("GetBalance = %s, want %s", gotBal.String(), newBal.String())
	}

	// Test 3: SetNonce and GetNonce
	err = sdb.SetNonce(addr, 5)
	if err != nil {
		t.Fatalf("SetNonce failed: %v", err)
	}

	gotNonce := sdb.GetNonce(addr)
	if gotNonce != 5 {
		t.Errorf("GetNonce = %d, want 5", gotNonce)
	}

	// Test 4: Retrieve updated full Account object
	acc, err := sdb.GetAccount(addr)
	if err != nil {
		t.Fatalf("GetAccount failed: %v", err)
	}

	if acc.Nonce != 5 {
		t.Errorf("acc.Nonce = %d, want 5", acc.Nonce)
	}

	if !bytes.Equal(acc.StorageRoot, EmptyRootHash) {
		t.Errorf("acc.StorageRoot mismatch, got %x", acc.StorageRoot)
	}

	if !bytes.Equal(acc.CodeHash, EmptyCodeHash) {
		t.Errorf("acc.CodeHash mismatch, got %x", acc.CodeHash)
	}
}
