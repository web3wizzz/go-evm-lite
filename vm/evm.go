package vm

import (
	"fmt"
	"math/big"
)

const (
	STOP   byte = 0x00
	PUSH1  byte = 0x60
	PUSH32 byte = 0x7f
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

		case opcode >= PUSH1 && opcode <= PUSH32:
			size := uint64(opcode-PUSH1) + 1
			data := make([]byte, size)

			// Copy immediate bytes; missing bytes remain zero.
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