// Package sorter provides sorting utilities for collections of environment
// variable entries.
//
// Entries can be sorted by key or value, in ascending or descending order.
// All sort operations are non-destructive — the original slice is never
// modified.
//
// Example usage:
//
//	result := sorter.SortByKey(entries)
//
//	opts := sorter.Options{Order: sorter.Descending, ByValue: false}
//	result := sorter.Sort(entries, opts)
package sorter
