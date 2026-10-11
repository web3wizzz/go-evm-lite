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

func TestEVMPop(t *testing.T) {
	evm := NewEVM(100)

	code := []byte{
		PUSH1, 0x2a,
		POP,
		STOP,
	}

	if err := evm.Run(code); err != nil {
		t.Fatal(err)
	}

	if evm.Context.Stack.Len() != 0 {
		t.Fatalf("expected empty stack, got %d item(s)", evm.Context.Stack.Len())
	}
}

func TestEVMPopUnderflow(t *testing.T) {
	evm := NewEVM(100)

	if err := evm.Run([]byte{POP}); err == nil {
		t.Fatal("expected POP on empty stack to fail")
	}
}

func TestEVMDup1(t *testing.T) {
	evm := NewEVM(100)

	code := []byte{
		PUSH1, 0x2a,
		DUP1,
		STOP,
	}

	if err := evm.Run(code); err != nil {
		t.Fatal(err)
	}

	top, err := evm.Context.Stack.Pop()
	if err != nil {
		t.Fatal(err)
	}

	below, err := evm.Context.Stack.Pop()
	if err != nil {
		t.Fatal(err)
	}

	if top.Cmp(big.NewInt(42)) != 0 || below.Cmp(big.NewInt(42)) != 0 {
		t.Fatalf("expected two copies of 42, got %s and %s", top, below)
	}
}

func TestEVMDup2(t *testing.T) {
	evm := NewEVM(100)

	code := []byte{
		PUSH1, 0x11,
		PUSH1, 0x22,
		DUP1 + 1,
		STOP,
	}

	if err := evm.Run(code); err != nil {
		t.Fatal(err)
	}

	first, _ := evm.Context.Stack.Pop()
	second, _ := evm.Context.Stack.Pop()
	third, _ := evm.Context.Stack.Pop()

	if first.Cmp(big.NewInt(0x11)) != 0 ||
		second.Cmp(big.NewInt(0x22)) != 0 ||
		third.Cmp(big.NewInt(0x11)) != 0 {
		t.Fatal("DUP2 did not duplicate the second stack item")
	}
}

func TestEVMDupUnderflow(t *testing.T) {
	evm := NewEVM(100)

	if err := evm.Run([]byte{DUP1}); err == nil {
		t.Fatal("expected DUP1 on an empty stack to fail")
	}
}
func TestEVMSwap1(t *testing.T) {
	evm := NewEVM(100)

	code := []byte{
		PUSH1, 0x11,
		PUSH1, 0x22,
		SWAP1,
		STOP,
	}

	if err := evm.Run(code); err != nil {
		t.Fatal(err)
	}

	top, _ := evm.Context.Stack.Pop()
	below, _ := evm.Context.Stack.Pop()

	if top.Cmp(big.NewInt(0x11)) != 0 ||
		below.Cmp(big.NewInt(0x22)) != 0 {
		t.Fatal("SWAP1 did not exchange the top two items")
	}
}

func TestEVMSwap2(t *testing.T) {
	evm := NewEVM(100)

	code := []byte{
		PUSH1, 0x11,
		PUSH1, 0x22,
		PUSH1, 0x33,
		SWAP1 + 1,
		STOP,
	}

	if err := evm.Run(code); err != nil {
		t.Fatal(err)
	}

	first, _ := evm.Context.Stack.Pop()
	second, _ := evm.Context.Stack.Pop()
	third, _ := evm.Context.Stack.Pop()

	if first.Cmp(big.NewInt(0x11)) != 0 ||
		second.Cmp(big.NewInt(0x22)) != 0 ||
		third.Cmp(big.NewInt(0x33)) != 0 {
		t.Fatal("SWAP2 did not exchange the correct items")
	}
}

func TestEVMSwapUnderflow(t *testing.T) {
	evm := NewEVM(100)

	if err := evm.Run([]byte{PUSH1, 0x01, SWAP1}); err == nil {
		t.Fatal("expected SWAP1 with one stack item to fail")
	}
}