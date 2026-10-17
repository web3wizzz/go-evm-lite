package vm

import "testing"

func TestEVMMStoreAndMLoad(t *testing.T) {
	code := NewAssembler().
		Push1(0x2a). // value: 42
		Push1(0x00). // memory offset
		Op(MSTORE).
		Push1(0x00).
		Op(MLOAD).
		Op(STOP).
		Bytes()

	evm := NewEVM(1_000_000)

	if err := evm.Run(code); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	got, err := evm.Context.Stack.Pop()
	if err != nil {
		t.Fatalf("Pop returned error: %v", err)
	}

	if got.Uint64() != 42 {
		t.Fatalf("result = %d, want 42", got.Uint64())
	}
}

func TestEVMMStore8(t *testing.T) {
	code := NewAssembler().
		Push1(0xab). // byte to store
		Push1(0x1f). // offset 31: last byte of a 32-byte word
		Op(MSTORE8).
		Push1(0x00).
		Op(MLOAD).
		Op(STOP).
		Bytes()

	evm := NewEVM(1_000_000)

	if err := evm.Run(code); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	got, err := evm.Context.Stack.Pop()
	if err != nil {
		t.Fatalf("Pop returned error: %v", err)
	}

	if got.Uint64() != 0xab {
		t.Fatalf("result = 0x%x, want 0xab", got.Uint64())
	}
}
