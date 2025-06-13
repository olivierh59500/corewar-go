package main

import (
	"os"
	"path/filepath"
	"strings"
)

// LoadWarriorFromFile loads a warrior from a .red file
func LoadWarriorFromFile(filename string, color WarriorColor) (*Warrior, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Extract name from metadata or filename
	name := extractName(string(content))
	if name == "" {
		name = strings.TrimSuffix(filepath.Base(filename), ".red")
	}

	// Extract author from metadata
	author := extractAuthor(string(content))

	warrior, err := LoadWarriorFromSource(name, string(content), color)
	if err != nil {
		return nil, err
	}

	warrior.Author = author
	return warrior, nil
}

// extractName extracts the warrior name from metadata comments
func extractName(source string) string {
	lines := strings.Split(source, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ";name ") {
			return strings.TrimSpace(line[6:])
		}
	}
	return ""
}

// extractAuthor extracts the author from metadata comments
func extractAuthor(source string) string {
	lines := strings.Split(source, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ";author ") {
			return strings.TrimSpace(line[8:])
		}
	}
	return "Unknown"
}

// LoadAllWarriors loads all warriors from the warriors directory
func LoadAllWarriors() ([]*Warrior, error) {
	warriors := make([]*Warrior, 0)
	colors := []WarriorColor{Red, Blue, Green, Yellow}
	colorIndex := 0

	files, err := filepath.Glob("warriors/*.red")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		warrior, err := LoadWarriorFromFile(file, colors[colorIndex%len(colors)])
		if err != nil {
			// Log error but continue loading other warriors
			continue
		}

		warriors = append(warriors, warrior)
		colorIndex++

		// Limit to 4 warriors for visualization
		if len(warriors) >= 4 {
			break
		}
	}

	return warriors, nil
}
