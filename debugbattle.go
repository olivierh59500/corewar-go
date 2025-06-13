package main

import (
	"fmt"
)

// DebugBattle runs a battle with detailed debugging
func DebugBattle() {
	fmt.Println("=== DEBUG: Imp vs Dwarf ===")

	// Create small core for debugging
	core := NewCore(100)
	vm := NewVM(core)

	// Create warriors
	imp := CreateImp()
	dwarf := CreateDwarf()

	// Load at known positions
	impStart := 0
	dwarfStart := 50

	// Load Imp
	for i, inst := range imp.Code {
		core.cells[impStart+i] = inst
		core.owners[impStart+i] = imp.Color
	}
	vm.AddProcess(imp, impStart)

	// Load Dwarf
	for i, inst := range dwarf.Code {
		core.cells[dwarfStart+i] = inst
		core.owners[dwarfStart+i] = dwarf.Color
	}
	vm.AddProcess(dwarf, dwarfStart)

	// Show initial state
	fmt.Printf("\nInitial Dwarf code at %d:\n", dwarfStart)
	for i := 0; i < 4; i++ {
		inst := core.cells[dwarfStart+i]
		fmt.Printf("  %d: %s A=%d B=%d\n", dwarfStart+i, OpCodeString(inst.Op), inst.A, inst.B)
	}

	// Run first 10 cycles
	for cycle := 0; cycle < 10; cycle++ {
		fmt.Printf("\n--- Cycle %d ---\n", cycle)

		// Show DAT pointer value
		datAddr := dwarfStart + 3
		dat := core.cells[datAddr]
		fmt.Printf("Dwarf DAT at %d: B=%d\n", datAddr, dat.B)

		// Execute one cycle
		vm.ExecuteCycle()

		// Check for bombs
		bombAddr := (datAddr + dat.B) % 100
		if core.cells[bombAddr].Op == DAT && core.owners[bombAddr] == dwarf.Color {
			fmt.Printf("BOMB placed at %d!\n", bombAddr)
		}

		// Show Imp position
		if len(vm.processes) > 0 && vm.processes[0].alive {
			fmt.Printf("Imp at position %d\n", vm.processes[0].pc)
		}
	}
}
