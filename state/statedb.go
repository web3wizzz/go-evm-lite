package state

import (
	"encoding/binary"
	"fmt"
	"math/big"
)

// StateDB manages account states, contract storage tries, and state modification journaling
type StateDB struct {
	trie         *Trie
	storageTries map[string]*Trie
	journal      []journalEntry
}

// NewStateDB initializes a new StateDB with empty state, storage tries, and journal
func NewStateDB() *StateDB {
	return &StateDB{
		trie:         NewTrie(),
		storageTries: make(map[string]*Trie),
		journal:      make([]journalEntry, 0),
	}
}

// Snapshot creates a revision point in the journal and returns its revision ID
func (s *StateDB) Snapshot() int {
	return len(s.journal)
}

// RevertToSnapshot rolls back all state changes made after the specified snapshot ID
func (s *StateDB) RevertToSnapshot(revertID int) {
	if revertID < 0 || revertID > len(s.journal) {
		return
	}

	// Revert changes in reverse order down to target snapshot ID
	for i := len(s.journal) - 1; i >= revertID; i-- {
		s.journal[i].revert(s)
	}

	s.journal = s.journal[:revertID]
}

// GetAccount retrieves and decodes an Account by its 20-byte address
func (s *StateDB) GetAccount(addr []byte) (*Account, error) {
	enc, err := s.trie.Get(addr)
	if err != nil {
		return NewAccount(0, new(big.Int)), nil
	}

	return DecodeAccount(enc)
}

// SetAccount serializes an Account and stores it in the state trie
func (s *StateDB) SetAccount(addr []byte, acc *Account) error {
	enc, err := EncodeAccount(acc)
	if err != nil {
		return fmt.Errorf("failed to encode account: %v", err)
	}
	s.trie.Put(addr, enc)
	return nil
}

// GetBalance returns an account's balance in wei
func (s *StateDB) GetBalance(addr []byte) *big.Int {
	acc, err := s.GetAccount(addr)
	if err != nil {
		return new(big.Int)
	}
	return acc.Balance
}

// SetBalance sets an account's balance in wei and logs a journal entry
func (s *StateDB) SetBalance(addr []byte, amount *big.Int) error {
	acc, err := s.GetAccount(addr)
	prev := new(big.Int)
	if err == nil && acc.Balance != nil {
		prev = new(big.Int).Set(acc.Balance)
	} else {
		acc = NewAccount(0, amount)
	}

	// Journal change before mutation
	s.journal = append(s.journal, balanceChange{account: addr, prev: prev})

	acc.Balance = amount
	return s.SetAccount(addr, acc)
}

// GetNonce returns an account's transaction nonce
func (s *StateDB) GetNonce(addr []byte) uint64 {
	acc, err := s.GetAccount(addr)
	if err != nil {
		return 0
	}
	return acc.Nonce
}

// SetNonce updates an account's transaction nonce and logs a journal entry
func (s *StateDB) SetNonce(addr []byte, nonce uint64) error {
	acc, err := s.GetAccount(addr)
	prev := uint64(0)
	if err == nil {
		prev = acc.Nonce
	} else {
		acc = NewAccount(nonce, new(big.Int))
	}

	// Journal change before mutation
	s.journal = append(s.journal, nonceChange{account: addr, prev: prev})

	acc.Nonce = nonce
	return s.SetAccount(addr, acc)
}

// GetState retrieves a 32-byte storage value from an account's storage trie
func (s *StateDB) GetState(addr, key []byte) []byte {
	st, ok := s.storageTries[string(addr)]
	if !ok {
		return make([]byte, 32)
	}
	val, err := st.Get(key)
	if err != nil {
		return make([]byte, 32)
	}
	return val
}

// SetState sets a 32-byte storage value in an account's storage trie and logs a journal entry
func (s *StateDB) SetState(addr, key, value []byte) error {
	st, ok := s.storageTries[string(addr)]
	if !ok {
		st = NewTrie()
		s.storageTries[string(addr)] = st
	}

	// Journal previous storage value before mutation
	prev := s.GetState(addr, key)
	s.journal = append(s.journal, storageChange{account: addr, key: key, prev: prev})

	st.Put(key, value)

	acc, err := s.GetAccount(addr)
	if err != nil {
		acc = NewAccount(0, new(big.Int))
	}
	return s.SetAccount(addr, acc)
}

const accountEncodedSize = 8 + 32 + 32 + 32

// EncodeAccount serializes an account using fixed-width fields:
// nonce (8 bytes), balance (32 bytes), storage root (32 bytes), code hash (32 bytes).
func EncodeAccount(acc *Account) ([]byte, error) {
	if acc == nil {
		return nil, fmt.Errorf("cannot encode nil account")
	}

	enc := make([]byte, accountEncodedSize)

	// Nonce: bytes 0–7
	binary.BigEndian.PutUint64(enc[:8], acc.Nonce)

	// Balance: bytes 8–39, left-padded to 32 bytes.
	if acc.Balance != nil {
		balanceBytes := acc.Balance.Bytes()
		if len(balanceBytes) > 32 {
			return nil, fmt.Errorf("balance exceeds 256 bits")
		}
		copy(enc[40-len(balanceBytes):40], balanceBytes)
	}

	// Storage root: bytes 40–71
	copy(enc[40:72], acc.StorageRoot)

	// Code hash: bytes 72–103
	copy(enc[72:104], acc.CodeHash)

	return enc, nil
}

// DecodeAccount reconstructs an Account from its fixed-width encoding.
func DecodeAccount(data []byte) (*Account, error) {
	if len(data) != accountEncodedSize {
		return nil, fmt.Errorf(
			"invalid account data length: got %d, want %d",
			len(data),
			accountEncodedSize,
		)
	}

	return &Account{
		Nonce:       binary.BigEndian.Uint64(data[:8]),
		Balance:     new(big.Int).SetBytes(data[8:40]),
		StorageRoot: append([]byte(nil), data[40:72]...),
		CodeHash:    append([]byte(nil), data[72:104]...),
	}, nil
}
