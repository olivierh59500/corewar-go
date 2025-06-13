package main

import (
	"fmt"
	"time"
)

// BattleStats tracks statistics for a battle
type BattleStats struct {
	StartTime       time.Time
	EndTime         time.Time
	TotalCycles     int
	MaxProcesses    map[*Warrior]int
	InstructionsRun map[*Warrior]int
	Winner          *Warrior
	IsDraw          bool
}

// BattleManager manages battles between warriors
type BattleManager struct {
	core        *Core
	vm          *VM
	warriors    []*Warrior
	stats       *BattleStats
	maxCycles   int
	currentGame *Game
}

// NewBattleManager creates a new battle manager
func NewBattleManager(coreSize, maxCycles int) *BattleManager {
	return &BattleManager{
		maxCycles: maxCycles,
	}
}

// SetupBattle prepares a new battle
func (bm *BattleManager) SetupBattle(warriors []*Warrior) {
	bm.core = NewCore(coreSize)
	bm.vm = NewVM(bm.core)
	bm.warriors = warriors
	bm.stats = &BattleStats{
		StartTime:       time.Now(),
		MaxProcesses:    make(map[*Warrior]int),
		InstructionsRun: make(map[*Warrior]int),
	}

	// Calculate starting positions (evenly distributed)
	spacing := coreSize / len(warriors)
	for i, warrior := range warriors {
		position := i * spacing
		bm.loadWarriorAt(warrior, position)
		bm.stats.MaxProcesses[warrior] = 1
		bm.stats.InstructionsRun[warrior] = 0
	}
}

// loadWarriorAt loads a warrior at a specific position
func (bm *BattleManager) loadWarriorAt(warrior *Warrior, position int) {
	warrior.StartPosition = position

	// Copy warrior code to core
	for i, inst := range warrior.Code {
		addr := (position + i) % coreSize
		bm.core.cells[addr] = inst
		bm.core.owners[addr] = warrior.Color
	}

	// Add initial process
	bm.vm.AddProcess(warrior, position)
}

// RunCycle executes one cycle and updates statistics
func (bm *BattleManager) RunCycle() bool {
	if bm.stats.TotalCycles >= bm.maxCycles {
		bm.stats.IsDraw = true
		bm.stats.EndTime = time.Now()
		return false
	}

	// Get current process before execution
	if len(bm.vm.processes) > 0 && bm.vm.current < len(bm.vm.processes) {
		currentProc := bm.vm.processes[bm.vm.current]
		if currentProc.alive {
			// Increment instruction count for the warrior executing
			bm.stats.InstructionsRun[currentProc.warrior]++
		}
	}

	// Execute cycle
	bm.vm.ExecuteCycle()
	bm.stats.TotalCycles++

	// Count processes after cycle
	processCounts := make(map[*Warrior]int)
	for _, proc := range bm.vm.processes {
		if proc.alive {
			processCounts[proc.warrior]++
		}
	}

	// Update max processes
	for warrior, count := range processCounts {
		if count > bm.stats.MaxProcesses[warrior] {
			bm.stats.MaxProcesses[warrior] = count
		}
	}

	// Check for winner
	aliveWarriors := make([]*Warrior, 0)
	for _, warrior := range bm.warriors {
		if bm.vm.IsWarriorAlive(warrior) {
			aliveWarriors = append(aliveWarriors, warrior)
		}
	}

	if len(aliveWarriors) <= 1 {
		if len(aliveWarriors) == 1 {
			bm.stats.Winner = aliveWarriors[0]
		} else {
			bm.stats.IsDraw = true
		}
		bm.stats.EndTime = time.Now()
		return false
	}

	return true
}

// GetBattleReport generates a battle report
func (bm *BattleManager) GetBattleReport() string {
	duration := bm.stats.EndTime.Sub(bm.stats.StartTime)

	report := fmt.Sprintf("=== BATTLE REPORT ===\n")
	report += fmt.Sprintf("Duration: %v\n", duration)
	report += fmt.Sprintf("Total Cycles: %d\n", bm.stats.TotalCycles)
	report += fmt.Sprintf("\n")

	if bm.stats.IsDraw {
		report += "Result: DRAW\n"
	} else if bm.stats.Winner != nil {
		report += fmt.Sprintf("Winner: %s by %s\n", bm.stats.Winner.Name, bm.stats.Winner.Author)
	}

	report += fmt.Sprintf("\nWarrior Statistics:\n")
	for _, warrior := range bm.warriors {
		report += fmt.Sprintf("\n%s:\n", warrior.Name)
		report += fmt.Sprintf("  Max Processes: %d\n", bm.stats.MaxProcesses[warrior])
		report += fmt.Sprintf("  Instructions Run: %d\n", bm.stats.InstructionsRun[warrior])
		report += fmt.Sprintf("  Starting Position: %d\n", warrior.StartPosition)

		// Calculate efficiency
		if bm.stats.InstructionsRun[warrior] > 0 {
			efficiency := float64(bm.stats.InstructionsRun[warrior]) / float64(bm.stats.TotalCycles)
			report += fmt.Sprintf("  Efficiency: %.2f%%\n", efficiency*100)
		}
	}

	return report
}

// Tournament runs a tournament between multiple warriors
type Tournament struct {
	warriors     []*Warrior
	rounds       int
	coreSize     int
	maxCycles    int
	wins         map[*Warrior]int
	draws        int
	totalBattles int
}

// NewTournament creates a new tournament
func NewTournament(warriors []*Warrior, rounds, coreSize, maxCycles int) *Tournament {
	return &Tournament{
		warriors:  warriors,
		rounds:    rounds,
		coreSize:  coreSize,
		maxCycles: maxCycles,
		wins:      make(map[*Warrior]int),
	}
}

// Run executes the tournament
func (t *Tournament) Run() {
	fmt.Println("Starting Tournament...")
	fmt.Printf("Warriors: %d, Rounds per match: %d\n", len(t.warriors), t.rounds)

	// Round-robin: each warrior fights each other warrior
	for i := 0; i < len(t.warriors); i++ {
		for j := i + 1; j < len(t.warriors); j++ {
			w1, w2 := t.warriors[i], t.warriors[j]

			// Run multiple rounds
			for round := 0; round < t.rounds; round++ {
				// Alternate starting positions
				if round%2 == 0 {
					t.runBattle([]*Warrior{w1, w2})
				} else {
					t.runBattle([]*Warrior{w2, w1})
				}
			}
		}
	}

	// Print results
	t.printResults()
}

// runBattle runs a single battle
func (t *Tournament) runBattle(warriors []*Warrior) {
	bm := NewBattleManager(t.coreSize, t.maxCycles)
	bm.SetupBattle(warriors)

	// Run battle to completion
	for bm.RunCycle() {
		// Battle continues
	}

	t.totalBattles++

	// Record results
	if bm.stats.IsDraw {
		t.draws++
	} else if bm.stats.Winner != nil {
		t.wins[bm.stats.Winner]++
	}
}

// printResults displays tournament results
func (t *Tournament) printResults() {
	fmt.Println("\n=== TOURNAMENT RESULTS ===")
	fmt.Printf("Total Battles: %d\n", t.totalBattles)
	fmt.Printf("Draws: %d (%.1f%%)\n", t.draws, float64(t.draws)/float64(t.totalBattles)*100)

	fmt.Println("\nWarrior Rankings:")

	// Sort warriors by wins
	ranked := make([]*Warrior, 0, len(t.warriors))
	for _, w := range t.warriors {
		ranked = append(ranked, w)
	}

	// Simple bubble sort
	for i := 0; i < len(ranked); i++ {
		for j := i + 1; j < len(ranked); j++ {
			if t.wins[ranked[j]] > t.wins[ranked[i]] {
				ranked[i], ranked[j] = ranked[j], ranked[i]
			}
		}
	}

	// Display rankings
	for i, w := range ranked {
		winRate := float64(t.wins[w]) / float64(t.totalBattles) * 100
		fmt.Printf("%d. %s: %d wins (%.1f%%)\n", i+1, w.Name, t.wins[w], winRate)
	}
}
