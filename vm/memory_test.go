package vm

import (
	"bytes"
	"testing"
)

func TestMemory_InitiallyEmpty(t *testing.T) {
	m := NewMemory()
	if m.Len() != 0 {
		t.Fatalf("expected 0, got %d", m.Len())
	}
}

func TestMemory_ResizeRoundsUpToWord(t *testing.T) {
	m := NewMemory()
	m.Resize(1)
	if m.Len() != 32 {
		t.Fatalf("expected 32 (word-aligned), got %d", m.Len())
	}

	m2 := NewMemory()
	m2.Resize(33)
	if m2.Len() != 64 {
		t.Fatalf("expected 64 (word-aligned), got %d", m2.Len())
	}
}

func TestMemory_ResizeNeverShrinks(t *testing.T) {
	m := NewMemory()
	m.Resize(64)
	m.Resize(1) // smaller request should be a no-op
	if m.Len() != 64 {
		t.Fatalf("expected memory to stay at 64, got %d", m.Len())
	}
}

func TestMemory_SetAndGet(t *testing.T) {
	m := NewMemory()
	data := []byte("hello world")
	if err := m.Set(0, uint64(len(data)), data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Get(0, uint64(len(data)))
	if !bytes.Equal(got, data) {
		t.Fatalf("expected %q, got %q", data, got)
	}
}

func TestMemory_SetZeroFillsShortValue(t *testing.T) {
	m := NewMemory()
	// Write a 3-byte value into a 5-byte region; the remaining 2 bytes
	// should be zero.
	if err := m.Set(0, 5, []byte{0x01, 0x02, 0x03}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Get(0, 5)
	want := []byte{0x01, 0x02, 0x03, 0x00, 0x00}
	if !bytes.Equal(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestMemory_SetTruncatesLongValue(t *testing.T) {
	m := NewMemory()
	if err := m.Set(0, 2, []byte{0xaa, 0xbb, 0xcc, 0xdd}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Get(0, 2)
	want := []byte{0xaa, 0xbb}
	if !bytes.Equal(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestMemory_Set8(t *testing.T) {
	m := NewMemory()
	if err := m.Set8(5, 0x42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Get(5, 1)
	if !bytes.Equal(got, []byte{0x42}) {
		t.Fatalf("expected [0x42], got %v", got)
	}
	// Memory should still have expanded by a full word, not just 1 byte.
	if m.Len() != 32 {
		t.Fatalf("expected memory len 32 after touching offset 5, got %d", m.Len())
	}
}

func TestMemory_GetPastBoundReturnsZeros(t *testing.T) {
	m := NewMemory() // len 0
	got := m.Get(0, 10)
	want := make([]byte, 10)
	if !bytes.Equal(got, want) {
		t.Fatalf("expected all-zero read past bound, got %v", got)
	}
}

func TestMemory_GetAutoExpandsRegionDoesNotMutate(t *testing.T) {
	// Get should never mutate memory, even when reading past the current
	// bound -- only Resize/Set/Set8 (and, at the context layer,
	// ExpandMemory) should grow the backing store.
	m := NewMemory()
	_ = m.Get(0, 100)
	if m.Len() != 0 {
		t.Fatalf("expected Get to leave memory untouched, got len %d", m.Len())
	}
}

func TestMemory_WriteOverflowDetected(t *testing.T) {
	m := NewMemory()
	// offset + size overflows uint64.
	err := m.Set(^uint64(0)-1, 10, []byte{0x01})
	if err == nil {
		t.Fatalf("expected overflow error, got nil")
	}
}
