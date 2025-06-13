package main

import (
	"fmt"
)

// TestDwarfBehavior tests the Dwarf warrior behavior
func TestDwarfBehavior() {
	fmt.Println("=== Testing Dwarf Behavior ===")

	// Create a small core for testing
	testCore := NewCore(100)
	testVM := NewVM(testCore)

	// Create and load Dwarf at position 10
	dwarf := CreateDwarf()
	startPos := 10

	// Load Dwarf code
	for i, inst := range dwarf.Code {
		addr := (startPos + i) % 100
		testCore.cells[addr] = inst
		testCore.owners[addr] = dwarf.Color
	}

	// Add process
	testVM.AddProcess(dwarf, startPos)

	// Execute first few cycles and show what happens
	for cycle := 0; cycle < 10; cycle++ {
		fmt.Printf("\nCycle %d:\n", cycle)

		// Show the DAT instruction that serves as pointer
		datAddr := (startPos + 3) % 100
		datInst := testCore.cells[datAddr]
		fmt.Printf("  DAT at %d: A=%d, B=%d\n", datAddr, datInst.A, datInst.B)

		// Execute one cycle
		proc := testVM.processes[0]
		fmt.Printf("  Executing at PC=%d\n", proc.pc)
		testVM.ExecuteCycle()

		// Check if any bombs were placed
		if cycle > 0 {
			// The bomb should be placed at position pointed by DAT.B
			bombAddr := (datAddr + datInst.B) % 100
			if testCore.cells[bombAddr].Op == DAT && testCore.owners[bombAddr] == dwarf.Color {
				fmt.Printf("  BOMB placed at %d!\n", bombAddr)
			}
		}
	}
}

// Add this function call in main() for testing
func RunDwarfTest() {
	TestDwarfBehavior()
}
