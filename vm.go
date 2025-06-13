package main

import "fmt"

// Process represents a single execution thread for a warrior
type Process struct {
	warrior *Warrior
	pc      int // program counter
	alive   bool
}

// VM represents the Core War virtual machine (MARS)
type VM struct {
	core      *Core
	processes []*Process
	current   int // current process index
}

// NewVM creates a new virtual machine
func NewVM(core *Core) *VM {
	return &VM{
		core:      core,
		processes: make([]*Process, 0),
	}
}

// AddProcess adds a new process for a warrior
func (vm *VM) AddProcess(warrior *Warrior, startAddr int) {
	vm.processes = append(vm.processes, &Process{
		warrior: warrior,
		pc:      startAddr,
		alive:   true,
	})
}

// ExecuteCycle executes one cycle of the VM
func (vm *VM) ExecuteCycle() {
	if len(vm.processes) == 0 {
		return
	}

	// Remove dead processes
	alive := make([]*Process, 0)
	for _, p := range vm.processes {
		if p.alive {
			alive = append(alive, p)
		}
	}
	vm.processes = alive

	if len(vm.processes) == 0 {
		return
	}

	// Get current process
	vm.current = vm.current % len(vm.processes)
	proc := vm.processes[vm.current]

	// Execute instruction
	vm.executeInstruction(proc)

	// Move to next process
	vm.current = (vm.current + 1) % len(vm.processes)
}

// executeInstruction executes a single instruction for a process
func (vm *VM) executeInstruction(proc *Process) {
	// Mark execution location
	vm.core.Execute(proc.pc)

	// Fetch instruction
	inst := vm.core.Read(proc.pc)

	// Check if this location still belongs to the warrior
	// If another warrior has overwritten this location, the process dies
	if vm.core.owners[proc.pc] != proc.warrior.Color && vm.core.owners[proc.pc] != Empty {
		// This location has been overwritten by another warrior
		proc.alive = false
		return
	}

	// Calculate next PC (will be overridden by jump instructions)
	nextPC := (proc.pc + 1) % vm.core.size

	// Execute based on opcode
	switch inst.Op {
	case DAT:
		// Data instruction kills the process
		proc.alive = false
		if debugMode {
			fmt.Printf("Process died executing DAT at PC=%d (warrior: %s)\n", proc.pc, proc.warrior.Name)
		}
		return

	case MOV:
		// Move instruction
		source := vm.evaluate(proc.pc, inst.AMode, inst.A, false)
		dest := vm.evaluate(proc.pc, inst.BMode, inst.B, true)

		if debugMode && proc.warrior.Name == "ImpKiller" {
			fmt.Printf("ImpKiller MOV: PC=%d, source=%d, dest=%d\n", proc.pc, source, dest)
			fmt.Printf("  Instruction: AMode=%d, BMode=%d, A=%d, B=%d\n", inst.AMode, inst.BMode, inst.A, inst.B)
		}

		if inst.AMode == IMMEDIATE {
			// Moving immediate value - creates a DAT instruction
			// This is typically a bomb!
			vm.core.Write(dest, Instruction{
				Op:    DAT,
				AMode: IMMEDIATE,
				BMode: IMMEDIATE,
				A:     source,
				B:     0,
			}, proc.warrior.Color)
			if debugMode && proc.warrior.Name == "ImpKiller" {
				fmt.Printf("  Placed DAT bomb at %d\n", dest)
			}
		} else {
			// Moving instruction
			srcInst := vm.core.Read(source)
			if debugMode && proc.warrior.Name == "ImpKiller" {
				fmt.Printf("  Moving instruction from %d: Op=%d\n", source, srcInst.Op)
			}
			vm.core.Write(dest, srcInst, proc.warrior.Color)
		}
		proc.pc = nextPC

	case ADD:
		// Add instruction
		source := vm.evaluate(proc.pc, inst.AMode, inst.A, false)
		dest := vm.evaluate(proc.pc, inst.BMode, inst.B, true)

		if debugMode && proc.warrior.Name == "Dwarf" {
			fmt.Printf("Dwarf ADD: PC=%d, source=%d, dest=%d\n", proc.pc, source, dest)
		}

		if inst.AMode == IMMEDIATE {
			// Add immediate value to B field
			destInst := vm.core.Read(dest)
			oldValue := destInst.B
			destInst.B = (destInst.B + source) % vm.core.size
			vm.core.Write(dest, destInst, proc.warrior.Color)
			if debugMode && proc.warrior.Name == "Dwarf" {
				fmt.Printf("  Added %d to B field at %d: %d -> %d\n", source, dest, oldValue, destInst.B)
			}
		} else {
			// Add instruction fields
			srcInst := vm.core.Read(source)
			destInst := vm.core.Read(dest)
			destInst.A = (destInst.A + srcInst.A) % vm.core.size
			destInst.B = (destInst.B + srcInst.B) % vm.core.size
			vm.core.Write(dest, destInst, proc.warrior.Color)
		}
		proc.pc = nextPC

	case SUB:
		// Subtract instruction
		source := vm.evaluate(proc.pc, inst.AMode, inst.A, false)
		dest := vm.evaluate(proc.pc, inst.BMode, inst.B, true)

		if inst.AMode == IMMEDIATE {
			// Subtract immediate value
			destInst := vm.core.Read(dest)
			destInst.B = (destInst.B - source + vm.core.size) % vm.core.size
			vm.core.Write(dest, destInst, proc.warrior.Color)
		} else {
			// Subtract instruction fields
			srcInst := vm.core.Read(source)
			destInst := vm.core.Read(dest)
			destInst.A = (destInst.A - srcInst.A + vm.core.size) % vm.core.size
			destInst.B = (destInst.B - srcInst.B + vm.core.size) % vm.core.size
			vm.core.Write(dest, destInst, proc.warrior.Color)
		}
		proc.pc = nextPC

	case JMP:
		// Jump instruction
		target := vm.evaluate(proc.pc, inst.AMode, inst.A, false)
		proc.pc = target

	case JMZ:
		// Jump if zero
		source := vm.evaluate(proc.pc, inst.AMode, inst.A, false)
		target := vm.evaluate(proc.pc, inst.BMode, inst.B, false)

		value := 0
		if inst.AMode == IMMEDIATE {
			value = source
		} else {
			srcInst := vm.core.Read(source)
			value = srcInst.B
		}

		if value == 0 {
			proc.pc = target
		} else {
			proc.pc = nextPC
		}

	case JMN:
		// Jump if not zero
		source := vm.evaluate(proc.pc, inst.AMode, inst.A, false)
		target := vm.evaluate(proc.pc, inst.BMode, inst.B, false)

		value := 0
		if inst.AMode == IMMEDIATE {
			value = source
		} else {
			srcInst := vm.core.Read(source)
			value = srcInst.B
		}

		if value != 0 {
			proc.pc = target
		} else {
			proc.pc = nextPC
		}

	case DJN:
		// Decrement and jump if zero (DJZ in original spec)
		source := vm.evaluate(proc.pc, inst.AMode, inst.A, true)
		target := vm.evaluate(proc.pc, inst.BMode, inst.B, false)

		// Decrement the content at location A
		if inst.AMode == IMMEDIATE {
			// Can't decrement immediate value
			proc.pc = nextPC
		} else {
			// According to spec: "Decrement contents of location A by 1"
			value := vm.core.Read(source)
			value.B = (value.B - 1 + vm.core.size) % vm.core.size
			vm.core.Write(source, value, proc.warrior.Color)

			// "If location A now holds 0, jump to location B"
			if value.B == 0 {
				proc.pc = target
			} else {
				proc.pc = nextPC
			}
		}

	case CMP:
		// Compare and skip if equal
		source1 := vm.evaluate(proc.pc, inst.AMode, inst.A, false)
		source2 := vm.evaluate(proc.pc, inst.BMode, inst.B, false)

		equal := false
		if inst.AMode == IMMEDIATE && inst.BMode == IMMEDIATE {
			equal = source1 == source2
		} else if inst.AMode == IMMEDIATE {
			inst2 := vm.core.Read(source2)
			equal = source1 == inst2.B
		} else if inst.BMode == IMMEDIATE {
			inst1 := vm.core.Read(source1)
			equal = inst1.B == source2
		} else {
			inst1 := vm.core.Read(source1)
			inst2 := vm.core.Read(source2)
			// For CMP, we typically just compare the B fields unless it's a full comparison
			equal = inst1.B == inst2.B
		}

		if equal {
			proc.pc = (nextPC + 1) % vm.core.size // Skip next instruction
		} else {
			proc.pc = nextPC
		}

	case SPL:
		// Split - create new process
		target := vm.evaluate(proc.pc, inst.AMode, inst.A, false)

		// Count current processes for this warrior
		processCount := 0
		for _, p := range vm.processes {
			if p.warrior == proc.warrior && p.alive {
				processCount++
			}
		}

		// Limit processes per warrior (optional)
		const maxProcessesPerWarrior = 64
		if processCount >= maxProcessesPerWarrior {
			// Cannot create more processes
			proc.pc = nextPC
			return
		}

		// Create new process at target
		newProc := &Process{
			warrior: proc.warrior,
			pc:      target,
			alive:   true,
		}

		// Insert new process after current one
		newProcesses := make([]*Process, 0, len(vm.processes)+1)
		for i, p := range vm.processes {
			newProcesses = append(newProcesses, p)
			if i == vm.current {
				newProcesses = append(newProcesses, newProc)
			}
		}
		vm.processes = newProcesses

		proc.pc = nextPC

	default:
		// Unknown instruction acts like NOP
		proc.pc = nextPC
	}
}

// evaluate resolves an address based on the addressing mode
func (vm *VM) evaluate(pc int, mode AddressMode, operand int, write bool) int {
	switch mode {
	case IMMEDIATE:
		return operand

	case DIRECT:
		return vm.core.normalize(pc + operand)

	case INDIRECT:
		pointer := vm.core.normalize(pc + operand)
		inst := vm.core.Read(pointer)
		// For indirect addressing, we use the B field as the offset
		return vm.core.normalize(pointer + inst.B)

	case PREDECREMENT:
		pointer := vm.core.normalize(pc + operand)
		inst := vm.core.Read(pointer)
		inst.B = (inst.B - 1 + vm.core.size) % vm.core.size
		if write {
			vm.core.Write(pointer, inst, Empty) // Don't change owner on indirect updates
		}
		return vm.core.normalize(pointer + inst.B)

	case POSTINCREMENT:
		pointer := vm.core.normalize(pc + operand)
		inst := vm.core.Read(pointer)
		result := vm.core.normalize(pointer + inst.B)
		inst.B = (inst.B + 1) % vm.core.size
		if write {
			vm.core.Write(pointer, inst, Empty) // Don't change owner on indirect updates
		}
		return result

	default:
		return vm.core.normalize(pc + operand)
	}
}

// IsWarriorAlive checks if a warrior has any alive processes
func (vm *VM) IsWarriorAlive(warrior *Warrior) bool {
	for _, proc := range vm.processes {
		if proc.warrior == warrior && proc.alive {
			return true
		}
	}
	return false
}
