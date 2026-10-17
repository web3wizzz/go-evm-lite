package vm

import (
	"bytes"
	"testing"
)

func TestAssemblerBuildsBytecode(t *testing.T) {
	code := NewAssembler().
		Push1(0x05).
		Push1(0x03).
		Op(ADD).
		Op(STOP).
		Bytes()

	want := []byte{
		PUSH1, 0x05,
		PUSH1, 0x03,
		ADD,
		STOP,
	}

	if !bytes.Equal(code, want) {
		t.Fatalf("assembled code = %x, want %x", code, want)
	}
}

func TestEVMRunsAssembledArithmetic(t *testing.T) {
	code := NewAssembler().
		Push1(0x05).
		Push1(0x03).
		Op(ADD).
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

	if got.Uint64() != 8 {
		t.Fatalf("result = %d, want 8", got.Uint64())
	}
}
