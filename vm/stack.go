package vm

import (
	"errors"
	"math/big"
)

const StackLimit = 1024

var (
	ErrStackUnderflow = errors.New("evm stack underflow")
	ErrStackOverflow  = errors.New("evm stack overflow")
	ErrInvalidWord    = errors.New("evm stack word must be an unsigned 256-bit integer")
)

type Stack struct {
	items []*big.Int
}

func NewStack() *Stack {
	return &Stack{
		items: make([]*big.Int, 0, StackLimit),
	}
}

func (s *Stack) Push(value *big.Int) error {
	if len(s.items) >= StackLimit {
		return ErrStackOverflow
	}

	if value == nil || value.Sign() < 0 || value.BitLen() > 256 {
		return ErrInvalidWord
	}

	// Keep an independent copy so outside code cannot mutate stack data.
	s.items = append(s.items, new(big.Int).Set(value))
	return nil
}

func (s *Stack) Pop() (*big.Int, error) {
	if len(s.items) == 0 {
		return nil, ErrStackUnderflow
	}

	last := len(s.items) - 1
	value := s.items[last]
	s.items[last] = nil
	s.items = s.items[:last]

	return new(big.Int).Set(value), nil
}

// Peek returns the item at depth n without removing it.
// Peek(0) is the top item.
func (s *Stack) Peek(n int) (*big.Int, error) {
	index := len(s.items) - 1 - n
	if n < 0 || index < 0 {
		return nil, ErrStackUnderflow
	}

	return new(big.Int).Set(s.items[index]), nil
}

func (s *Stack) Len() int {
	return len(s.items)
}
