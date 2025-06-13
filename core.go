package main

import (
	"image/color"
)

// OpCode represents the instruction operation codes
type OpCode int

const (
	DAT OpCode = iota // Data (kills process)
	MOV               // Move
	ADD               // Add
	SUB               // Subtract
	JMP               // Jump
	JMZ               // Jump if Zero
	JMN               // Jump if Not Zero
	DJN               // Decrement and Jump if Not Zero
	CMP               // Compare (skip if equal)
	SPL               // Split (create new process)
	NOP               // No operation
)

// AddressMode represents the addressing modes
type AddressMode int

const (
	IMMEDIATE     AddressMode = iota // #
	DIRECT                           // $
	INDIRECT                         // @
	PREDECREMENT                     // <
	POSTINCREMENT                    // >
)

// Instruction represents a single Core War instruction
type Instruction struct {
	Op    OpCode
	AMode AddressMode
	BMode AddressMode
	A     int
	B     int
}

// WarriorColor represents the color assigned to a warrior
type WarriorColor int

const (
	Empty WarriorColor = iota
	Red
	Blue
	Green
	Yellow
)

// Core represents the memory core where warriors battle
type Core struct {
	cells       []Instruction
	owners      []WarriorColor // Track which warrior owns each cell
	readEffect  []float32      // Visual effect for read operations
	writeEffect []float32      // Visual effect for write operations
	execEffect  []float32      // Visual effect for execution
	size        int
}

// NewCore creates a new memory core
func NewCore(size int) *Core {
	c := &Core{
		cells:       make([]Instruction, size),
		owners:      make([]WarriorColor, size),
		readEffect:  make([]float32, size),
		writeEffect: make([]float32, size),
		execEffect:  make([]float32, size),
		size:        size,
	}

	// Initialize all cells with DAT #0, #0
	for i := 0; i < size; i++ {
		c.cells[i] = Instruction{
			Op:    DAT,
			AMode: IMMEDIATE,
			BMode: IMMEDIATE,
			A:     0,
			B:     0,
		}
		c.owners[i] = Empty
	}

	return c
}

// Read returns the instruction at the given address
func (c *Core) Read(addr int) Instruction {
	addr = c.normalize(addr)
	c.readEffect[addr] = 1.0
	return c.cells[addr]
}

// Write stores an instruction at the given address
func (c *Core) Write(addr int, inst Instruction, owner WarriorColor) {
	addr = c.normalize(addr)
	c.cells[addr] = inst
	c.owners[addr] = owner
	c.writeEffect[addr] = 1.0
}

// Execute marks a cell as being executed
func (c *Core) Execute(addr int) {
	addr = c.normalize(addr)
	c.execEffect[addr] = 1.0
}

// normalize ensures address is within core bounds
func (c *Core) normalize(addr int) int {
	addr = addr % c.size
	if addr < 0 {
		addr += c.size
	}
	return addr
}

// DecayEffects reduces the intensity of visual effects over time
func (c *Core) DecayEffects() {
	decay := float32(0.95)
	for i := range c.cells {
		c.readEffect[i] *= decay
		c.writeEffect[i] *= decay
		c.execEffect[i] *= decay

		// Clamp to zero if too small
		if c.readEffect[i] < 0.01 {
			c.readEffect[i] = 0
		}
		if c.writeEffect[i] < 0.01 {
			c.writeEffect[i] = 0
		}
		if c.execEffect[i] < 0.01 {
			c.execEffect[i] = 0
		}
	}
}

// GetColor returns the color for a given warrior
func GetWarriorColor(wc WarriorColor) color.RGBA {
	switch wc {
	case Red:
		return color.RGBA{255, 0, 0, 255}
	case Blue:
		return color.RGBA{0, 0, 255, 255}
	case Green:
		return color.RGBA{0, 255, 0, 255}
	case Yellow:
		return color.RGBA{255, 255, 0, 255}
	default:
		return color.RGBA{50, 50, 50, 255}
	}
}

// OpCodeString returns the string representation of an opcode
func OpCodeString(op OpCode) string {
	switch op {
	case DAT:
		return "DAT"
	case MOV:
		return "MOV"
	case ADD:
		return "ADD"
	case SUB:
		return "SUB"
	case JMP:
		return "JMP"
	case JMZ:
		return "JMZ"
	case JMN:
		return "JMN"
	case DJN:
		return "DJN"
	case CMP:
		return "CMP"
	case SPL:
		return "SPL"
	case NOP:
		return "NOP"
	default:
		return "???"
	}
}

// AddressModeString returns the string representation of an address mode
func AddressModeString(mode AddressMode) string {
	switch mode {
	case IMMEDIATE:
		return "#"
	case DIRECT:
		return "$"
	case INDIRECT:
		return "@"
	case PREDECREMENT:
		return "<"
	case POSTINCREMENT:
		return ">"
	default:
		return "?"
	}
}
