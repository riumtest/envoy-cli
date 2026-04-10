// Package renamer provides functionality for bulk-renaming keys within a set
// of environment file entries.
//
// A Rule maps a source key (From) to a destination key (To). Rename processes
// all rules against the supplied entries and returns:
//
//   - The updated entry slice with keys renamed in-place.
//   - A Result describing which rules succeeded, were skipped (source key
//     absent), or caused a conflict (destination key already exists).
//
// The original slice is never mutated; a copy is returned.
package renamer
