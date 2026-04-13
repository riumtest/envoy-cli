// Package masker provides utilities for detecting and masking sensitive
// environment variable values based on key name patterns.
//
// A Masker can be created with default patterns (which cover common secrets
// such as passwords, tokens, and API keys) or with custom patterns.
//
// Example usage:
//
//	m := masker.New()
//	if m.IsSensitive("DATABASE_PASSWORD") {
//		fmt.Println(m.Mask("DATABASE_PASSWORD", "s3cr3t"))
//	}
package masker
