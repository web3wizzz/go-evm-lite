package vm

import "fmt"

// ExecutionContext holds the per-call execution state: the memory, gas
// remaining, program counter, and (from Day 16) the 256-bit stack. This
// file focuses on the memory/gas wiring; Stack is assumed to already exist
// alongside this from your Day 16 work and is referenced here by name only.
type ExecutionContext struct {
	Memory       *Memory
	GasRemaining uint64
	PC           uint64
	Stack        *Stack
}

// NewExecutionContext creates a context with empty memory and the given
// starting gas.
func NewExecutionContext(gasLimit uint64) *ExecutionContext {
	return &ExecutionContext{
		Memory:       NewMemory(),
		Stack:        NewStack(),
		GasRemaining: gasLimit,
		PC:           0,
	}
}

// useGas deducts `amount` from GasRemaining, returning an out-of-gas error
// (and leaving GasRemaining at 0, never negative) if there isn't enough.
func (ctx *ExecutionContext) useGas(amount uint64) error {
	if amount > ctx.GasRemaining {
		have := ctx.GasRemaining
		ctx.GasRemaining = 0
		return fmt.Errorf("vm: out of gas: need %d, have %d", amount, have)
	}
	ctx.GasRemaining -= amount
	return nil
}

// ExpandMemory charges the gas cost of growing memory to cover
// [offset, offset+size), then performs the expansion. Every opcode that
// touches memory (MLOAD, MSTORE, MSTORE8, CODECOPY, CALLDATACOPY, the CALL
// family's return-data handling, etc.) should route through this rather
// than calling Memory.Resize directly, or memory will silently become free.
//
// size == 0 is a no-op (and costs nothing) regardless of offset — reading
// or writing zero bytes never touches memory, per the Yellow Paper.
func (ctx *ExecutionContext) ExpandMemory(offset, size uint64) error {
	if size == 0 {
		return nil
	}
	end, ok := addUint64(offset, size)
	if !ok {
		return fmt.Errorf("vm: memory expansion overflow: offset=%d size=%d", offset, size)
	}

	cost := MemoryExpansionGasCost(ctx.Memory.Len(), end)
	if err := ctx.useGas(cost); err != nil {
		return err
	}

	ctx.Memory.Resize(end)
	return nil
}

// MStore is the MSTORE opcode primitive: writes a 32-byte word at offset,
// charging memory expansion gas first.
func (ctx *ExecutionContext) MStore(offset uint64, word []byte) error {
	if err := ctx.ExpandMemory(offset, 32); err != nil {
		return err
	}
	return ctx.Memory.Set(offset, 32, word)
}

// MStore8 is the MSTORE8 opcode primitive: writes a single byte at offset,
// charging memory expansion gas first.
func (ctx *ExecutionContext) MStore8(offset uint64, value byte) error {
	if err := ctx.ExpandMemory(offset, 1); err != nil {
		return err
	}
	return ctx.Memory.Set8(offset, value)
}

// MLoad is the MLOAD opcode primitive: reads a 32-byte word at offset,
// charging memory expansion gas first (reading past the current bound
// expands memory too, since the EVM treats unread-but-touched memory as
// zero-initialized, not nonexistent).
func (ctx *ExecutionContext) MLoad(offset uint64) ([]byte, error) {
	if err := ctx.ExpandMemory(offset, 32); err != nil {
		return nil, err
	}
	return ctx.Memory.Get(offset, 32), nil
}
