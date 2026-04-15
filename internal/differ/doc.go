// Package differ provides functionality for comparing two sets of environment
// file entries and producing a structured diff result.
//
// It identifies keys that have been added, removed, changed, or remain
// unchanged between a base and target environment. Results can be summarised
// by change type and formatted via the formatter sub-package.
//
// Example usage:
//
//	base := []envfile.Entry{{Key: "FOO", Value: "bar"}}
//	target := []envfile.Entry{{Key: "FOO", Value: "baz"}, {Key: "NEW", Value: "val"}}
//
//	result := differ.Compare(base, target)
//	summary := result.Summary()
//	// summary[differ.Changed] == 1
//	// summary[differ.Added]   == 1
package differ
