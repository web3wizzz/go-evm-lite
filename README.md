# evm-lite

A lightweight Ethereum Virtual Machine built from scratch in Go.

`evm-lite` is an educational EVM implementation focused on understanding how Ethereum bytecode execution works: stack operations, memory, storage, gas, jumps, hashing, and contract execution context.

> Status: Core interpreter working and covered by tests.

## Features

### Core EVM execution

- Bytecode fetch-decode-execute loop
- Program counter (`PC`)
- Invalid opcode detection
- 1024-item stack limit
- 256-bit EVM words using `math/big`

### Arithmetic and comparison opcodes

- `ADD`, `SUB`, `MUL`, `DIV`, `MOD`
- `SDIV`, `SMOD`, `EXP`, `SIGNEXTEND`
- `LT`, `GT`, `SLT`, `SGT`, `EQ`, `ISZERO`

### Bitwise opcodes

- `AND`, `OR`, `XOR`, `NOT`
- `BYTE`, `SHL`, `SHR`, `SAR`

### Stack opcodes

- `PUSH1` through `PUSH32`
- `POP`
- `DUP1` through `DUP16`
- `SWAP1` through `SWAP16`

### Memory and storage

- Expandable byte-addressable memory
- Memory expansion gas calculation
- `MLOAD`, `MSTORE`, `MSTORE8`
- Persistent per-contract storage
- `SLOAD`, `SSTORE`

### Control flow and environment

- `JUMP`, `JUMPI`, `JUMPDEST`
- Jump destination validation
- `KECCAK256`
- `ADDRESS`, `BALANCE`, `CALLER`, `CALLVALUE`

### State layer

- Accounts with nonce, balance, storage root, and code hash
- Merkle Patricia Trie foundation
- `StateDB` account and storage abstraction
- State snapshots and rollback journal

## Project Structure

```text
evm-lite/
├── crypto/     # Keccak-256 hashing
├── state/      # Accounts, StateDB, journal, and trie implementation
├── vm/         # EVM interpreter, stack, memory, gas, and opcode tests
├── go.mod
└── README.md