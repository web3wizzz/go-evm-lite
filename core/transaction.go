package core

type Transaction struct {
	Nonce    uint64
	GasPrice uint64
	GasLimit uint64
	To       *[8]byte
	Value    uint64
	Data     []byte
}
