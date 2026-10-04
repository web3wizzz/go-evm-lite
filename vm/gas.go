package vm

// This file implements the Yellow Paper's quadratic memory expansion cost
// (Appendix G): for memory covering `a` 32-byte words, the total cost is
//
//	C_mem(a) = 3*a + floor(a^2 / 512)
//
// Opcodes never pay this total directly — they pay the *marginal* cost of
// growing from the current word count to a new, larger word count. That
// marginal cost is what MemoryExpansionGasCost computes.

// memoryWordCount returns the number of 32-byte words needed to cover
// `size` bytes (i.e. ceil(size / 32)).
func memoryWordCount(size uint64) uint64 {
	return (size + wordSize - 1) / wordSize
}

// memoryCost returns the total (non-marginal) gas cost of memory that is
// `words` words long, per C_mem(a) = 3a + floor(a^2/512).
//
// Note: words^2 can overflow uint64 for sufficiently adversarial inputs
// (e.g. an attacker-controlled offset near 2^64). Real implementations cap
// memory size well below that point (geth effectively treats anything over
// ~2^32 words as "out of gas" immediately); callers of this package should
// impose a similar sanity cap on offsets/sizes before reaching this
// function. This implementation does not itself re-derive that cap.
func memoryCost(words uint64) uint64 {
	return 3*words + (words*words)/512
}

// MemoryExpansionGasCost returns the marginal gas cost of growing memory
// from covering `currentSize` bytes to covering `newSize` bytes (i.e. the
// highest offset an operation is about to touch, exclusive). If newSize
// does not require any additional 32-byte words beyond what currentSize
// already covers, the cost is 0 — this is why repeated MLOADs/MSTOREs
// within already-allocated memory are "free" (beyond their opcode's own
// flat gas cost).
//
// Callers must call this BEFORE performing the memory operation (and before
// calling Memory.Resize/Set/Set8), using the memory's size at that moment.
func MemoryExpansionGasCost(currentSize, newSize uint64) uint64 {
	if newSize <= currentSize {
		return 0
	}
	currentWords := memoryWordCount(currentSize)
	newWords := memoryWordCount(newSize)
	if newWords <= currentWords {
		return 0
	}
	return memoryCost(newWords) - memoryCost(currentWords)
}