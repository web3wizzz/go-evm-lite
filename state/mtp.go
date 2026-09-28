package state

import (
	"bytes"
	"fmt"
)

// Trie represents a Merkle Patricia Trie
type Trie struct {
	root Node
}

// NewTrie creates a new empty MPT
func NewTrie() *Trie {
	return &Trie{root: nil}
}

// Root returns the root node of the trie
func (t *Trie) Root() Node {
	return t.root
}

// Get retrieves the value associated with key
func (t *Trie) Get(key []byte) ([]byte, error) {
	nibbles := BytesToNibbles(key)
	return t.get(t.root, nibbles)
}

func (t *Trie) get(node Node, nibbles []byte) ([]byte, error) {
	if node == nil {
		return nil, fmt.Errorf("key not found")
	}

	switch n := node.(type) {
	case *LeafNode:
		if bytes.Equal(n.Path, nibbles) {
			return n.Value, nil
		}
		return nil, fmt.Errorf("key not found")

	case *ExtensionNode:
		matchLen := prefixMatchLen(n.Path, nibbles)
		if matchLen == len(n.Path) {
			return t.get(n.Child, nibbles[matchLen:])
		}
		return nil, fmt.Errorf("key not found")

	case *BranchNode:
		if len(nibbles) == 0 {
			if n.Value != nil {
				return n.Value, nil
			}
			return nil, fmt.Errorf("key not found")
		}
		idx := nibbles[0]
		return t.get(n.ChildrenNode[idx], nibbles[1:])

	default:
		return nil, fmt.Errorf("unknown node type")
	}
}

// Put inserts or updates a key-value pair in the trie
func (t *Trie) Put(key, value []byte) {
	nibbles := BytesToNibbles(key)
	t.root = t.put(t.root, nibbles, value)
}

func (t *Trie) put(node Node, nibbles []byte, value []byte) Node {
	if node == nil {
		return &LeafNode{Path: nibbles, Value: value}
	}

	switch n := node.(type) {
	case *LeafNode:
		matchLen := prefixMatchLen(n.Path, nibbles)

		// Exact match: update value
		if matchLen == len(n.Path) && matchLen == len(nibbles) {
			return &LeafNode{Path: nibbles, Value: value}
		}

		// Split leaf into a branch
		branch := &BranchNode{}

		// Handle remaining leaf path
		if matchLen == len(n.Path) {
			branch.Value = n.Value
		} else {
			leafIdx := n.Path[matchLen]
			branch.ChildrenNode[leafIdx] = &LeafNode{
				Path:  n.Path[matchLen+1:],
				Value: n.Value,
			}
		}

		// Handle new key/value
		if matchLen == len(nibbles) {
			branch.Value = value
		} else {
			newIdx := nibbles[matchLen]
			branch.ChildrenNode[newIdx] = &LeafNode{
				Path:  nibbles[matchLen+1:],
				Value: value,
			}
		}

		// If there is a shared prefix before divergence, wrap branch in ExtensionNode
		if matchLen > 0 {
			return &ExtensionNode{
				Path:  n.Path[:matchLen],
				Child: branch,
			}
		}
		return branch

	case *ExtensionNode:
		matchLen := prefixMatchLen(n.Path, nibbles)

		// Entire extension path matches: recurse into child
		if matchLen == len(n.Path) {
			n.Child = t.put(n.Child, nibbles[matchLen:], value)
			return n
		}

		// Split extension node
		branch := &BranchNode{}

		// Existing extension child branch
		extIdx := n.Path[matchLen]
		if matchLen+1 == len(n.Path) {
			branch.ChildrenNode[extIdx] = n.Child
		} else {
			branch.ChildrenNode[extIdx] = &ExtensionNode{
				Path:  n.Path[matchLen+1:],
				Child: n.Child,
			}
		}

		// New key insertion
		if matchLen == len(nibbles) {
			branch.Value = value
		} else {
			newIdx := nibbles[matchLen]
			branch.ChildrenNode[newIdx] = &LeafNode{
				Path:  nibbles[matchLen+1:],
				Value: value,
			}
		}

		if matchLen > 0 {
			return &ExtensionNode{
				Path:  n.Path[:matchLen],
				Child: branch,
			}
		}
		return branch

	case *BranchNode:
		if len(nibbles) == 0 {
			n.Value = value
			return n
		}
		idx := nibbles[0]
		n.ChildrenNode[idx] = t.put(n.ChildrenNode[idx], nibbles[1:], value)
		return n

	default:
		panic("unknown node type")
	}
}

// Helper: Calculate common prefix length between two nibble slices
func prefixMatchLen(a, b []byte) int {
	maxLen := len(a)
	if len(b) < maxLen {
		maxLen = len(b)
	}
	i := 0
	for i < maxLen && a[i] == b[i] {
		i++
	}
	return i
}
