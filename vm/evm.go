package vm

import (
	"fmt"
	"math/big"
)

const (
	STOP       byte = 0x00
	ADD        byte = 0x01
	MUL        byte = 0x02
	SUB        byte = 0x03
	DIV        byte = 0x04
	SDIV       byte = 0x05
	MOD        byte = 0x06
	SMOD       byte = 0x07
	EXP        byte = 0x0a
	SIGNEXTEND byte = 0x0b
	LT         byte = 0x10
	GT         byte = 0x11
	SLT        byte = 0x12
	SGT        byte = 0x13
	EQ         byte = 0x14
	ISZERO     byte = 0x15
	POP        byte = 0x50
	PUSH1      byte = 0x60
	PUSH32     byte = 0x7f
	DUP1       byte = 0x80
	DUP16      byte = 0x8f
	SWAP1      byte = 0x90
	SWAP16     byte = 0x9f
	AND        byte = 0x16
	OR         byte = 0x17
	XOR        byte = 0x18
	NOT        byte = 0x19
	BYTE       byte = 0x1a
	SHL        byte = 0x1b
	SHR        byte = 0x1c
	SAR        byte = 0x1d
	MLOAD      byte = 0x51
	MSTORE     byte = 0x52
	MSTORE8    byte = 0x53
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
var signedBoundary = new(big.Int).Lsh(big.NewInt(1), 255)

// signedWord converts an unsigned 256-bit EVM word to its signed form.
func signedWord(value *big.Int) *big.Int {
	if value.Cmp(signedBoundary) >= 0 {
		return new(big.Int).Sub(value, wordModulus)
	}

	return new(big.Int).Set(value)
}

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

func (e *EVM) binaryCompare(
	name string,
	comparison func(a, b *big.Int) bool,
) error {
	a, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: %s failed: %w", name, err)
	}

	b, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: %s failed: %w", name, err)
	}

	result := big.NewInt(0)
	if comparison(a, b) {
		result.SetInt64(1)
	}

	if err := e.Context.Stack.Push(result); err != nil {
		return fmt.Errorf("vm: %s failed: %w", name, err)
	}

	return nil
}

func (e *EVM) isZero() error {
	value, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: ISZERO failed: %w", err)
	}

	result := big.NewInt(0)
	if value.Sign() == 0 {
		result.SetInt64(1)
	}

	if err := e.Context.Stack.Push(result); err != nil {
		return fmt.Errorf("vm: ISZERO failed: %w", err)
	}

	return nil
}

func (e *EVM) signExtend() error {
	byteIndex, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: SIGNEXTEND failed: %w", err)
	}

	value, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: SIGNEXTEND failed: %w", err)
	}

	if byteIndex.Cmp(big.NewInt(32)) >= 0 {
		return e.Context.Stack.Push(value)
	}

	bits := uint(byteIndex.Uint64()*8 + 8)

	mask := new(big.Int).Sub(
		new(big.Int).Lsh(big.NewInt(1), bits),
		big.NewInt(1),
	)

	result := new(big.Int).And(value, mask)

	if result.Bit(int(bits-1)) == 1 {
		extension := new(big.Int).Xor(
			new(big.Int).Sub(wordModulus, big.NewInt(1)),
			mask,
		)
		result.Or(result, extension)
	}

	return e.Context.Stack.Push(result)
}

const memoryOpcodeGas uint64 = 3

func (e *EVM) popMemoryOffset(op string) (uint64, error) {
	value, err := e.Context.Stack.Pop()
	if err != nil {
		return 0, fmt.Errorf("vm: %s failed: %w", op, err)
	}

	if !value.IsUint64() {
		return 0, fmt.Errorf("vm: %s offset does not fit uint64", op)
	}

	return value.Uint64(), nil
}

func wordToBytes(value *big.Int) []byte {
	word := make([]byte, 32)
	valueBytes := value.Bytes()

	copy(word[32-len(valueBytes):], valueBytes)

	return word
}

func (e *EVM) mload() error {
	if err := e.Context.useGas(memoryOpcodeGas); err != nil {
		return fmt.Errorf("vm: MLOAD failed: %w", err)
	}

	offset, err := e.popMemoryOffset("MLOAD")
	if err != nil {
		return err
	}

	word, err := e.Context.MLoad(offset)
	if err != nil {
		return fmt.Errorf("vm: MLOAD failed: %w", err)
	}

	if err := e.Context.Stack.Push(new(big.Int).SetBytes(word)); err != nil {
		return fmt.Errorf("vm: MLOAD failed: %w", err)
	}

	return nil
}

func (e *EVM) mstore() error {
	if err := e.Context.useGas(memoryOpcodeGas); err != nil {
		return fmt.Errorf("vm: MSTORE failed: %w", err)
	}

	offset, err := e.popMemoryOffset("MSTORE")
	if err != nil {
		return err
	}

	value, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: MSTORE failed: %w", err)
	}

	if err := e.Context.MStore(offset, wordToBytes(value)); err != nil {
		return fmt.Errorf("vm: MSTORE failed: %w", err)
	}

	return nil
}

func (e *EVM) mstore8() error {
	if err := e.Context.useGas(memoryOpcodeGas); err != nil {
		return fmt.Errorf("vm: MSTORE8 failed: %w", err)
	}

	offset, err := e.popMemoryOffset("MSTORE8")
	if err != nil {
		return err
	}

	value, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: MSTORE8 failed: %w", err)
	}

	if err := e.Context.MStore8(offset, byte(value.Uint64())); err != nil {
		return fmt.Errorf("vm: MSTORE8 failed: %w", err)
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

		case opcode == SDIV:
			if err := e.binaryOp("SDIV", func(a, b *big.Int) *big.Int {
				dividend := signedWord(a)
				divisor := signedWord(b)

				if divisor.Sign() == 0 {
					return new(big.Int)
				}

				return new(big.Int).Quo(dividend, divisor)
			}); err != nil {
				return err
			}

		case opcode == SMOD:
			if err := e.binaryOp("SMOD", func(a, b *big.Int) *big.Int {
				dividend := signedWord(a)
				divisor := signedWord(b)

				if divisor.Sign() == 0 {
					return new(big.Int)
				}

				return new(big.Int).Rem(dividend, divisor)
			}); err != nil {
				return err
			}

		case opcode == EXP:
			if err := e.binaryOp("EXP", func(a, b *big.Int) *big.Int {
				return new(big.Int).Exp(a, b, wordModulus)
			}); err != nil {
				return err
			}

		case opcode == SIGNEXTEND:
			if err := e.signExtend(); err != nil {
				return err
			}

		case opcode == LT:
			if err := e.binaryCompare("LT", func(a, b *big.Int) bool {
				return a.Cmp(b) < 0
			}); err != nil {
				return err
			}

		case opcode == GT:
			if err := e.binaryCompare("GT", func(a, b *big.Int) bool {
				return a.Cmp(b) > 0
			}); err != nil {
				return err
			}

		case opcode == SLT:
			if err := e.binaryCompare("SLT", func(a, b *big.Int) bool {
				return signedWord(a).Cmp(signedWord(b)) < 0
			}); err != nil {
				return err
			}

		case opcode == SGT:
			if err := e.binaryCompare("SGT", func(a, b *big.Int) bool {
				return signedWord(a).Cmp(signedWord(b)) > 0
			}); err != nil {
				return err
			}

		case opcode == EQ:
			if err := e.binaryCompare("EQ", func(a, b *big.Int) bool {
				return a.Cmp(b) == 0
			}); err != nil {
				return err
			}

		case opcode == ISZERO:
			if err := e.isZero(); err != nil {
				return err
			}

		case opcode == AND:
			if err := e.binaryOp("AND", func(a, b *big.Int) *big.Int {
				return new(big.Int).And(a, b)
			}); err != nil {
				return err
			}

		case opcode == OR:
			if err := e.binaryOp("OR", func(a, b *big.Int) *big.Int {
				return new(big.Int).Or(a, b)
			}); err != nil {
				return err
			}

		case opcode == MLOAD:
			if err := e.mload(); err != nil {
				return err
			}

		case opcode == MSTORE:
			if err := e.mstore(); err != nil {
				return err
			}

		case opcode == MSTORE8:
			if err := e.mstore8(); err != nil {
				return err
			}

		case opcode == XOR:
			if err := e.binaryOp("XOR", func(a, b *big.Int) *big.Int {
				return new(big.Int).Xor(a, b)
			}); err != nil {
				return err
			}

		case opcode == NOT:
			if err := e.bitwiseNot(); err != nil {
				return err
			}

		case opcode == BYTE:
			if err := e.byteOp(); err != nil {
				return err
			}

		case opcode == SHL:
			if err := e.shl(); err != nil {
				return err
			}

		case opcode == SHR:
			if err := e.shr(); err != nil {
				return err
			}

		case opcode == SAR:
			if err := e.sar(); err != nil {
				return err
			}

		case opcode == POP:
			if _, err := e.Context.Stack.Pop(); err != nil {
				return fmt.Errorf("vm: POP failed: %w", err)
			}

		case opcode >= PUSH1 && opcode <= PUSH32:
			size := uint64(opcode-PUSH1) + 1
			data := make([]byte, size)

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
