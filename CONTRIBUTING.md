

```markdown
# Contributing to evm-lite

Thank you for your interest in contributing to **evm-lite**, a lightweight, zero-dependency Ethereum Virtual Machine (EVM) built in Go from scratch based on the Ethereum Yellow Paper.

---

## 🚀 Getting Started

### Prerequisites
- **Go**: Version 1.22 or higher
- **Git**: For version control
- **Environment**: Linux, macOS, or GitHub Codespaces (recommended)

### Local Setup

1. **Fork and clone the repository:**
   ```bash
      git clone https://github.com/your-username/evm-lite.git
         cd evm-lite
            ```

            2. **Verify tests pass:**
               ```bash
                  go test -v ./...
                     ```

                     ---

                     ## 🛠️ Project Structure

                     The codebase strictly follows modular packages aligned with core Ethereum protocol specifications:

                     ```
                     evm-lite/
                     ├── crypto/          # Keccak-256 hashing, RLP encoding & decoding
                     ├── state/           # Merkle Patricia Trie (MPT) nodes, Hex-Prefix encoding
                     ├── vm/              # EVM opcode interpreter, execution context, stack, memory
                     └── main.go          # CLI entry point
                     ```

                     ---

                     ## 📋 Development & Coding Standards

                     To maintain standard protocol compliance and clean architecture:

                     1. **Standard Library First**: Do not add external third-party dependencies unless explicitly discussed. Use Go standard packages (`crypto`, `encoding/binary`, `fmt`, `testing`).
                     2. **Code Formatting**: Always run `go fmt` before committing code:
                        ```bash
                           go fmt ./...
                              ```
                              3. **Spec Alignment**: Ensure function signatures and mathematical logic strictly align with Gavin Wood's Ethereum Yellow Paper (e.g., Appendix B for RLP, Appendix C for HP encoding).
                              4. **Testing**: Every package must include unit tests (`*_test.go`). Roundtrip tests are required for encoders/decoders.

                              ---

                              ## 🔀 Pull Request Process

                              1. **Create a Feature Branch**:
                                 ```bash
                                    git checkout -b feat/day-X-short-description
                                       ```
                                       2. **Commit Message Format**:
                                          Use conventional commits:
                                             - `feat(crypto): implement RLP decoder`
                                                - `test(state): add hex-prefix unit test cases`
                                                   - `fix(vm): resolve stack overflow on PUSH32`
                                                   3. **Run Tests**: Ensure all package tests pass locally:
                                                      ```bash
                                                         go test -v ./...
                                                            ```
                                                            4. **Submit PR**: Open a Pull Request targeting `main`. Include test execution outputs in your PR description.

                                                            ---

                                                            ## 🧪 Running Benchmarks & Tests

                                                            Run verbose tests for specific packages:
                                                            ```bash
                                                            # Test crypto package
                                                            go test -v ./crypto

                                                            # Test state management
                                                            go test -v ./state
                                                            ```

                                                            ---

                                                            Thank you for helping build **evm-lite**!
                                                            ```