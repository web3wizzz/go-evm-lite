package state

type Account struct {
	Nonce       uint64 `json:"nonce"`
	Balance     uint64 `json:"balance"`
	StorageRoot []byte `json:"storageRoot"`
	CodeHash    []byte `json:"codeHash"`
}
