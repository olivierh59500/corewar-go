package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1024
	screenHeight = 768
	coreSize     = 8000  // Standard core size
	maxCycles    = 80000 // Maximum cycles before draw
	debugMode    = false // Set to true for debug output
	cellSpacing  = 1     // Spacing between cells
	infoHeight   = 100   // Height of info panel
)

// Game represents the main game state
type Game struct {
	core         *Core
	vm           *VM
	warriors     []*Warrior
	paused       bool
	speed        int // cycles per frame
	cycle        int
	keyPressed   map[ebiten.Key]bool
	battleMgr    *BattleManager
	gameOver     bool
	battleReport string
}

// NewGame creates a new game instance
func NewGame(warriors []*Warrior) *Game {
	g := &Game{
		speed:      10,
		keyPressed: make(map[ebiten.Key]bool),
		battleMgr:  NewBattleManager(coreSize, maxCycles),
	}

	// Setup battle
	g.battleMgr.SetupBattle(warriors)
	g.core = g.battleMgr.core
	g.vm = g.battleMgr.vm
	g.warriors = warriors

	return g
}

// Update handles game logic updates
func (g *Game) Update() error {
	// Handle keyboard input
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if !g.keyPressed[ebiten.KeySpace] {
			g.paused = !g.paused
			g.keyPressed[ebiten.KeySpace] = true
		}
	} else {
		g.keyPressed[ebiten.KeySpace] = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if g.speed < 100 {
			g.speed++
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if g.speed > 1 {
			g.speed--
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		// Restart the game
		*g = *NewGame(g.warriors)
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// Execute cycles if not paused and game not over
	if !g.paused && !g.gameOver {
		for i := 0; i < g.speed; i++ {
			if !g.battleMgr.RunCycle() {
				// Battle ended
				g.gameOver = true
				g.battleReport = g.battleMgr.GetBattleReport()
				log.Print(g.battleReport)
				break
			}
			g.cycle = g.battleMgr.stats.TotalCycles
		}
	}

	// Decay visual effects
	g.core.DecayEffects()

	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawCore(screen)
	g.DrawUI(screen)
}

// Layout returns the game's screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Parse command line flags
	mode := flag.String("mode", "visual", "Game mode: visual, battle, or tournament")
	warrior1 := flag.String("w1", "", "Path to first warrior file")
	warrior2 := flag.String("w2", "", "Path to second warrior file")
	rounds := flag.Int("rounds", 10, "Number of rounds for tournament mode")
	flag.Parse()

	// Load warriors based on mode
	var warriors []*Warrior

	switch *mode {
	case "visual":
		// Visual mode - interactive graphics
		if *warrior1 != "" && *warrior2 != "" {
			// Load specified warriors
			w1, err := LoadWarriorFromFile(*warrior1, Red)
			if err != nil {
				log.Fatalf("Error loading warrior 1: %v", err)
			}
			w2, err := LoadWarriorFromFile(*warrior2, Blue)
			if err != nil {
				log.Fatalf("Error loading warrior 2: %v", err)
			}
			warriors = []*Warrior{w1, w2}
		} else {
			// Use default warriors
			warriors = []*Warrior{CreateImp(), CreateDwarf()}
		}

		// Run visual game
		ebiten.SetWindowSize(screenWidth, screenHeight)
		ebiten.SetWindowTitle("Core War")

		game := NewGame(warriors)

		if err := ebiten.RunGame(game); err != nil {
			log.Fatal(err)
		}

	case "battle":
		// Battle mode - single battle with statistics
		if *warrior1 == "" || *warrior2 == "" {
			fmt.Println("Battle mode requires two warriors: -w1 <file> -w2 <file>")
			os.Exit(1)
		}

		w1, err := LoadWarriorFromFile(*warrior1, Red)
		if err != nil {
			log.Fatalf("Error loading warrior 1: %v", err)
		}
		w2, err := LoadWarriorFromFile(*warrior2, Blue)
		if err != nil {
			log.Fatalf("Error loading warrior 2: %v", err)
		}

		// Run battle
		bm := NewBattleManager(coreSize, maxCycles)
		bm.SetupBattle([]*Warrior{w1, w2})

		fmt.Printf("Battle: %s vs %s\n", w1.Name, w2.Name)
		fmt.Println("Running battle...")

		cycles := 0
		for bm.RunCycle() {
			cycles++
			if cycles%1000 == 0 {
				fmt.Printf(".")
			}
		}
		fmt.Println()

		fmt.Println(bm.GetBattleReport())

	case "tournament":
		// Tournament mode - round-robin tournament
		// Load all warriors from warriors directory or specified files
		allWarriors, err := LoadAllWarriors()
		if err != nil {
			log.Fatalf("Error loading warriors: %v", err)
		}

		if len(allWarriors) < 2 {
			// Use built-in warriors
			allWarriors = []*Warrior{
				CreateImp(),
				CreateDwarf(),
				CreateBomber(),
				CreateQuickScan(),
			}
		}

		// Run tournament
		tournament := NewTournament(allWarriors, *rounds, coreSize, maxCycles)
		tournament.Run()

	default:
		fmt.Printf("Unknown mode: %s\n", *mode)
		fmt.Println("Available modes: visual, battle, tournament")
		os.Exit(1)
	}
}
