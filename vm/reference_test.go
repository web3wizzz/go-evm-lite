package vm

import (
	"math/big"
	"strings"
	"testing"
)

func TestReferenceArithmeticSemantics(t *testing.T) {
	tests := []struct {
		name     string
		code     []byte
		expected *big.Int
	}{
		{
			name: "SUB uses top minus next stack value",
			code: []byte{
				PUSH1, 0x05,
				PUSH1, 0x03,
				SUB,
				STOP,
			},
			expected: new(big.Int).Sub(wordModulus, big.NewInt(2)),
		},
		{
			name: "DIV uses top divided by next stack value",
			code: []byte{
				PUSH1, 0x06,
				PUSH1, 0x03,
				DIV,
				STOP,
			},
			expected: big.NewInt(0),
		},
		{
			name: "MOD uses top modulo next stack value",
			code: []byte{
				PUSH1, 0x06,
				PUSH1, 0x03,
				MOD,
				STOP,
			},
			expected: big.NewInt(3),
		},
		{
			name: "LT compares top against next stack value",
			code: []byte{
				PUSH1, 0x05,
				PUSH1, 0x03,
				LT,
				STOP,
			},
			expected: big.NewInt(1),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evm := NewEVM(1_000_000)

			if err := evm.Run(test.code); err != nil {
				t.Fatal(err)
			}

			got, err := evm.Context.Stack.Pop()
			if err != nil {
				t.Fatal(err)
			}

			if got.Cmp(test.expected) != 0 {
				t.Fatalf("result = %x, want %x", got, test.expected)
			}
		})
	}
}

func TestReferenceJumpDestinationRules(t *testing.T) {
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

func TestReferenceStoragePersistsPerContract(t *testing.T) {
	stateDB := NewEVM(1_000_000).State
	contract := [20]byte{0xca, 0xfe}

	writer := NewEVMWithState(1_000_000, stateDB, contract)
	reader := NewEVMWithState(1_000_000, stateDB, contract)

	if err := writer.Run([]byte{
		PUSH1, 0x2a,
		PUSH1, 0x01,
		SSTORE,
		STOP,
	}); err != nil {
		t.Fatal(err)
	}

	if err := reader.Run([]byte{
		PUSH1, 0x01,
		SLOAD,
		STOP,
	}); err != nil {
		t.Fatal(err)
	}

	got, err := reader.Context.Stack.Pop()
	if err != nil {
		t.Fatal(err)
	}

	if got.Uint64() != 42 {
		t.Fatalf("storage value = %d, want 42", got.Uint64())
	}
}
