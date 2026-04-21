// Package redactor provides utilities for sanitising environment variable
// entries by replacing sensitive values with a configurable placeholder.
//
// It builds on top of the masker package to identify sensitive keys and
// supports both structured ([]envfile.Entry) and raw string redaction.
//
// Basic usage:
//
//	entries := []envfile.Entry{{Key: "DB_PASSWORD", Value: "s3cr3t"}, ...}
//	redacted := redactor.Redact(entries, redactor.Options{})
//
// Custom placeholder and extra patterns:
//
//	opts := redactor.Options{
//		Placeholder:   "<REDACTED>",
//		ExtraPatterns: []string{"token", "cert"},
//	}
//	redacted := redactor.Redact(entries, opts)
//
// Raw string redaction:
//
//	// RedactString replaces sensitive values found in a plain string,
//	// which is useful for sanitising log output or error messages that
//	// may inadvertently contain secret values.
//	clean := redactor.RedactString(rawOutput, entries, redactor.Options{})
package redactor
