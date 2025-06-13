#!/bin/bash

# Core War Game Runner Script

echo "Core War Game - Go Implementation"
echo "================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.19 or higher."
    exit 1
fi

# Initialize module if needed
if [ ! -f "go.mod" ]; then
    echo "Initializing Go module..."
    go mod init corewar
fi

# Download dependencies
echo "Downloading dependencies..."
go mod download
go mod tidy

# Create warriors directory if it doesn't exist
if [ ! -d "warriors" ]; then
    echo "Creating warriors directory..."
    mkdir warriors
fi

# Create warrior files if they don't exist
if [ ! -f "warriors/imp.red" ]; then
    echo "Creating example warrior files..."
    
    cat > warriors/imp.red << 'EOF'
;redcode
;name Imp
;author A.K. Dewdney
;strategy Copies itself through core memory

imp:    MOV 0, 1

end
EOF

    cat > warriors/dwarf.red << 'EOF'
;redcode
;name Dwarf
;author A.K. Dewdney
;strategy Bombs core with DAT instructions

        ADD #4, 3
        MOV #0, @-1
        JMP -2
        DAT #0, #0

end
EOF

    cat > warriors/vampire.red << 'EOF'
;redcode
;name Vampire
;author Core War Community
;strategy Converts enemy processes

pit:    JMP pit
        MOV pit, @fang
        ADD #10, fang
        JMP -2
fang:   DAT #0, #10

end
EOF

    cat > warriors/scanner.red << 'EOF'
;redcode
;name Scanner
;author Core War Community
;strategy Scans and bombs

scan:   CMP -10, -20
        JMP found
        ADD #1, -2
        JMP scan
found:  MOV #0, @-4
        JMP scan

end
EOF
fi

# Run the game
echo "Starting Core War..."
go run .