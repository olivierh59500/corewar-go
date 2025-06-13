package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Assembler converts Redcode source to instructions
type Assembler struct {
	labels map[string]int
}

// NewAssembler creates a new assembler instance
func NewAssembler() *Assembler {
	return &Assembler{
		labels: make(map[string]int),
	}
}

// Parse converts Redcode source into instructions
func (a *Assembler) Parse(source string) ([]Instruction, error) {
	lines := strings.Split(source, "\n")
	instructions := make([]Instruction, 0)

	// First pass: collect labels
	lineNum := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		// Skip END directive
		tokens := strings.Fields(line)
		if len(tokens) > 0 && strings.ToUpper(tokens[0]) == "END" {
			continue
		}

		// Check for label
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			label := strings.TrimSpace(parts[0])
			a.labels[label] = lineNum
			line = strings.TrimSpace(parts[1])

			// If nothing after label, continue
			if line == "" {
				continue
			}
		}

		lineNum++
	}

	// Second pass: parse instructions
	lineNum = 0
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		// Remove label if present
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			line = strings.TrimSpace(parts[1])
			if line == "" {
				continue
			}
		}

		// Parse instruction
		inst, err := a.parseInstruction(line, lineNum)
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", lineNum+1, err)
		}

		// Skip NOP instructions (used for END directive)
		if inst.Op != NOP {
			instructions = append(instructions, inst)
			lineNum++
		}
	}

	return instructions, nil
}

// parseInstruction parses a single instruction line
func (a *Assembler) parseInstruction(line string, currentLine int) (Instruction, error) {
	// Remove inline comments
	if idx := strings.Index(line, ";"); idx >= 0 {
		line = line[:idx]
	}
	line = strings.TrimSpace(line)

	// Split into tokens
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return Instruction{}, fmt.Errorf("empty instruction")
	}

	// Check for END directive
	if strings.ToUpper(tokens[0]) == "END" {
		// END directive - skip it
		return Instruction{Op: NOP}, nil
	}

	// Parse opcode
	op, err := parseOpCode(tokens[0])
	if err != nil {
		return Instruction{}, err
	}

	inst := Instruction{Op: op}

	// Parse operands based on instruction type
	switch op {
	case DAT:
		// DAT can have 0, 1, or 2 operands
		if len(tokens) > 1 {
			// Remove trailing comma
			tokens[1] = strings.TrimSuffix(tokens[1], ",")
			mode, value, err := a.parseOperand(tokens[1], currentLine)
			if err != nil {
				return Instruction{}, err
			}
			inst.AMode = mode
			inst.A = value
		}
		if len(tokens) > 2 {
			tokens[2] = strings.TrimSuffix(tokens[2], ",")
			mode, value, err := a.parseOperand(tokens[2], currentLine)
			if err != nil {
				return Instruction{}, err
			}
			inst.BMode = mode
			inst.B = value
		}

	case JMP, SPL:
		// Single operand instructions
		if len(tokens) < 2 {
			return Instruction{}, fmt.Errorf("%s requires an operand", tokens[0])
		}
		tokens[1] = strings.TrimSuffix(tokens[1], ",")
		mode, value, err := a.parseOperand(tokens[1], currentLine)
		if err != nil {
			return Instruction{}, err
		}
		inst.AMode = mode
		inst.A = value

	default:
		// Two operand instructions
		if len(tokens) < 3 {
			// Some warriors might have operands separated by comma without space
			// Try to split by comma
			if len(tokens) == 2 && strings.Contains(tokens[1], ",") {
				parts := strings.Split(tokens[1], ",")
				if len(parts) >= 2 {
					tokens = []string{tokens[0], parts[0], parts[1]}
				}
			}

			if len(tokens) < 3 {
				return Instruction{}, fmt.Errorf("%s requires two operands", tokens[0])
			}
		}

		// Remove comma if present
		tokens[1] = strings.TrimSuffix(tokens[1], ",")

		mode1, value1, err := a.parseOperand(tokens[1], currentLine)
		if err != nil {
			return Instruction{}, err
		}
		inst.AMode = mode1
		inst.A = value1

		tokens[2] = strings.TrimSuffix(tokens[2], ",")
		mode2, value2, err := a.parseOperand(tokens[2], currentLine)
		if err != nil {
			return Instruction{}, err
		}
		inst.BMode = mode2
		inst.B = value2
	}

	return inst, nil
}

// parseOpCode converts string to OpCode
func parseOpCode(s string) (OpCode, error) {
	switch strings.ToUpper(s) {
	case "DAT":
		return DAT, nil
	case "MOV":
		return MOV, nil
	case "ADD":
		return ADD, nil
	case "SUB":
		return SUB, nil
	case "JMP":
		return JMP, nil
	case "JMZ":
		return JMZ, nil
	case "JMN":
		return JMN, nil
	case "DJN":
		return DJN, nil
	case "CMP", "SEQ":
		return CMP, nil
	case "SPL":
		return SPL, nil
	case "NOP":
		return NOP, nil
	default:
		return DAT, fmt.Errorf("unknown opcode: %s", s)
	}
}

// parseOperand parses an operand into addressing mode and value
func (a *Assembler) parseOperand(s string, currentLine int) (AddressMode, int, error) {
	s = strings.TrimSpace(s)

	// Default mode
	mode := DIRECT

	// Check for addressing mode prefix
	if len(s) > 0 {
		switch s[0] {
		case '#':
			mode = IMMEDIATE
			s = s[1:]
		case '$':
			mode = DIRECT
			s = s[1:]
		case '@':
			mode = INDIRECT
			s = s[1:]
		case '<':
			mode = PREDECREMENT
			s = s[1:]
		case '>':
			mode = POSTINCREMENT
			s = s[1:]
		}
	}

	// Parse value
	value := 0

	// Check if it's a label
	if label, ok := a.labels[s]; ok {
		value = label - currentLine
	} else {
		// Try to parse as number
		var err error
		value, err = strconv.Atoi(s)
		if err != nil {
			return mode, 0, fmt.Errorf("invalid operand: %s", s)
		}
	}

	return mode, value, nil
}

// LoadWarriorFromSource creates a warrior from Redcode source
func LoadWarriorFromSource(name, source string, color WarriorColor) (*Warrior, error) {
	assembler := NewAssembler()
	instructions, err := assembler.Parse(source)
	if err != nil {
		return nil, err
	}

	return &Warrior{
		Name:  name,
		Code:  instructions,
		Color: color,
	}, nil
}
