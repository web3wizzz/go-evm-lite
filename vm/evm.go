package vm

import (
	"fmt"
	"math/big"
)

const (
	STOP   byte = 0x00
	POP    byte = 0x50
	PUSH1  byte = 0x60
	PUSH32 byte = 0x7f
	DUP1   byte = 0x80
	DUP16  byte = 0x8f
	SWAP1  byte = 0x90
	SWAP16 byte = 0x9f
	ADD    byte = 0x01
	MUL    byte = 0x02
	SUB    byte = 0x03
	DIV    byte = 0x04
	MOD    byte = 0x06
)

type EVM struct {
	Context *ExecutionContext
	Code    []byte
}

func NewEVM(gasLimit uint64) *EVM {
	return &EVM{
		Context: NewExecutionContext(gasLimit),
	}
}

var wordModulus = new(big.Int).Lsh(big.NewInt(1), 256)

func (e *EVM) binaryOp(
	name string,
	operation func(a, b *big.Int) *big.Int,
) error {
	a, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: %s failed: %w", name, err)
	}

	b, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: %s failed: %w", name, err)
	}

	result := operation(a, b)
	result.Mod(result, wordModulus)

	if err := e.Context.Stack.Push(result); err != nil {
		return fmt.Errorf("vm: %s failed: %w", name, err)
	}

	return nil
}

// Run executes bytecode until STOP, an invalid opcode, or the end of the code.
func (e *EVM) Run(code []byte) error {
	e.Code = append([]byte(nil), code...)
	e.Context.PC = 0

	for e.Context.PC < uint64(len(e.Code)) {
		opcode := e.Code[e.Context.PC]
		e.Context.PC++

		switch {
		case opcode == STOP:
			return nil

		case opcode == ADD:
			if err := e.binaryOp("ADD", func(a, b *big.Int) *big.Int {
				return new(big.Int).Add(a, b)
			}); err != nil {
				return err
			}

		case opcode == SUB:
			if err := e.binaryOp("SUB", func(a, b *big.Int) *big.Int {
				return new(big.Int).Sub(a, b)
			}); err != nil {
				return err
			}

		case opcode == MUL:
			if err := e.binaryOp("MUL", func(a, b *big.Int) *big.Int {
				return new(big.Int).Mul(a, b)
			}); err != nil {
				return err
			}

		case opcode == DIV:
			if err := e.binaryOp("DIV", func(a, b *big.Int) *big.Int {
				if b.Sign() == 0 {
					return new(big.Int)
				}
				return new(big.Int).Div(a, b)
			}); err != nil {
				return err
			}

		case opcode == MOD:
			if err := e.binaryOp("MOD", func(a, b *big.Int) *big.Int {
				if b.Sign() == 0 {
					return new(big.Int)
				}
				return new(big.Int).Mod(a, b)
			}); err != nil {
				return err
			}

		case opcode == POP:
			if _, err := e.Context.Stack.Pop(); err != nil {
				return fmt.Errorf("vm: POP failed: %w", err)
			}

		case opcode >= PUSH1 && opcode <= PUSH32:
			size := uint64(opcode-PUSH1) + 1
			data := make([]byte, size)

			// Missing PUSH bytes are treated as zero.
			if e.Context.PC < uint64(len(e.Code)) {
				available := uint64(len(e.Code)) - e.Context.PC
				if available > size {
					available = size
				}

				copy(data, e.Code[e.Context.PC:e.Context.PC+available])
			}

			e.Context.PC += size

			if err := e.Context.Stack.Push(new(big.Int).SetBytes(data)); err != nil {
				return fmt.Errorf("vm: PUSH%d failed: %w", size, err)
			}

		case opcode >= DUP1 && opcode <= DUP16:
			depth := int(opcode-DUP1) + 1

			value, err := e.Context.Stack.Peek(depth - 1)
			if err != nil {
				return fmt.Errorf("vm: DUP%d failed: %w", depth, err)
			}

			if err := e.Context.Stack.Push(value); err != nil {
				return fmt.Errorf("vm: DUP%d failed: %w", depth, err)
			}
		case opcode >= SWAP1 && opcode <= SWAP16:
			depth := int(opcode-SWAP1) + 1

			if err := e.Context.Stack.Swap(depth); err != nil {
				return fmt.Errorf("vm: SWAP%d failed: %w", depth, err)
			}

		default:
			return fmt.Errorf(
				"vm: invalid opcode 0x%02x at pc %d",
				opcode,
				e.Context.PC-1,
			)
		}
	}

	return nil
}
