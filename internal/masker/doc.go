// Package masker provides utilities for detecting and masking sensitive
// environment variable values based on key name patterns.
//
// A Masker can be created with default patterns (e.g. keys containing
// "SECRET", "PASSWORD", "TOKEN", "KEY", "AUTH") or with additional
// custom patterns supplied by the caller.
//
// Example usage:
//
//	m := masker.New()
//	if m.IsSensitive("DB_PASSWORD") {
//		fmt.Println(m.Mask("DB_PASSWORD", "hunter2")) // ***
//	}
package masker
