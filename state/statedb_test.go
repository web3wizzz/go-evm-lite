
package state

import (
	"bytes"
	"testing"
)

func TestStateDBStorageSlots(t *testing.T) {
	sdb := NewStateDB()
	addr1 := []byte("0x11111111111111111111")
	addr2 := []byte("0x22222222222222222222")

	slotKey := []byte("storage_slot_0001")
	value1 := []byte("value_for_contract_1___________")
	value2 := []byte("value_for_contract_2___________")

	// Test 1: Uninitialized storage slot returns 32 zero bytes
	uninitVal := sdb.GetState(addr1, slotKey)
	if len(uninitVal) != 32 {
		t.Errorf("expected 32 bytes for uninitialized state, got %d", len(uninitVal))
	}
	empty32 := make([]byte, 32)
	if !bytes.Equal(uninitVal, empty32) {
		t.Errorf("expected zero bytes for uninitialized storage, got %x", uninitVal)
	}

	// Test 2: SetState and GetState for addr1
	err := sdb.SetState(addr1, slotKey, value1)
	if err != nil {
		t.Fatalf("SetState failed: %v", err)
	}

	gotVal1 := sdb.GetState(addr1, slotKey)
	if !bytes.Equal(gotVal1, value1) {
		t.Errorf("GetState(addr1) = %s, want %s", string(gotVal1), string(value1))
	}

	// Test 3: Storage isolation - addr2 should still return uninitialized zero bytes for same slotKey
	gotVal2 := sdb.GetState(addr2, slotKey)
	if !bytes.Equal(gotVal2, empty32) {
		t.Errorf("storage trie leakage between accounts! addr2 got %x", gotVal2)
	}

	// Test 4: SetState for addr2 and verify independent value
	err = sdb.SetState(addr2, slotKey, value2)
	if err != nil {
		t.Fatalf("SetState(addr2) failed: %v", err)
	}

	gotVal2Updated := sdb.GetState(addr2, slotKey)
	if !bytes.Equal(gotVal2Updated, value2) {
		t.Errorf("GetState(addr2) = %s, want %s", string(gotVal2Updated), string(value2))
	}
}
