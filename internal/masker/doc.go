// Package masker provides key-based secret detection and value masking
// for environment variable processing in envoy-cli.
//
// A Masker inspects environment variable keys against a configurable list
// of sensitive patterns (e.g. "SECRET", "PASSWORD", "TOKEN") and replaces
// matching values with a redaction string before they are displayed or
// written to output.
//
// Basic usage:
//
//	m := masker.New()
//
//	// Check whether a key is sensitive.
//	if m.IsSensitive("DB_PASSWORD") {
//		fmt.Println("this key is sensitive")
//	}
//
//	// Mask a value — returns "****" for sensitive keys, original otherwise.
//	display := m.Mask("DB_PASSWORD", "s3cr3t") // → "****"
//	display  = m.Mask("APP_ENV",    "prod")    // → "prod"
//
// Custom patterns and mask strings can be supplied via NewWithPatterns.
package masker
