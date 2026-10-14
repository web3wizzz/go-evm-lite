package vm

import (
	"math/big"
	"testing"
)

func TestEVMBasicBitwiseOps(t *testing.T) {
	tests := []struct {
		name     string
		opcode   byte
		expected int64
	}{
		{"AND", AND, 0x08},
		{"OR", OR, 0x0e},
		{"XOR", XOR, 0x06},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evm := NewEVM(100)

			code := []byte{
				PUSH1, 0x0c,
				PUSH1, 0x0a,
				test.opcode,
				STOP,
			}

			if err := evm.Run(code); err != nil {
				t.Fatal(err)
			}

			value, err := evm.Context.Stack.Pop()
			if err != nil {
				t.Fatal(err)
			}

			if value.Cmp(big.NewInt(test.expected)) != 0 {
				t.Fatalf("expected %d, got %s", test.expected, value)
			}
		})
	}
}

func TestEVMNot(t *testing.T) {
	evm := NewEVM(100)

	if err := evm.Run([]byte{PUSH1, 0x00, NOT, STOP}); err != nil {
		t.Fatal(err)
	}

	value, err := evm.Context.Stack.Pop()
	if err != nil {
		t.Fatal(err)
	}

	if value.Cmp(maxWord) != 0 {
		t.Fatalf("expected max 256-bit word, got %s", value)
	}
}

func TestEVMByte(t *testing.T) {
	evm := NewEVM(100)

	word := make([]byte, 32)
	word[0] = 0xab

	code := append([]byte{PUSH32}, word...)
	code = append(code, PUSH1, 0x00, BYTE, STOP)

	if err := evm.Run(code); err != nil {
		t.Fatal(err)
	}

	value, err := evm.Context.Stack.Pop()
	if err != nil {
		t.Fatal(err)
	}

	if value.Cmp(big.NewInt(0xab)) != 0 {
		t.Fatalf("expected 0xab, got %s", value)
	}
}

func TestEVMShifts(t *testing.T) {
	tests := []struct {
		name     string
		code     []byte
		expected int64
	}{
		{
			"SHL",
			[]byte{PUSH1, 0x03, PUSH1, 0x02, SHL, STOP},
			12,
		},
		{
			"SHR",
			[]byte{PUSH1, 0x10, PUSH1, 0x02, SHR, STOP},
			4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evm := NewEVM(100)

			if err := evm.Run(test.code); err != nil {
				t.Fatal(err)
			}

			value, err := evm.Context.Stack.Pop()
			if err != nil {
				t.Fatal(err)
			}

			if value.Cmp(big.NewInt(test.expected)) != 0 {
				t.Fatalf("expected %d, got %s", test.expected, value)
			}
		})
	}
}

func TestEVMSarNegativeOne(t *testing.T) {
	evm := NewEVM(100)

	negativeOne := make([]byte, 32)
	for i := range negativeOne {
		negativeOne[i] = 0xff
	}

	code := append([]byte{PUSH32}, negativeOne...)
	code = append(code, PUSH1, 0x01, SAR, STOP)

	if err := evm.Run(code); err != nil {
		t.Fatal(err)
	}

	value, err := evm.Context.Stack.Pop()
	if err != nil {
		t.Fatal(err)
	}

	if value.Cmp(maxWord) != 0 {
		t.Fatalf("expected negative one after SAR, got %s", value)
	}
}
