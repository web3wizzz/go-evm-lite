package vm

import "fmt"

// wordSize is the EVM's word size in bytes. Memory always expands in
// multiples of this size, even when an individual write (e.g. MSTORE8)
// touches only a single byte.
const wordSize = 32

// Memory models the EVM's byte-addressable, linearly expanding memory.
// It is conceptually infinite but allocated lazily: it only ever grows to
// cover the highest offset actually touched, and never shrinks.
//
// Memory itself does NOT charge gas. Gas accounting for expansion is the
// caller's responsibility (see MemoryExpansionGasCost in gas.go) — callers
// must compute and charge the cost *before* calling Resize/Set/Set8, using
// the memory's size prior to the call.
type Memory struct {
	store []byte
}

// NewMemory returns an empty Memory.
func NewMemory() *Memory {
	return &Memory{store: []byte{}}
}

// Len returns the current size of memory in bytes (always a multiple of 32).
func (m *Memory) Len() uint64 {
	return uint64(len(m.store))
}

// Resize grows memory so it covers at least `size` bytes, zero-filling the
// new region. The actual allocation is rounded up to the next 32-byte word
// boundary, matching the EVM's word-based expansion. Resize never shrinks
// memory — shrinking is not part of the EVM memory model.
func (m *Memory) Resize(size uint64) {
	aligned := wordAlign(size)
	if aligned <= uint64(len(m.store)) {
		return
	}
	grown := make([]byte, aligned)
	copy(grown, m.store)
	m.store = grown
}

// Set writes `value` into memory at [offset, offset+size), resizing memory
// first if needed. If len(value) < size, the remaining bytes are zero-filled.
// If len(value) > size, only the first `size` bytes of value are written.
// This is the primitive behind MSTORE (size == 32) and CALLDATACOPY-style
// opcodes (arbitrary size).
func (m *Memory) Set(offset, size uint64, value []byte) error {
	if size == 0 {
		return nil
	}
	end, ok := addUint64(offset, size)
	if !ok {
		return fmt.Errorf("vm: memory write overflow: offset=%d size=%d", offset, size)
	}

	m.Resize(end)

	n := size
	if uint64(len(value)) < n {
		n = uint64(len(value))
	}
	copy(m.store[offset:offset+n], value[:n])

	// Zero-fill any remainder (covers both "value shorter than size" and
	// stale bytes from a previous, now-overwritten use of this region).
	for i := offset + n; i < end; i++ {
		m.store[i] = 0
	}
	return nil
}

// Set8 writes a single byte at `offset`, resizing memory first if needed.
// This is the primitive behind MSTORE8.
func (m *Memory) Set8(offset uint64, value byte) error {
	end, ok := addUint64(offset, 1)
	if !ok {
		return fmt.Errorf("vm: memory write overflow: offset=%d", offset)
	}
	m.Resize(end)
	m.store[offset] = value
	return nil
}

// Get returns a copy of the `size` bytes starting at `offset`. If the
// requested range extends past the current memory bound, the missing
// portion is returned as zero bytes rather than causing an error — this
// mirrors EVM semantics where any address is "readable" (memory is
// conceptually infinite, just mostly zero). Callers that need gas-correct
// behavior should charge expansion cost and call Resize before Get, exactly
// as they must before Set.
func (m *Memory) Get(offset, size uint64) []byte {
	out := make([]byte, size)
	if size == 0 || offset >= uint64(len(m.store)) {
		return out
	}
	end := offset + size
	if end > uint64(len(m.store)) {
		end = uint64(len(m.store))
	}
	copy(out, m.store[offset:end])
	return out
}

// wordAlign rounds size up to the next multiple of wordSize (32).
func wordAlign(size uint64) uint64 {
	if size%wordSize == 0 {
		return size
	}
	return (size/wordSize + 1) * wordSize
}

// addUint64 adds a and b, reporting overflow instead of silently wrapping.
// Offsets/sizes in crafted bytecode can be adversarially large, so every
// offset+size computation in this package goes through this rather than a
// bare `+`.
func addUint64(a, b uint64) (uint64, bool) {
	sum := a + b
	if sum < a {
		return 0, false
	}
	return sum, true
}
