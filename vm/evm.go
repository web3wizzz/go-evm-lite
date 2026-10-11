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