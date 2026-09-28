
package state

import (
	"bytes"
	"testing"
)

func TestTriePutAndGet(t *testing.T) {
	trie := NewTrie()

	// Test 1: Single insertion and retrieval
	trie.Put([]byte("cat"), []byte("dog"))
	val, err := trie.Get([]byte("cat"))
	if err != nil {
		t.Fatalf("Get('cat') failed: %v", err)
	}
	if !bytes.Equal(val, []byte("dog")) {
		t.Errorf("Get('cat') = %s; want 'dog'", string(val))
	}

	// Test 2: Diverging keys (triggers leaf splitting into branch)
	trie.Put([]byte("cab"), []byte("car"))
	val, err = trie.Get([]byte("cab"))
	if err != nil {
		t.Fatalf("Get('cab') failed: %v", err)
	}
	if !bytes.Equal(val, []byte("car")) {
		t.Errorf("Get('cab') = %s; want 'car'", string(val))
	}

	// Test 3: Retrieve original key after branch split
	val, err = trie.Get([]byte("cat"))
	if err != nil {
		t.Fatalf("Get('cat') failed after split: %v", err)
	}
	if !bytes.Equal(val, []byte("dog")) {
		t.Errorf("Get('cat') = %s; want 'dog'", string(val))
	}

	// Test 4: Key not found
	_, err = trie.Get([]byte("nonexistent"))
	if err == nil {
		t.Errorf("Get('nonexistent') expected error, got nil")
	}
}
