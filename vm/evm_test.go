package vm

import (
	"math/big"
	"testing"
)

func TestEVMStop(t *testing.T) {
	evm := NewEVM(100)

	if err := evm.Run([]byte{STOP}); err != nil {
		t.Fatalf("expected STOP to succeed, got %v", err)
	}

	if evm.Context.PC != 1 {
		t.Fatalf("expected PC 1, got %d", evm.Context.PC)
	}
}

func TestEVMEmptyCodeStops(t *testing.T) {
	evm := NewEVM(100)

	if err := evm.Run(nil); err != nil {
		t.Fatalf("expected empty code to stop, got %v", err)
	}
}

func TestEVMRejectsInvalidOpcode(t *testing.T) {
	evm := NewEVM(100)

	if err := evm.Run([]byte{0xfe}); err == nil {
		t.Fatal("expected invalid opcode error")
	}
}

func TestEVMPush1(t *testing.T) {
	evm := NewEVM(100)

	if err := evm.Run([]byte{PUSH1, 0x2a, STOP}); err != nil {
		t.Fatal(err)
	}

	value, err := evm.Context.Stack.Pop()
	if err != nil {
		t.Fatal(err)
	}

	if value.Cmp(big.NewInt(42)) != 0 {
		t.Fatalf("expected 42, got %s", value)
	}
}

func TestEVMPush32(t *testing.T) {
	evm := NewEVM(100)

	data := make([]byte, 32)
	for i := range data {
		data[i] = byte(i + 1)
	}

	code := append([]byte{PUSH32}, data...)
	code = append(code, STOP)

	if err := evm.Run(code); err != nil {
		t.Fatal(err)
	}

	value, err := evm.Context.Stack.Pop()
	if err != nil {
		t.Fatal(err)
	}

	expected := new(big.Int).SetBytes(data)
	if value.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, value)
	}
}

func TestEVMPushPadsMissingBytesWithZero(t *testing.T) {
	evm := NewEVM(100)

	// PUSH2 has only one provided byte, so 0x01 becomes 0x0100.
	if err := evm.Run([]byte{PUSH1 + 1, 0x01}); err != nil {
		t.Fatal(err)
	}

	value, err := evm.Context.Stack.Pop()
	if err != nil {
		t.Fatal(err)
	}

	if value.Cmp(big.NewInt(0x100)) != 0 {
		t.Fatalf("expected 256, got %s", value)
	}
}