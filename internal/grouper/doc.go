// Package grouper provides utilities for grouping environment variable entries
// by a common key prefix.
//
// Grouping is determined by splitting each key on a configurable delimiter
// (default "_") and using the first segment as the group name. Keys that
// contain no delimiter are placed in a special "_other" group.
//
// Example usage:
//
//	groups := grouper.GroupByPrefix(entries, grouper.DefaultOptions())
//	for _, g := range groups {
//		fmt.Printf("[%s] %d keys\n", g.Name, len(g.Entries))
//	}
package grouper
