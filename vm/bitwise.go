package vm

import (
	"fmt"
	"math/big"
)

var maxWord = new(big.Int).Sub(wordModulus, big.NewInt(1))

func (e *EVM) bitwiseNot() error {
	value, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: NOT failed: %w", err)
	}

	result := new(big.Int).Xor(value, maxWord)

	if err := e.Context.Stack.Push(result); err != nil {
		return fmt.Errorf("vm: NOT failed: %w", err)
	}

	return nil
}

func (e *EVM) byteOp() error {
	index, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: BYTE failed: %w", err)
	}

	value, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: BYTE failed: %w", err)
	}

	result := new(big.Int)

	if index.Cmp(big.NewInt(32)) < 0 {
		shift := uint((31 - index.Uint64()) * 8)
		result.Rsh(value, shift)
		result.And(result, big.NewInt(0xff))
	}

	if err := e.Context.Stack.Push(result); err != nil {
		return fmt.Errorf("vm: BYTE failed: %w", err)
	}

	return nil
}

func (e *EVM) shl() error {
	shift, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: SHL failed: %w", err)
	}

	value, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: SHL failed: %w", err)
	}

	result := new(big.Int)

	if shift.Cmp(big.NewInt(256)) < 0 {
		result.Lsh(value, uint(shift.Uint64()))
		result.Mod(result, wordModulus)
	}

	if err := e.Context.Stack.Push(result); err != nil {
		return fmt.Errorf("vm: SHL failed: %w", err)
	}

	return nil
}

func (e *EVM) shr() error {
	shift, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: SHR failed: %w", err)
	}

	value, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: SHR failed: %w", err)
	}

	result := new(big.Int)

	if shift.Cmp(big.NewInt(256)) < 0 {
		result.Rsh(value, uint(shift.Uint64()))
	}

	if err := e.Context.Stack.Push(result); err != nil {
		return fmt.Errorf("vm: SHR failed: %w", err)
	}

	return nil
}

func (e *EVM) sar() error {
	shift, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: SAR failed: %w", err)
	}

	value, err := e.Context.Stack.Pop()
	if err != nil {
		return fmt.Errorf("vm: SAR failed: %w", err)
	}

	signedValue := signedWord(value)
	result := new(big.Int)

	if shift.Cmp(big.NewInt(256)) >= 0 {
		if signedValue.Sign() < 0 {
			result.Set(maxWord)
		}
	} else {
		result.Rsh(signedValue, uint(shift.Uint64()))
		result.Mod(result, wordModulus)
	}

	if err := e.Context.Stack.Push(result); err != nil {
		return fmt.Errorf("vm: SAR failed: %w", err)
	}

	return nil
}
