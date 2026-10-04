package vm

import (
	"errors"
	"math/big"
	"testing"
)

func TestStackPushPopLIFO(t *testing.T) {
	stack := NewStack()

	if err := stack.Push(big.NewInt(5)); err != nil {
		t.Fatal(err)
	}
	if err := stack.Push(big.NewInt(10)); err != nil {
		t.Fatal(err)
	}

	value, err := stack.Pop()
	if err != nil {
		t.Fatal(err)
	}
	if value.Cmp(big.NewInt(10)) != 0 {
		t.Fatalf("expected 10, got %s", value)
	}

	value, err = stack.Pop()
	if err != nil {
		t.Fatal(err)
	}
	if value.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("expected 5, got %s", value)
	}
}

func TestStackUnderflow(t *testing.T) {
	_, err := NewStack().Pop()

	if !errors.Is(err, ErrStackUnderflow) {
		t.Fatalf("expected stack underflow, got %v", err)
	}
}

func TestStackRejectsInvalidWords(t *testing.T) {
	stack := NewStack()

	if err := stack.Push(big.NewInt(-1)); !errors.Is(err, ErrInvalidWord) {
		t.Fatalf("expected invalid negative word error, got %v", err)
	}

	tooLarge := new(big.Int).Lsh(big.NewInt(1), 256)
	if err := stack.Push(tooLarge); !errors.Is(err, ErrInvalidWord) {
		t.Fatalf("expected invalid 257-bit word error, got %v", err)
	}
}

func TestStackOverflow(t *testing.T) {
	stack := NewStack()

	for i := 0; i < StackLimit; i++ {
		if err := stack.Push(big.NewInt(int64(i))); err != nil {
			t.Fatalf("push %d failed: %v", i, err)
		}
	}

	if err := stack.Push(big.NewInt(1)); !errors.Is(err, ErrStackOverflow) {
		t.Fatalf("expected stack overflow, got %v", err)
	}
}
