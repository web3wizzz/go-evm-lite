package vm

// Assembler builds EVM bytecode without manually writing opcode bytes.
type Assembler struct {
	code []byte
}

// NewAssembler creates an empty bytecode builder.
func NewAssembler() *Assembler {
	return &Assembler{}
}

// Op appends one opcode.
func (a *Assembler) Op(opcode byte) *Assembler {
	a.code = append(a.code, opcode)
	return a
}

// Push1 appends a PUSH1 instruction and its one-byte value.
func (a *Assembler) Push1(value byte) *Assembler {
	a.code = append(a.code, PUSH1, value)
	return a
}

// Bytes returns a safe copy of the assembled bytecode.
func (a *Assembler) Bytes() []byte {
	return append([]byte(nil), a.code...)
}
