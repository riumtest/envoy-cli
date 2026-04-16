// Package pitcher promotes (pitches) environment entries from a source
// env file into a destination env file.
//
// It supports:
//   - Selective key promotion via an allowlist
//   - Optional key prefixing on the destination
//   - Overwrite control for existing destination keys
//
// Example usage:
//
//	result, err := pitcher.Pitch(srcEntries, dstEntries, pitcher.Options{
//		Keys:      []string{"DB_HOST", "DB_PORT"},
//		Overwrite: true,
//		Prefix:    "PROD",
//	})
package pitcher
