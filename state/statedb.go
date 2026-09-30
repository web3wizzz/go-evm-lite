package state

import (
	"fmt"
	"math/big"

	"evm-lite/crypto"
)

// StateDB manages account states using an underlying Merkle Patricia Trie
type StateDB struct {
	trie *Trie
}

// NewStateDB initializes a new StateDB with an empty MPT
func NewStateDB() *StateDB {
	return &StateDB{
		trie: NewTrie(),
	}
}

// GetAccount retrieves and decodes an Account by its 20-byte address
func (s *StateDB) GetAccount(addr []byte) (*Account, error) {
	enc, err := s.trie.Get(addr)
	if err != nil {
		// Return a fresh default account if address does not exist yet
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

// SetBalance sets an account's balance in wei
func (s *StateDB) SetBalance(addr []byte, amount *big.Int) error {
	acc, err := s.GetAccount(addr)
	if err != nil {
		acc = NewAccount(0, amount)
	} else {
		acc.Balance = amount
	}
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

// SetNonce updates an account's transaction nonce
func (s *StateDB) SetNonce(addr []byte, nonce uint64) error {
	acc, err := s.GetAccount(addr)
	if err != nil {
		acc = NewAccount(nonce, new(big.Int))
	} else {
		acc.Nonce = nonce
	}
	return s.SetAccount(addr, acc)
}

// Helper: EncodeAccount serializes an Account struct into an RLP-encoded state tuple.
func EncodeAccount(acc *Account) ([]byte, error) {
	if acc == nil {
		return nil, fmt.Errorf("account is nil")
	}

	balance := new(big.Int)
	if acc.Balance != nil {
		balance = new(big.Int).Set(acc.Balance)
	}

	storageRoot := append([]byte(nil), acc.StorageRoot...)
	if len(storageRoot) == 0 {
		storageRoot = append([]byte(nil), EmptyRootHash...)
	}

	codeHash := append([]byte(nil), acc.CodeHash...)
	if len(codeHash) == 0 {
		codeHash = append([]byte(nil), EmptyCodeHash...)
	}

	items := [][]byte{
		crypto.EncodeUint(acc.Nonce),
		crypto.EncodeBytes(balance.Bytes()),
		crypto.EncodeBytes(storageRoot),
		crypto.EncodeBytes(codeHash),
	}

	return crypto.EncodeList(items), nil
}

// Helper: DecodeAccount reconstructs an Account struct from the RLP-encoded state tuple.
func DecodeAccount(data []byte) (*Account, error) {
	items, _, err := crypto.DecodeList(data)
	if err != nil {
		return nil, fmt.Errorf("invalid account data: %w", err)
	}
	if len(items) != 4 {
		return nil, fmt.Errorf("invalid account data length")
	}

	nonce := uint64(0)
	if len(items[0]) > 0 {
		nonce = new(big.Int).SetBytes(items[0]).Uint64()
	}

	balance := new(big.Int).SetBytes(items[1])
	storageRoot := append([]byte(nil), items[2]...)
	codeHash := append([]byte(nil), items[3]...)

	if len(storageRoot) == 0 {
		storageRoot = append([]byte(nil), EmptyRootHash...)
	}
	if len(codeHash) == 0 {
		codeHash = append([]byte(nil), EmptyCodeHash...)
	}

	return &Account{
		Nonce:       nonce,
		Balance:     balance,
		StorageRoot: storageRoot,
		CodeHash:    codeHash,
	}, nil
}
