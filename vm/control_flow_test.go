package vm

import (
	"strings"
	"testing"
)

func TestEVMJump(t *testing.T) {
	evm := NewEVM(1_000_000)

	code := []byte{
		PUSH1, 0x04,
		JUMP,
		STOP,
		JUMPDEST,
		PUSH1, 0x2a,
		STOP,
	}

	if err := evm.Run(code); err != nil {
		t.Fatal(err)
	}

	got, err := evm.Context.Stack.Pop()
	if err != nil {
		t.Fatal(err)
	}

	if got.Uint64() != 42 {
		t.Fatalf("result = %d, want 42", got.Uint64())
	}
}

func TestEVMJumpi(t *testing.T) {
	t.Run("true condition jumps", func(t *testing.T) {
		evm := NewEVM(1_000_000)

		code := []byte{
			PUSH1, 0x01,
			PUSH1, 0x06,
			JUMPI,
			STOP,
			JUMPDEST,
			PUSH1, 0x2a,
			STOP,
		}

		if err := evm.Run(code); err != nil {
			t.Fatal(err)
		}

		got, err := evm.Context.Stack.Pop()
		if err != nil {
			t.Fatal(err)
		}

		if got.Uint64() != 42 {
			t.Fatalf("result = %d, want 42", got.Uint64())
		}
	})

	t.Run("false condition does not jump", func(t *testing.T) {
		evm := NewEVM(1_000_000)

		code := []byte{
			PUSH1, 0x00,
			PUSH1, 0x08,
			JUMPI,
			PUSH1, 0x2a,
			STOP,
			JUMPDEST,
			PUSH1, 0xff,
			STOP,
		}

		if err := evm.Run(code); err != nil {
			t.Fatal(err)
		}

		got, err := evm.Context.Stack.Pop()
		if err != nil {
			t.Fatal(err)
		}

		if got.Uint64() != 42 {
			t.Fatalf("result = %d, want 42", got.Uint64())
		}
	})
}

func TestEVMRejectsInvalidJumpDestination(t *testing.T) {
	evm := NewEVM(1_000_000)

	err := evm.Run([]byte{
		PUSH1, 0x03,
		JUMP,
		STOP,
	})

	if err == nil {
		t.Fatal("expected invalid jump destination error")
	}

	if !strings.Contains(err.Error(), "invalid jump destination") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEVMRejectsJumpIntoPushData(t *testing.T) {
	evm := NewEVM(1_000_000)

	err := evm.Run([]byte{
		PUSH1, 0x01,
		JUMP,
		STOP,
	})

	if err == nil {
		t.Fatal("expected invalid jump destination error")
	}

	if !strings.Contains(err.Error(), "invalid jump destination") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEVMPC(t *testing.T) {
	evm := NewEVM(1_000_000)

	if err := evm.Run([]byte{PC, STOP}); err != nil {
		t.Fatal(err)
	}

	got, err := evm.Context.Stack.Pop()
	if err != nil {
		t.Fatal(err)
	}

	if got.Uint64() != 0 {
		t.Fatalf("PC = %d, want 0", got.Uint64())
	}
}
