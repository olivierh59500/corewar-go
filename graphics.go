package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// DrawCore renders the memory core visualization
func (g *Game) DrawCore(screen *ebiten.Image) {
	// Calculate optimal grid dimensions to fill the screen
	margin := 10
	availableWidth := screenWidth - 2*margin
	availableHeight := screenHeight - infoHeight - 2*margin

	// Calculate optimal cell size
	// Try different column counts to find the best fit
	bestCellSize := 1
	bestCols := 1

	for cols := 32; cols <= 200; cols++ {
		rows := (coreSize + cols - 1) / cols

		// Calculate cell size that would fit
		maxCellWidth := availableWidth / cols
		maxCellHeight := availableHeight / rows

		// Account for spacing
		if cols > 1 {
			maxCellWidth = (availableWidth - (cols-1)*cellSpacing) / cols
		}
		if rows > 1 {
			maxCellHeight = (availableHeight - (rows-1)*cellSpacing) / rows
		}

		// Use the smaller of the two to ensure it fits
		cellSize := maxCellWidth
		if maxCellHeight < cellSize {
			cellSize = maxCellHeight
		}

		// Check if this is better than what we have
		if cellSize >= bestCellSize && cellSize >= 2 {
			bestCellSize = cellSize
			bestCols = cols
		}
	}

	// Use the calculated values
	actualCellSize := bestCellSize
	actualCols := bestCols
	gridRows := (coreSize + actualCols - 1) / actualCols

	// Draw background
	vector.DrawFilledRect(screen, 0, 0, float32(screenWidth), float32(screenHeight), color.RGBA{20, 20, 20, 255}, false)

	// Center the grid
	totalWidth := actualCols*actualCellSize + (actualCols-1)*cellSpacing
	totalHeight := gridRows*actualCellSize + (gridRows-1)*cellSpacing
	offsetX := (screenWidth - totalWidth) / 2
	offsetY := (screenHeight - infoHeight - totalHeight) / 2

	// Draw memory cells
	for i := 0; i < coreSize; i++ {
		row := i / actualCols
		col := i % actualCols

		x := float32(offsetX + col*(actualCellSize+cellSpacing))
		y := float32(offsetY + row*(actualCellSize+cellSpacing))

		// Skip if outside screen bounds
		if y > float32(screenHeight-infoHeight) {
			continue
		}

		// Get base color
		baseColor := GetWarriorColor(g.core.owners[i])

		// Check if it's a DAT instruction (bomb)
		isDat := g.core.cells[i].Op == DAT && g.core.owners[i] != Empty

		// Apply effects
		r, g_, b := float32(baseColor.R), float32(baseColor.G), float32(baseColor.B)

		// Make DAT instructions darker to distinguish them
		if isDat {
			r = r * 0.5
			g_ = g_ * 0.5
			b = b * 0.5
		}

		// Execution effect (brighten)
		if g.core.execEffect[i] > 0 {
			factor := 1 + g.core.execEffect[i]
			r = min(255, r*factor)
			g_ = min(255, g_*factor)
			b = min(255, b*factor)
		}

		// Write effect (yellow tint)
		if g.core.writeEffect[i] > 0 {
			r = min(255, r+200*g.core.writeEffect[i])
			g_ = min(255, g_+200*g.core.writeEffect[i])
		}

		// Read effect (green tint)
		if g.core.readEffect[i] > 0 {
			g_ = min(255, g_+200*g.core.readEffect[i])
		}

		cellColor := color.RGBA{uint8(r), uint8(g_), uint8(b), 255}

		// Draw cell
		vector.DrawFilledRect(screen, x, y, float32(actualCellSize), float32(actualCellSize), cellColor, false)
	}

	// Draw grid info
	infoY := float32(screenHeight - infoHeight + 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Core Size: %d cells (%d x %d grid)",
		coreSize, actualCols, gridRows), 10, int(infoY))

	// Draw cell ownership stats
	cellCounts := make(map[WarriorColor]int)
	for i := 0; i < coreSize; i++ {
		cellCounts[g.core.owners[i]]++
	}

	x := 300
	for _, w := range g.warriors {
		count := cellCounts[w.Color]
		percentage := float64(count) / float64(coreSize) * 100
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s cells: %d (%.1f%%)",
			w.Name, count, percentage), x, int(infoY))
		x += 200
	}
}

// DrawUI renders the user interface
func (g *Game) DrawUI(screen *ebiten.Image) {
	// Draw info panel background
	infoY := float32(screenHeight - infoHeight)
	vector.DrawFilledRect(screen, 0, infoY, float32(screenWidth), float32(infoHeight),
		color.RGBA{30, 30, 30, 255}, false)

	// Draw separator line
	vector.DrawFilledRect(screen, 0, infoY, float32(screenWidth), 2,
		color.RGBA{100, 100, 100, 255}, false)

	// Draw warrior info
	y := int(infoY) + 30
	for _, w := range g.warriors {
		colorBox := GetWarriorColor(w.Color)

		// Draw color indicator
		vector.DrawFilledRect(screen, 10, float32(y), 20, 20, colorBox, false)

		// Count processes
		processCount := 0
		for _, p := range g.vm.processes {
			if p.warrior == w && p.alive {
				processCount++
			}
		}

		// Draw warrior info
		status := "ALIVE"
		if processCount == 0 {
			status = "DEAD"
		}

		// Count DAT bombs placed by this warrior
		datCount := 0
		for j := 0; j < g.core.size; j++ {
			if g.core.owners[j] == w.Color && g.core.cells[j].Op == DAT {
				datCount++
			}
		}

		info := fmt.Sprintf("%s - %s (Processes: %d, DATs: %d)", w.Name, status, processCount, datCount)
		ebitenutil.DebugPrintAt(screen, info, 40, y+5)

		y += 25
	}

	// Draw controls
	controlsX := screenWidth / 2
	y = int(infoY) + 10

	pauseText := "RUNNING"
	if g.paused {
		pauseText = "PAUSED"
	}
	if g.gameOver {
		pauseText = "GAME OVER"
		if g.battleMgr.stats.Winner != nil {
			pauseText = fmt.Sprintf("WINNER: %s", g.battleMgr.stats.Winner.Name)
		} else if g.battleMgr.stats.IsDraw {
			pauseText = "DRAW"
		}
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Status: %s", pauseText), controlsX, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Speed: %d cycles/frame", g.speed), controlsX, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Cycle: %d / %d", g.cycle, maxCycles), controlsX, y)

	// Draw control instructions
	instructionsX := screenWidth - 300
	y = int(infoY) + 10

	ebitenutil.DebugPrintAt(screen, "Controls:", instructionsX, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, "SPACE - Pause/Resume", instructionsX, y)
	y += 15
	ebitenutil.DebugPrintAt(screen, "↑/↓ - Adjust Speed", instructionsX, y)
	y += 15
	ebitenutil.DebugPrintAt(screen, "R - Restart", instructionsX, y)
	y += 15
	ebitenutil.DebugPrintAt(screen, "ESC - Exit", instructionsX, y)
}

// min returns the minimum of two float32 values
func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
