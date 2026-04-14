// Package masker provides utilities for detecting and masking sensitive
// environment variable values based on key name patterns.
//
// Sensitive keys are identified by matching against a set of configurable
// patterns (e.g. keys containing "SECRET", "PASSWORD", "TOKEN", etc.).
// Matched values are replaced with a fixed placeholder string.
//
// Example usage:
//
//	m := masker.New()
//	if m.IsSensitive("DB_PASSWORD") {
//		fmt.Println(m.Mask("DB_PASSWORD", "hunter2")) // "***"
//	}
package masker
