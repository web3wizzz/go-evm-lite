package vm

import "testing"

func TestMemoryWordCount(t *testing.T) {
	cases := []struct {
		size  uint64
		words uint64
	}{
		{0, 0},
		{1, 1},
		{32, 1},
		{33, 2},
		{64, 2},
		{65, 3},
	}
	for _, c := range cases {
		got := memoryWordCount(c.size)
		if got != c.words {
			t.Fatalf("memoryWordCount(%d): expected %d, got %d", c.size, c.words, got)
		}
	}
}

func TestMemoryCost_KnownValues(t *testing.T) {
	// Hand-computed from C_mem(a) = 3a + floor(a^2/512).
	cases := []struct {
		words uint64
		cost  uint64
	}{
		{0, 0},
		{1, 3},                       // 3*1 + 0
		{32, 3*32 + (32*32)/512},     // 96 + 2 = 98
		{100, 3*100 + (100*100)/512}, // 300 + 19 = 319
	}
	for _, c := range cases {
		got := memoryCost(c.words)
		if got != c.cost {
			t.Fatalf("memoryCost(%d): expected %d, got %d", c.words, c.cost, got)
		}
	}
}

func TestMemoryExpansionGasCost_NoExpansionIsFree(t *testing.T) {
	// Growing from 0 to 0, or requesting something already covered,
	// costs nothing.
	if got := MemoryExpansionGasCost(64, 32); got != 0 {
		t.Fatalf("expected 0 for non-expanding request, got %d", got)
	}
	if got := MemoryExpansionGasCost(64, 64); got != 0 {
		t.Fatalf("expected 0 for same-size request, got %d", got)
	}
}

func TestMemoryExpansionGasCost_FirstWord(t *testing.T) {
	// From empty memory, touching the first word (32 bytes) costs
	// memoryCost(1) - memoryCost(0) = 3 - 0 = 3.
	got := MemoryExpansionGasCost(0, 32)
	want := uint64(3)
	if got != want {
		t.Fatalf("expected %d, got %d", want, got)
	}
}

func TestMemoryExpansionGasCost_MarginalNotTotal(t *testing.T) {
	// Expanding from 1 word to 32 words should cost memoryCost(32) -
	// memoryCost(1), NOT memoryCost(32) outright.
	got := MemoryExpansionGasCost(32, 32*32) // currentSize=32 (1 word), newSize=1024 (32 words)
	want := memoryCost(32) - memoryCost(1)
	if got != want {
		t.Fatalf("expected %d, got %d", want, got)
	}

	// Sanity: this must be strictly less than the full, non-marginal cost
	// of 32 words, since 1 word's cost was already "paid for".
	if got >= memoryCost(32) {
		t.Fatalf("marginal cost (%d) should be less than total cost (%d)", got, memoryCost(32))
	}
}

func TestMemoryExpansionGasCost_PartialWordRoundsUp(t *testing.T) {
	// Requesting newSize=33 still only needs 2 words (ceil(33/32)=2), same
	// as requesting newSize=64 -- the cost should be identical for both,
	// since both land in the same word count.
	costAt33 := MemoryExpansionGasCost(0, 33)
	costAt64 := MemoryExpansionGasCost(0, 64)
	if costAt33 != costAt64 {
		t.Fatalf("expected equal cost for 33 and 64 (both round to 2 words): got %d vs %d", costAt33, costAt64)
	}
}

func TestExecutionContext_MStoreChargesGasOnce(t *testing.T) {
	ctx := NewExecutionContext(1_000_000)
	startGas := ctx.GasRemaining

	word := make([]byte, 32)
	word[31] = 0xff

	if err := ctx.MStore(0, word); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	firstCost := startGas - ctx.GasRemaining
	if firstCost != 3 { // memoryCost(1) - memoryCost(0)
		t.Fatalf("expected first MSTORE to cost 3 gas, cost %d", firstCost)
	}

	// A second MSTORE to the SAME word should cost 0 additional gas --
	// memory is already expanded to cover it.
	gasBeforeSecond := ctx.GasRemaining
	if err := ctx.MStore(0, word); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.GasRemaining != gasBeforeSecond {
		t.Fatalf("expected re-writing the same word to cost 0 gas, cost %d", gasBeforeSecond-ctx.GasRemaining)
	}
}

func TestExecutionContext_OutOfGas(t *testing.T) {
	ctx := NewExecutionContext(2) // not enough even for the first word (costs 3)
	word := make([]byte, 32)
	err := ctx.MStore(0, word)
	if err == nil {
		t.Fatalf("expected out-of-gas error, got nil")
	}
	if ctx.GasRemaining != 0 {
		t.Fatalf("expected GasRemaining to be clamped to 0, got %d", ctx.GasRemaining)
	}
}

func TestExecutionContext_MStore8ExpandsByFullWord(t *testing.T) {
	ctx := NewExecutionContext(1_000_000)
	if err := ctx.MStore8(0, 0x01); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.Memory.Len() != 32 {
		t.Fatalf("expected memory to expand to a full word (32), got %d", ctx.Memory.Len())
	}
}

func TestExecutionContext_MLoadRoundTrip(t *testing.T) {
	ctx := NewExecutionContext(1_000_000)
	word := make([]byte, 32)
	word[0] = 0xde
	word[31] = 0xef

	if err := ctx.MStore(64, word); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := ctx.MLoad(64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := range word {
		if got[i] != word[i] {
			t.Fatalf("mismatch at byte %d: expected %02x, got %02x", i, word[i], got[i])
		}
	}
}
