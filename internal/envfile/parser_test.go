package envfile

import (
	"strings"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Entry
	}{
		{
			name:     "simple key-value",
			input:    "DB_HOST=localhost",
			expected: Entry{Key: "DB_HOST", Value: "localhost", IsEmpty: false},
		},
		{
			name:     "quoted value",
			input:    `API_KEY="secret123"`,
			expected: Entry{Key: "API_KEY", Value: "secret123", IsEmpty: false},
		},
		{
			name:     "empty line",
			input:    "",
			expected: Entry{IsEmpty: true},
		},
		{
			name:     "comment line",
			input:    "# This is a comment",
			expected: Entry{Comment: "# This is a comment", IsEmpty: true},
		},
		{
			name:     "inline comment",
			input:    "PORT=3000 # Application port",
			expected: Entry{Key: "PORT", Value: "3000", Comment: "# Application port", IsEmpty: false},
		},
		{
			name:     "value with spaces",
			input:    "MESSAGE=Hello World",
			expected: Entry{Key: "MESSAGE", Value: "Hello World", IsEmpty: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLine(tt.input)
			if result.Key != tt.expected.Key || result.Value != tt.expected.Value ||
				result.IsEmpty != tt.expected.IsEmpty {
				t.Errorf("parseLine(%q) = %+v, want %+v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParse(t *testing.T) {
	input := `# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin

# API Keys
API_KEY="secret123"
API_SECRET=xyz789 # Keep this safe
`

	reader := strings.NewReader(input)
	envFile, err := Parse(reader, ".env")

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if envFile.Path != ".env" {
		t.Errorf("Path = %v, want .env", envFile.Path)
	}

	if len(envFile.Entries) != 8 {
		t.Errorf("Got %d entries, want 8", len(envFile.Entries))
	}

	envMap := envFile.ToMap()
	if envMap["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST = %v, want localhost", envMap["DB_HOST"])
	}

	if envMap["API_KEY"] != "secret123" {
		t.Errorf("API_KEY = %v, want secret123", envMap["API_KEY"])
	}

	if len(envMap) != 5 {
		t.Errorf("Map has %d entries, want 5", len(envMap))
	}
}
