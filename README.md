# go-evm-lite

A lightweight, educational Ethereum Virtual Machine built from scratch in Go.

**Status:** Educational Implementation | **Roadmap:** 12-week Development Cycle

---

## Table of Contents

- [Overview](#overview)
- [Architecture & Directory Structure](#architecture--directory-structure)
- [Core Concepts & Specification](#core-concepts--specification)
- [Feature Comparison](#feature-comparison)
- [Getting Started](#getting-started)
- [Usage Examples](#usage-examples)
- [Project Roadmap](#project-roadmap)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

`evm-lite` is a minimal Ethereum Virtual Machine implementation designed for educational purposes. Built entirely in Go, it demonstrates core EVM concepts including bytecode execution, gas metering, state management, and smart contract execution.

### Goals

- **Educational**: Understand EVM internals through a clean, understandable implementation
- **Minimal**: Implement only essential EVM components without unnecessary complexity
- **Functional**: Execute real smart contract bytecode and manage blockchain state
- **Well-Documented**: Clear code and comprehensive documentation for learning

### Key Features

- Full EVM execution loop (fetch-decode-execute)
- Stack-based computation (1024 depth, 256-bit words)
- Expandable byte-addressable memory
- Persistent account storage with Merkle Patricia Trie
- Gas metering and cost accounting
- Smart contract creation and execution
- Keccak-256 cryptographic hashing

---

## Architecture & Directory Structure

```
go-evm-lite/
├── core/              # Transaction and block execution
├── state/             # Account state and storage
├── vm/                # Virtual machine execution engine
├── crypto/            # Cryptographic utilities
├── tests/             # Integration and unit tests
├── examples/          # Example contracts and usage
├── go.mod
└── README.md
```

### Package Responsibilities

#### `core/`
Handles transaction processing, block headers, and the execution engine state transitions.

**Key Components:**
- `Block`: Block header and metadata (number, timestamp, coinbase, gas limit)
- `Transaction`: Transaction structure (to, from, value, data, gas price, gas limit)
- `ExecutionEngine`: Orchestrates state transitions and applies transactions to state
- `Receipt`: Transaction receipt (status, gas used, logs)

#### `state/`
Manages the account state model and persistent storage layer using Merkle Patricia Trie.

**Key Components:**
- `Account`: Account structure with fields:
  - `Nonce`: Transaction count (prevents replay attacks)
  - `Balance`: Account ether balance (in Wei)
  - `StorageRoot`: Root hash of account's storage trie
  - `CodeHash`: Hash of account bytecode (empty for EOAs)
- `StateDB`: CRUD interface for account and storage operations
  - Get/Set account balance, nonce, code
  - Get/Set contract storage values
  - Commit state changes to persistent trie
- `MerklePatriciaTrie`: Efficient storage proof structure
  - Path compression using hex prefixes
  - Keccak-256 hashing of node values

#### `vm/`
The execution engine core — interprets and executes EVM bytecode.

**Key Components:**
- `ExecutionContext`: Runtime context containing:
  - Stack: 1024-depth stack of 256-bit words (using `github.com/holiman/uint256`)
  - Memory: Dynamically expanding byte array (32-byte word aligned)
  - Storage: Persistent key-value store for contract data
  - Caller: Transaction sender information
  - CallValue: Ether sent with transaction
- `Interpreter`: The fetch-decode-execute loop
  - Program counter (PC) management
  - Opcode dispatch and execution
  - Exception handling (out of gas, stack underflow, etc.)
- `GasMeter`: Gas accounting and cost calculation
  - Base costs per opcode
  - Memory expansion costs (quadratic)
  - Storage operation costs (cold/warm accesses)
- `Opcodes`: Implementation of 140+ EVM operations:
  - Arithmetic: ADD, SUB, MUL, DIV, MOD, ADDMOD, MULMOD
  - Comparison: LT, GT, EQ, ISZERO
  - Bitwise: AND, OR, XOR, NOT, SHL, SHR, SAR
  - Cryptographic: SHA3 (Keccak-256)
  - Environmental: CALLER, CALLVALUE, CODESIZE, CODECOPY
  - Storage: SLOAD, SSTORE
  - Memory: MLOAD, MSTORE, MSTORE8
  - Control Flow: JUMP, JUMPI, JUMPDEST
  - Contract: CALL, DELEGATECALL, CREATE
  - Halt: STOP, REVERT, SELFDESTRUCT

#### `crypto/`
Cryptographic utilities supporting the EVM.

**Key Components:**
- `Keccak256`: Keccak-256 hashing (using `golang.org/x/crypto/sha3`)
- `Hash`: Type-safe hash wrapper with utility methods
- `Address`: 20-byte Ethereum address with checksum validation

---

## Core Concepts & Specification

### State Transition Function

The EVM operates as a deterministic state machine. The state transition function is defined as:

```
Υ(σ, T) = σ'
```

Where:
- **σ** (sigma): Current world state (mapping of addresses to accounts)
- **T**: Transaction to execute
- **σ'** (sigma prime): Resulting state after transaction execution

For a given transaction T = (nonce, gasPrice, gasLimit, to, value, data):

1. **Validation**: Check nonce, sender balance, and gas
2. **Intrinsic Gas**: Calculate base transaction cost
3. **Execution**: 
   - If `to` is null: Contract creation with bytecode from `data`
   - Otherwise: Code execution (or value transfer if `to` is EOA)
4. **State Update**: Apply state changes (storage writes, balance transfers)
5. **Finalization**: Log events, calculate refunds, return receipt

### Account Model

Ethereum uses an account-based model with two account types:

#### Externally Owned Accounts (EOAs)
- Controlled by private key
- Cannot execute code
- Have nonce, balance; storage is empty
- CodeHash = Keccak256(empty bytes)

#### Contract Accounts
- Created by transaction with no recipient
- Controlled by bytecode
- Have nonce, balance, storage, and associated code
- CodeHash = Keccak256(bytecode)

### EVM Memory Model

**Stack:**
- LIFO data structure, 1024 items maximum
- Each item is a 256-bit (32-byte) word
- Operations are word-oriented (not byte-oriented)

**Memory:**
- Byte-addressable array of bytes
- Expands dynamically as needed (32-byte word aligned)
- Cost grows quadratically with expansion: `memCost = (size² / 512) + (32 × size)`
- Volatile: reset for each transaction

**Storage:**
- Persistent key-value store (256-bit key, 256-bit value)
- Associated with contract account
- Part of world state, survives across transactions
- Cold access costs 2100 gas, warm access 100 gas (post-EIP-2929)

### Opcode Execution

Each EVM instruction is a 1-byte opcode (0x00-0xff). Execution follows:

```
while pc < len(code) {
    opcode = code[pc]
    (stackDelta, gasCost) = opInfo[opcode]
    gasRemaining -= gasCost
    pc = execute(opcode)
}
```

Gas prevents infinite loops and resource exhaustion. If `gasRemaining < 0`, execution halts with `OutOfGas` exception.

---

## Feature Comparison

| Feature | `evm-lite` | `go-ethereum` (Geth) | Notes |
|---------|-----------|---------------------|-------|
| **Execution** | | | |
| Opcode Interpreter | ✅ Full | ✅ Full + JIT | evm-lite uses pure Go interpreter |
| Stack Operations | ✅ Full (1024 depth) | ✅ Full | Identical semantics |
| Memory Model | ✅ Full | ✅ Full | Dynamic expansion with gas cost |
| Storage (SLOAD/SSTORE) | ✅ Full | ✅ Full + Cold/Warm | evm-lite simplified gas costs |
| **State Management** | | | |
| Account Model | ✅ EOA + Contracts | ✅ Full | Identical |
| Nonce Tracking | ✅ Full | ✅ Full | Prevents replay attacks |
| Balance Management | ✅ Full | ✅ Full | Wei precision (big.Int) |
| Storage Trie | ✅ Merkle Patricia | ✅ Merkle Patricia | State proofs not verified |
| Code Storage | ✅ Full | ✅ Full | Contract bytecode storage |
| **Gas** | | | |
| Intrinsic Gas | ✅ Full | ✅ Full | 21,000 base + data cost |
| Opcode Costs | ✅ Core opcodes | ✅ 140+ opcodes | All essential opcodes implemented |
| Memory Costs | ✅ Quadratic | ✅ Quadratic | Identical formula |
| Storage Costs | ✅ Simplified | ✅ Full (EIP-2929) | evm-lite: no cold/warm distinction |
| Refunds | ✅ Basic | ✅ Full | SSTORE refunds only |
| **Cryptography** | | | |
| Keccak-256 | ✅ Full | ✅ Full | Uses golang.org/x/crypto/sha3 |
| Address Derivation | ✅ Full | ✅ Full | RLP encoding of contract create |
| **Consensus** | | | |
| PoW | ❌ Not implemented | ✅ Full (deprecated) | Not needed for execution |
| PoS/Beacon | ❌ Not implemented | ✅ Full | Outside EVM scope |
| State Finality | ✅ Implicit | ✅ Full + Consensus rules | evm-lite assumes valid state |
| **Precompiles** | | | |
| ECRECOVER | ❌ Stub | ✅ Full | Signature recovery |
| SHA256 | ❌ Stub | ✅ Full | Hash function |
| RIPEMD-160 | ❌ Stub | ✅ Full | Hash function |
| Identity | ✅ Full | ✅ Full | Input passthrough |
| Modexp | ❌ Stub | ✅ Full | Modular exponentiation |
| **Advanced Features** | | | |
| Contract Creation | ✅ Basic | ✅ Full (EIP-170) | Creates account, deploys code |
| CALL/DELEGATECALL | ✅ Basic | ✅ Full + Analysis | Simplified gas accounting |
| Logs/Events | ❌ Not implemented | ✅ Full (Bloom filters) | Observable in go-ethereum |
| Reverts | ✅ REVERT opcode | ✅ Full + Reason strings | Partial state atomicity |
| State Channels | ❌ Not applicable | ✅ Sidechains/Rollups | Outside EVM scope |

---

## Getting Started

### Prerequisites

- **Go 1.21+** ([Download](https://golang.org/dl))
- **Git**

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/go-evm-lite.git
cd go-evm-lite

# Download dependencies
go mod download

# Verify installation
go test ./...
```

### Basic Usage

```go
package main

import (
    "go-evm-lite/core"
    "go-evm-lite/state"
    "go-evm-lite/vm"
)

func main() {
    // Initialize state database
    stateDB := state.NewStateDB()
    
    // Create execution engine
    engine := core.NewExecutionEngine(stateDB)
    
    // Create and execute a transaction
    tx := &core.Transaction{
        Nonce:    0,
        GasPrice: big.NewInt(1),
        GasLimit: 21000,
        To:       nil, // Contract creation
        Value:    big.NewInt(0),
        Data:     bytecode, // Smart contract bytecode
    }
    
    receipt, err := engine.ExecuteTransaction(tx, &core.Block{})
    if err != nil {
        panic(err)
    }
    
    println("Gas Used:", receipt.GasUsed)
    println("Status:", receipt.Status)
}
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test -v ./vm/...

# Run with coverage
go test -cover ./...
```

---

## Usage Examples

### Example 1: Simple Counter Contract

Deploy a contract that stores and increments a counter:

```go
// Counter contract bytecode (simplified)
// PUSH1 0x00 (slot 0)
// SLOAD
// PUSH1 0x01
// ADD
// PUSH1 0x00
// SSTORE

counterBytecode := []byte{
    0x60, 0x00,       // PUSH1 0
    0x54,             // SLOAD
    0x60, 0x01,       // PUSH1 1
    0x01,             // ADD
    0x60, 0x00,       // PUSH1 0
    0x55,             // SSTORE
    0x60, 0x01,       // PUSH1 1
    0x60, 0x00,       // PUSH1 0
    0xf3,             // RETURN
}

// Execute transaction to deploy contract
tx := &core.Transaction{
    Data:     counterBytecode,
    GasLimit: 100000,
}
```

### Example 2: Token Transfer

Simulate an ETH transfer between accounts:

```go
alice := common.HexToAddress("0x1234...") 
bob := common.HexToAddress("0x5678...")

// Set Alice's balance
stateDB.SetBalance(alice, big.NewInt(1e18)) // 1 ETH

// Create transfer transaction
tx := &core.Transaction{
    From:  alice,
    To:    bob,
    Value: big.NewInt(5e17), // 0.5 ETH
}

// Execute
receipt, _ := engine.ExecuteTransaction(tx, block)
println("Alice balance:", stateDB.GetBalance(alice))
println("Bob balance:", stateDB.GetBalance(bob))
```

### Example 3: Storage Operations

Interact with contract storage:

```go
contractAddr := common.HexToAddress("0xaaaa...")

// Write to storage
key := common.HexToHash("0x01")
value := common.HexToHash("0xdeadbeef")
stateDB.SetState(contractAddr, key, value)

// Read from storage
retrieved := stateDB.GetState(contractAddr, key)
println("Stored value:", retrieved.Hex())
```

---

## Project Roadmap

### Week 1-2: Foundation
- [x] Account model and state structure
- [x] Stack and memory implementations
- [x] Basic opcodes (arithmetic, logical)
- [ ] Transaction model and execution context

### Week 3-4: Core Execution
- [ ] Full opcode set (stack, memory, storage)
- [ ] Gas metering and cost calculation
- [ ] Program counter and control flow
- [ ] Exception handling

### Week 5-6: State Management
- [ ] Merkle Patricia Trie implementation
- [ ] StateDB CRUD operations
- [ ] Account creation and code deployment
- [ ] Storage state transitions

### Week 7-8: Advanced Features
- [ ] Contract-to-contract calls (CALL, DELEGATECALL)
- [ ] Contract creation (CREATE)
- [ ] Revert semantics and atomicity
- [ ] Event logging (optional)

### Week 9-10: Testing & Optimization
- [ ] Comprehensive test suite
- [ ] Integration with Ethereum test vectors
- [ ] Performance profiling and optimization
- [ ] Bytecode validation

### Week 11-12: Documentation & Polish
- [ ] Complete API documentation
- [ ] Tutorial examples and guides
- [ ] Performance benchmarks
- [ ] Release preparation

---

## Contributing

Contributions are welcome! Areas for contribution:

1. **Opcode Implementations**: Complete all 140+ EVM opcodes
2. **Precompiles**: Add cryptographic precompiled contracts
3. **Testing**: Expand test coverage with Ethereum test vectors
4. **Documentation**: Improve code comments and guides
5. **Optimization**: Performance improvements and refactoring
6. **Examples**: More real-world contract examples

### Development Guidelines

- Follow Go best practices and idioms
- Write tests for all new features
- Update documentation with API changes
- Keep commits atomic and well-described

---

## License

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for details.

Educational use of this software is encouraged. For production use, always use a battle-tested implementation like [go-ethereum](https://github.com/ethereum/go-ethereum).

---

## Resources

### Learning Materials
- [Ethereum Yellow Paper](https://ethereum.org/en/whitepaper/) — Formal EVM specification
- [EVM.codes](https://www.evm.codes/) — Interactive opcode reference
- [go-ethereum Source](https://github.com/ethereum/go-ethereum/tree/master/core/vm) — Production reference implementation

### Related Projects
- [go-ethereum](https://github.com/ethereum/go-ethereum) — Full Ethereum client
- [ethereum-tests](https://github.com/ethereum/tests) — Ethereum test vectors
- [solc](https://github.com/ethereum/solidity) — Solidity compiler

### Community
- [Ethereum Research](https://ethresear.ch/)
- [Ethereum Discord](https://discord.gg/ethereum)
- [Go Community](https://golang.org/)

---

**Built with ❤️ for Ethereum education**
