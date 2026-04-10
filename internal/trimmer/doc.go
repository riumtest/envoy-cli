// Package trimmer provides utilities for cleaning up parsed .env entries.
//
// It supports three orthogonal operations that can be combined via Options:
//
//   - TrimValues: strips leading and trailing whitespace from every value.
//   - RemoveEmpty: discards entries whose value is the empty string (after
//     optional trimming).
//   - DeduplicateKeys: when the same key appears more than once, the last
//     occurrence wins and earlier ones are removed, preserving order.
//
// The original slice is never mutated; Trim always returns a fresh copy.
package trimmer
