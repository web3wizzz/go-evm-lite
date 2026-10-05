package vm

import "fmt"

const STOP byte = 0x00

type EVM struct {
	Context *ExecutionContext
	Code    []byte
}

func NewEVM(gasLimit uint64) *EVM {
	return &EVM{
		Context: NewExecutionContext(gasLimit),
	}
}

func (e *EVM) Run(code []byte) error {
	e.Code = append([]byte(nil), code...)
	e.Context.PC = 0

	for e.Context.PC < uint64(len(e.Code)) {
		opcode := e.Code[e.Context.PC]
		e.Context.PC++

		switch opcode {
		case STOP:
			return nil
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
