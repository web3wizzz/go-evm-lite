package vm

import (
	"testing"

	"evm-lite/state"
)

func TestEVMSStorePersistsAcrossRuns(t *testing.T) {
	stateDB := state.NewStateDB()
	contract := [20]byte{0xca, 0xfe, 0x01}

	writer := NewEVMWithState(1_000_000, stateDB, contract)

	writeCode := NewAssembler().
		Push1(0x2a). // value 42
		Push1(0x01). // storage slot 1
		Op(SSTORE).
		Op(STOP).
		Bytes()

	if err := writer.Run(writeCode); err != nil {
		t.Fatalf("write Run returned error: %v", err)
	}

	reader := NewEVMWithState(1_000_000, stateDB, contract)

	readCode := NewAssembler().
		Push1(0x01).
		Op(SLOAD).
		Op(STOP).
		Bytes()

	if err := reader.Run(readCode); err != nil {
		t.Fatalf("read Run returned error: %v", err)
	}

	got, err := reader.Context.Stack.Pop()
	if err != nil {
		t.Fatalf("Pop returned error: %v", err)
	}

	if got.Uint64() != 42 {
		t.Fatalf("storage value = %d, want 42", got.Uint64())
	}
}

func TestEVMSLoadMissingSlotReturnsZero(t *testing.T) {
	evm := NewEVMWithState(
		1_000_000,
		state.NewStateDB(),
		[20]byte{0xca, 0xfe, 0x02},
	)

	code := NewAssembler().
		Push1(0x09).
		Op(SLOAD).
		Op(STOP).
		Bytes()

	if err := evm.Run(code); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	got, err := evm.Context.Stack.Pop()
	if err != nil {
		t.Fatalf("Pop returned error: %v", err)
	}

	if got.Sign() != 0 {
		t.Fatalf("missing storage slot = %d, want 0", got.Uint64())
	}
}

func TestEVMStorageIsIsolatedByContract(t *testing.T) {
	stateDB := state.NewStateDB()

	contractA := [20]byte{0x01}
	contractB := [20]byte{0x02}

	writer := NewEVMWithState(1_000_000, stateDB, contractA)

	if err := writer.Run(NewAssembler().
		Push1(0x2a).
		Push1(0x01).
		Op(SSTORE).
		Op(STOP).
		Bytes()); err != nil {
		t.Fatalf("write Run returned error: %v", err)
	}

	reader := NewEVMWithState(1_000_000, stateDB, contractB)

	if err := reader.Run(NewAssembler().
		Push1(0x01).
		Op(SLOAD).
		Op(STOP).
		Bytes()); err != nil {
		t.Fatalf("read Run returned error: %v", err)
	}

	got, err := reader.Context.Stack.Pop()
	if err != nil {
		t.Fatalf("Pop returned error: %v", err)
	}

	if got.Sign() != 0 {
		t.Fatalf("other contract storage = %d, want 0", got.Uint64())
	}
}
