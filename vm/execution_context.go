package vm

type ExecutionContext struct {
	Origin   [8]byte
	GasPrice uint64
	Data     []byte
	Value    uint64
}

type EVM struct {
	Ctx ExecutionContext
}
