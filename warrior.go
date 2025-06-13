package main

// Warrior represents a Core War program
type Warrior struct {
	Name          string
	Author        string
	Code          []Instruction
	StartPosition int
	Color         WarriorColor
}

// Some classic Core War warriors as examples

// CreateImp creates the classic Imp warrior
func CreateImp() *Warrior {
	return &Warrior{
		Name:   "Imp",
		Author: "A.K. Dewdney",
		Code: []Instruction{
			{Op: MOV, AMode: IMMEDIATE, BMode: DIRECT, A: 0, B: 1},
		},
		Color: Red,
	}
}

// CreateDwarf creates the classic Dwarf warrior
func CreateDwarf() *Warrior {
	return &Warrior{
		Name:   "Dwarf",
		Author: "A.K. Dewdney",
		Code: []Instruction{
			{Op: ADD, AMode: IMMEDIATE, BMode: DIRECT, A: 4, B: 3},
			{Op: MOV, AMode: IMMEDIATE, BMode: INDIRECT, A: 0, B: -2},
			{Op: JMP, AMode: DIRECT, BMode: DIRECT, A: -2, B: 0},
			{Op: DAT, AMode: IMMEDIATE, BMode: IMMEDIATE, A: 0, B: 0},
		},
		Color: Blue,
	}
}

// CreateStone creates a simple Stone warrior
func CreateStone() *Warrior {
	return &Warrior{
		Name:   "Stone",
		Author: "Core War Community",
		Code: []Instruction{
			{Op: MOV, AMode: DIRECT, BMode: INDIRECT, A: 2, B: -1},
			{Op: JMP, AMode: DIRECT, BMode: DIRECT, A: -1, B: 0},
			{Op: DAT, AMode: IMMEDIATE, BMode: IMMEDIATE, A: 0, B: 0},
		},
		Color: Green,
	}
}

// CreateGate creates a Gate warrior (simple scanner)
func CreateGate() *Warrior {
	return &Warrior{
		Name:   "Gate",
		Author: "Core War Community",
		Code: []Instruction{
			{Op: CMP, AMode: DIRECT, BMode: DIRECT, A: 9, B: 19},
			{Op: JMP, AMode: DIRECT, BMode: DIRECT, A: 3, B: 0},
			{Op: ADD, AMode: IMMEDIATE, BMode: DIRECT, A: 1, B: -1},
			{Op: JMP, AMode: DIRECT, BMode: DIRECT, A: -3, B: 0},
			{Op: SPL, AMode: DIRECT, BMode: DIRECT, A: 0, B: 0},
			{Op: MOV, AMode: DIRECT, BMode: INDIRECT, A: -4, B: -3},
			{Op: ADD, AMode: IMMEDIATE, BMode: DIRECT, A: 1, B: -1},
			{Op: JMP, AMode: DIRECT, BMode: DIRECT, A: -2, B: 0},
		},
		Color: Yellow,
	}
}

// CreateClearImp creates a self-clearing Imp
func CreateClearImp() *Warrior {
	return &Warrior{
		Name:   "Clear Imp",
		Author: "Core War Community",
		Code: []Instruction{
			{Op: MOV, AMode: IMMEDIATE, BMode: DIRECT, A: 0, B: 1},
			{Op: DAT, AMode: IMMEDIATE, BMode: IMMEDIATE, A: 0, B: 0},
		},
		Color: Red,
	}
}

// CreateSilkWarrior creates a Silk-style warrior (paper/stone hybrid)
func CreateSilkWarrior() *Warrior {
	return &Warrior{
		Name:   "Silk Warrior",
		Author: "Core War Community",
		Code: []Instruction{
			{Op: SPL, AMode: DIRECT, BMode: DIRECT, A: 1, B: 0},
			{Op: MOV, AMode: INDIRECT, BMode: DIRECT, A: -1, B: 0},
			{Op: MOV, AMode: DIRECT, BMode: POSTINCREMENT, A: -2, B: -2},
			{Op: DJN, AMode: DIRECT, BMode: DIRECT, A: -1, B: -3},
			{Op: SPL, AMode: DIRECT, BMode: DIRECT, A: -3, B: 0},
			{Op: SPL, AMode: INDIRECT, BMode: DIRECT, A: -1, B: 0},
			{Op: MOV, AMode: INDIRECT, BMode: DIRECT, A: -3, B: -5},
		},
		Color: Blue,
	}
}

// CreateBomber creates a spiral bomber warrior
func CreateBomber() *Warrior {
	return &Warrior{
		Name:   "Bomber",
		Author: "Core War Community",
		Code: []Instruction{
			{Op: MOV, AMode: DIRECT, BMode: INDIRECT, A: 2, B: 1},     // 0: MOV bomb, @ptr
			{Op: ADD, AMode: IMMEDIATE, BMode: DIRECT, A: 3, B: -1},   // 1: ADD #step, ptr
			{Op: JMP, AMode: DIRECT, BMode: DIRECT, A: -2, B: 0},      // 2: JMP start
			{Op: DAT, AMode: IMMEDIATE, BMode: IMMEDIATE, A: 0, B: 0}, // 3: bomb
			// ptr is at -1 (relative to ADD instruction)
		},
		Color: Red,
	}
}

// CreatePaperOne creates a paper/replicator warrior
func CreatePaperOne() *Warrior {
	return &Warrior{
		Name:   "Paper One",
		Author: "Core War Community",
		Code: []Instruction{
			// Boot phase
			{Op: SPL, AMode: DIRECT, BMode: DIRECT, A: 1, B: 0},    // 0: SPL 1
			{Op: MOV, AMode: INDIRECT, BMode: DIRECT, A: -1, B: 0}, // 1: MOV @-1, 0

			// Copy phase
			{Op: MOV, AMode: POSTINCREMENT, BMode: POSTINCREMENT, A: 1, B: -1}, // 2: MOV >1, >-1
			{Op: SPL, AMode: INDIRECT, BMode: DIRECT, A: -1, B: 0},             // 3: SPL @-1
			{Op: DJN, AMode: DIRECT, BMode: DIRECT, A: -2, B: -3},              // 4: DJN -2, -3

			// Data
			{Op: DAT, AMode: IMMEDIATE, BMode: IMMEDIATE, A: 0, B: 100}, // 5: DAT #0, #100
		},
		Color: Blue,
	}
}

// CreateVampire creates a vampire warrior (converts enemy processes)
func CreateVampire() *Warrior {
	return &Warrior{
		Name:   "Vampire",
		Author: "Core War Community",
		Code: []Instruction{
			// Vampire pit
			{Op: JMP, AMode: DIRECT, BMode: DIRECT, A: 0, B: 0}, // 0: pit JMP pit

			// Main loop
			{Op: MOV, AMode: DIRECT, BMode: INDIRECT, A: -1, B: 3},     // 1: MOV pit, @fang
			{Op: ADD, AMode: IMMEDIATE, BMode: DIRECT, A: 2, B: 2},     // 2: ADD #step, fang
			{Op: JMP, AMode: DIRECT, BMode: DIRECT, A: -2, B: 0},       // 3: JMP loop
			{Op: DAT, AMode: IMMEDIATE, BMode: IMMEDIATE, A: 0, B: 10}, // 4: fang
		},
		Color: Green,
	}
}

// CreateQuickScan creates a fast scanner warrior
func CreateQuickScan() *Warrior {
	return &Warrior{
		Name:   "QuickScan",
		Author: "Core War Community",
		Code: []Instruction{
			// Quick scan with bombing
			{Op: ADD, AMode: DIRECT, BMode: DIRECT, A: 3, B: 1},        // 0: ADD step, scan
			{Op: CMP, AMode: DIRECT, BMode: DIRECT, A: 12, B: 2},       // 1: scan CMP 12, 2
			{Op: JMP, AMode: DIRECT, BMode: DIRECT, A: 4, B: 0},        // 2: JMP attack
			{Op: JMP, AMode: DIRECT, BMode: DIRECT, A: -3, B: 0},       // 3: JMP loop
			{Op: MOV, AMode: DIRECT, BMode: DIRECT, A: 5, B: -1},       // 4: attack MOV bomb, scan-1
			{Op: JMP, AMode: DIRECT, BMode: DIRECT, A: -5, B: 0},       // 5: JMP loop
			{Op: DAT, AMode: IMMEDIATE, BMode: IMMEDIATE, A: 0, B: 10}, // 6: step
			{Op: DAT, AMode: IMMEDIATE, BMode: IMMEDIATE, A: 0, B: 0},  // 7: bomb
		},
		Color: Yellow,
	}
}

// CreateSilkPaper creates an advanced paper warrior
func CreateSilkPaper() *Warrior {
	return &Warrior{
		Name:   "Silk Paper",
		Author: "Core War Community",
		Code: []Instruction{
			// Silk-style paper with anti-imp
			{Op: SPL, AMode: DIRECT, BMode: DIRECT, A: 1, B: 0},                 // 0
			{Op: MOV, AMode: INDIRECT, BMode: DIRECT, A: -1, B: 0},              // 1
			{Op: SPL, AMode: DIRECT, BMode: DIRECT, A: 1, B: 0},                 // 2
			{Op: MOV, AMode: POSTINCREMENT, BMode: POSTINCREMENT, A: -3, B: -2}, // 3
			{Op: SPL, AMode: INDIRECT, BMode: DIRECT, A: -3, B: 0},              // 4
			{Op: MOV, AMode: DIRECT, BMode: POSTINCREMENT, A: 4, B: -5},         // 5
			{Op: DJN, AMode: DIRECT, BMode: DIRECT, A: -3, B: -5},               // 6
			{Op: SPL, AMode: INDIRECT, BMode: DIRECT, A: -1, B: 0},              // 7
			{Op: DAT, AMode: IMMEDIATE, BMode: IMMEDIATE, A: 0, B: 0},           // 8
		},
		Color: Red,
	}
}

// CreateOneShot creates a one-shot scanner
func CreateOneShot() *Warrior {
	return &Warrior{
		Name:   "OneShot",
		Author: "Core War Community",
		Code: []Instruction{
			// One-shot scanner with SPL/DAT clear
			{Op: ADD, AMode: IMMEDIATE, BMode: DIRECT, A: 5, B: 4},     // 0
			{Op: CMP, AMode: DIRECT, BMode: DIRECT, A: 2, B: -3},       // 1
			{Op: JMP, AMode: DIRECT, BMode: DIRECT, A: 4, B: 0},        // 2
			{Op: JMP, AMode: DIRECT, BMode: DIRECT, A: -3, B: 0},       // 3
			{Op: SPL, AMode: DIRECT, BMode: DIRECT, A: 0, B: 0},        // 4: ptr
			{Op: MOV, AMode: DIRECT, BMode: INDIRECT, A: 2, B: -1},     // 5
			{Op: MOV, AMode: DIRECT, BMode: PREDECREMENT, A: 1, B: -2}, // 6
			{Op: JMP, AMode: DIRECT, BMode: DIRECT, A: -7, B: 0},       // 7
			{Op: DAT, AMode: IMMEDIATE, BMode: IMMEDIATE, A: 0, B: 0},  // 8
		},
		Color: Blue,
	}
}
