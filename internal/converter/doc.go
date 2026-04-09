// Package converter provides format conversion utilities for .env file entries.
//
// Supported output formats:
//
//   - dotenv  — standard KEY=VALUE pairs
//   - export  — shell-compatible export KEY="VALUE" syntax
//   - json    — JSON object representation
//   - yaml    — YAML key-value mapping
//
// Example usage:
//
//	result, err := converter.Convert(entries, converter.FormatJSON)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Print(result)
package converter
