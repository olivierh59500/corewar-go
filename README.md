# Core War Game

A Go implementation of the classic Core War game with graphical visualization using Ebiten.

## Overview

Core War is a programming game where programs (called "warriors") battle in a virtual computer's memory. Each warrior tries to eliminate the others by causing them to execute invalid instructions. This implementation features:

- Real-time graphical visualization of the memory core
- Support for Redcode assembly language
- Multiple example warriors
- Pause/resume functionality
- Speed control
- Visual representation of read/write operations
- Battle statistics and tournament mode
- Three game modes: Visual, Battle, and Tournament

## Features

- **Visual Memory Core**: Watch programs battle in real-time with color-coded visualization
- **Redcode Assembler**: Built-in assembler for the Redcode language
- **Multiple Warriors**: Includes several classic warriors (Imp, Dwarf, Vampire, Scanner, etc.)
- **Interactive Controls**: 
  - Space: Pause/Resume
  - Up/Down: Adjust execution speed
  - R: Restart battle
  - ESC: Exit
- **Game Modes**:
  - **Visual Mode**: Interactive graphical battle viewer
  - **Battle Mode**: Single battle with detailed statistics
  - **Tournament Mode**: Round-robin tournament between multiple warriors

## Requirements

- Go 1.19 or higher
- Ebiten v2

## Installation

```bash
# Clone the repository
git clone https://github.com/olivierh59500/corewar-go
cd corewar-go

# Install dependencies
go mod init corewar
go get github.com/hajimehoshi/ebiten/v2

# Run the game
go run .
```

## Usage

### Visual Mode (Default)
Watch battles with real-time graphics:
```bash
# Run with default warriors
go run .

# Run with custom warriors
go run . -w1 warriors/imp.red -w2 warriors/dwarf.red
```

### Battle Mode
Run a single battle with statistics:
```bash
go run . -mode battle -w1 warriors/vampire.red -w2 warriors/scanner.red
```

### Tournament Mode
Run a round-robin tournament:
```bash
go run . -mode tournament -rounds 20
```

## Project Structure

```
corewar-go/
├── main.go           # Entry point and game loop
├── core.go           # Memory core implementation
├── vm.go             # Virtual machine (MARS)
├── assembler.go      # Redcode assembler
├── warrior.go        # Warrior structure and loading
├── graphics.go       # Ebiten graphics rendering
├── battle.go         # Battle manager and statistics
├── loader.go         # File loading utilities
├── warriors/         # Example warrior programs
│   ├── imp.red
│   ├── dwarf.red
│   ├── mice.red
│   ├── scanner.red
│   ├── vampire.red
│   └── gate.red
├── go.mod
├── go.sum
├── run.sh           # Helper script
└── README.md
```

## How to Play

1. Run the game with `go run .`
2. The game will automatically load two warriors and start the battle
3. Watch as the warriors execute instructions and try to eliminate each other
4. Use keyboard controls to interact with the simulation

## Understanding the Display

- **Memory Grid**: Each cell represents a memory location (8000 total)
- **Colors**:
  - Gray: Empty memory
  - Red/Blue/Green/Yellow: Instructions belonging to different warriors
  - Bright colors: Currently executing instruction
  - Yellow tint: Recent write operation
  - Green tint: Recent read operation

## Writing Warriors

Warriors are written in Redcode, a simplified assembly language. Here's a simple example:

```redcode
;redcode
;name MyWarrior
;author Your Name
;strategy Description of strategy

start:  MOV bomb, @ptr    ; Copy bomb to target
        ADD #4, ptr       ; Increment pointer
        JMP start         ; Loop
bomb:   DAT #0, #0        ; The bomb
ptr:    DAT #0, #100      ; Starting target

end start
```

Save your warrior as a `.red` file in the `warriors/` directory.

## Redcode Instructions

- **MOV**: Copy data from source to destination
- **ADD**: Add source to destination
- **SUB**: Subtract source from destination
- **JMP**: Jump to address
- **JMZ**: Jump if zero
- **JMN**: Jump if not zero
- **DJN**: Decrement and jump if not zero
- **CMP/SEQ**: Compare and skip next instruction if equal
- **SPL**: Split (create new process)
- **DAT**: Data (terminates execution)
- **NOP**: No operation

## Addressing Modes

- `#`: Immediate (e.g., `#5`) - Use the number itself
- `: Direct (e.g., `$5`) - Use the address
- `@`: Indirect (e.g., `@5`) - Use the address pointed to
- `<`: Pre-decrement indirect - Decrement pointer before use
- `>`: Post-increment indirect - Increment pointer after use

## Example Warriors

### Imp
The simplest warrior - copies itself through memory:
```redcode
MOV 0, 1
```

### Dwarf
Bombs memory with DAT instructions:
```redcode
ADD #4, 3
MOV #0, @-1
JMP -2
DAT #0, #0
```

### Vampire
Converts enemy processes by making them jump to a trap:
```redcode
pit: JMP pit
MOV pit, @fang
ADD #10, fang
JMP -2
fang: DAT #0, #10
```

## Battle Statistics

In battle and tournament modes, the game tracks:
- Total cycles executed
- Maximum processes per warrior
- Instructions executed per warrior
- Efficiency ratings
- Win/loss records (tournament mode)

## Contributing

Feel free to contribute by:
- Adding new warriors
- Improving the graphics
- Optimizing the VM
- Adding new features
- Implementing more Redcode instructions

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Original Core War concept by A.K. Dewdney
- Inspired by the ICWS (International Core War Society) standards
- Redcode specification based on ICWS '94 standard