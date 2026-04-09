package envfile

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry represents a single environment variable entry
type Entry struct {
	Key     string
	Value   string
	Comment string
	IsEmpty bool // true for empty lines or comment-only lines
}

// EnvFile represents a parsed .env file
type EnvFile struct {
	Entries []Entry
	Path    string
}

// Parse reads and parses an .env file from a reader
func Parse(r io.Reader, path string) (*EnvFile, error) {
	scanner := bufio.NewScanner(r)
	var entries []Entry

	for scanner.Scan() {
		line := scanner.Text()
		entry := parseLine(line)
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return &EnvFile{
		Entries: entries,
		Path:    path,
	}, nil
}

// parseLine parses a single line from an .env file
func parseLine(line string) Entry {
	trimmed := strings.TrimSpace(line)

	// Empty line
	if trimmed == "" {
		return Entry{IsEmpty: true}
	}

	// Comment line
	if strings.HasPrefix(trimmed, "#") {
		return Entry{Comment: trimmed, IsEmpty: true}
	}

	// Key-value pair with optional inline comment
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		// Invalid format, treat as comment
		return Entry{Comment: line, IsEmpty: true}
	}

	key := strings.TrimSpace(parts[0])
	valuePart := parts[1]

	// Check for inline comment
	value := valuePart
	comment := ""
	if idx := strings.Index(valuePart, " #"); idx != -1 {
		value = strings.TrimSpace(valuePart[:idx])
		comment = strings.TrimSpace(valuePart[idx:])
	} else {
		value = strings.TrimSpace(valuePart)
	}

	// Remove quotes if present
	value = strings.Trim(value, `"`)
	value = strings.Trim(value, `'`)

	return Entry{
		Key:     key,
		Value:   value,
		Comment: comment,
		IsEmpty: false,
	}
}

// ToMap converts the EnvFile to a map of key-value pairs
func (e *EnvFile) ToMap() map[string]string {
	result := make(map[string]string)
	for _, entry := range e.Entries {
		if !entry.IsEmpty && entry.Key != "" {
			result[entry.Key] = entry.Value
		}
	}
	return result
}
