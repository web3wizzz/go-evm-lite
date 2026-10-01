
package state

import (
	"math/big"
)

// journalEntry represents a reversible atomic change in StateDB
type journalEntry interface {
	revert(s *StateDB)
}

// balanceChange records a balance modification
type balanceChange struct {
	account []byte
	prev    *big.Int
}

func (ch balanceChange) revert(s *StateDB) {
	acc, err := s.GetAccount(ch.account)
	if err == nil {
		acc.Balance = ch.prev
		enc, _ := EncodeAccount(acc)
		s.trie.Put(ch.account, enc)
	}
}

// nonceChange records a nonce modification
type nonceChange struct {
	account []byte
	prev    uint64
}

func (ch nonceChange) revert(s *StateDB) {
	acc, err := s.GetAccount(ch.account)
	if err == nil {
		acc.Nonce = ch.prev
		enc, _ := EncodeAccount(acc)
		s.trie.Put(ch.account, enc)
	}
}

// storageChange records a contract storage slot modification
type storageChange struct {
	account []byte
	key     []byte
	prev    []byte
}

func (ch storageChange) revert(s *StateDB) {
	st, ok := s.storageTries[string(ch.account)]
	if ok {
		st.Put(ch.key, ch.prev)
	}
}
