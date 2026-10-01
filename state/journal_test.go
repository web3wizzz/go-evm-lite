package state

import (
	"bytes"
	"math/big"
	"testing"
)

func TestStateJournalAndRevert(t *testing.T) {
	sdb := NewStateDB()
	addr := []byte("0x12345678901234567890")
	slotKey := []byte("slot_0")

	// Set initial state values
	initBal := big.NewInt(1000)
	sdb.SetBalance(addr, initBal)
	sdb.SetNonce(addr, 1)
	sdb.SetState(addr, slotKey, []byte("val_0___________________________"))

	// Take Snapshot 0
	snap0 := sdb.Snapshot()

	// Make mutations after Snapshot 0
	sdb.SetBalance(addr, big.NewInt(5000))
	sdb.SetNonce(addr, 2)
	sdb.SetState(addr, slotKey, []byte("val_1___________________________"))

	// Take Snapshot 1
	snap1 := sdb.Snapshot()

	// Make further mutations after Snapshot 1
	sdb.SetBalance(addr, big.NewInt(9000))
	sdb.SetNonce(addr, 3)
	sdb.SetState(addr, slotKey, []byte("val_2___________________________"))

	// Verify values before revert
	if sdb.GetBalance(addr).Cmp(big.NewInt(9000)) != 0 {
		t.Errorf("expected balance 9000 before revert, got %s", sdb.GetBalance(addr).String())
	}

	// Revert to Snapshot 1
	sdb.RevertToSnapshot(snap1)

	if sdb.GetBalance(addr).Cmp(big.NewInt(5000)) != 0 {
		t.Errorf("expected balance 5000 after revert to snap1, got %s", sdb.GetBalance(addr).String())
	}
	if sdb.GetNonce(addr) != 2 {
		t.Errorf("expected nonce 2 after revert to snap1, got %d", sdb.GetNonce(addr))
	}
	if !bytes.Equal(sdb.GetState(addr, slotKey), []byte("val_1___________________________")) {
		t.Errorf("expected slot val_1 after revert to snap1, got %s", string(sdb.GetState(addr, slotKey)))
	}

	// Revert to Snapshot 0
	sdb.RevertToSnapshot(snap0)

	if sdb.GetBalance(addr).Cmp(initBal) != 0 {
		t.Errorf("expected balance 1000 after revert to snap0, got %s", sdb.GetBalance(addr).String())
	}
	if sdb.GetNonce(addr) != 1 {
		t.Errorf("expected nonce 1 after revert to snap0, got %d", sdb.GetNonce(addr))
	}
	if !bytes.Equal(sdb.GetState(addr, slotKey), []byte("val_0___________________________")) {
		t.Errorf("expected slot val_0 after revert to snap0, got %s", string(sdb.GetState(addr, slotKey)))
	}
}
