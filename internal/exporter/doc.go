// Package exporter provides functionality for serialising parsed .env entries
// into various output formats.
//
// Supported formats:
//
//   - dotenv  – standard KEY=VALUE pairs, one per line (default)
//   - export  – shell-compatible `export KEY='VALUE'` statements
//   - json    – a JSON array of {"key", "value"} objects, sorted by key
//
// Sensitive values (as determined by the masker package) can be automatically
// redacted before output by setting Options.MaskSensitive to true and
// supplying a masker.Masker instance.
//
// Example usage:
//
//	entries, _ := envfile.Parse(r)
//	m := masker.New()
//	exporter.Write(os.Stdout, entries, exporter.Options{
//		Format:        exporter.FormatExport,
//		MaskSensitive: true,
//		Masker:        m,
//	})
package exporter
