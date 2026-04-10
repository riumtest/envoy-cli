// Package comparator provides utilities for performing a deep structural
// comparison between two sets of environment file entries.
//
// Unlike the differ package which focuses on change detection for display,
// the comparator package produces a structured Result that can be used
// programmatically — for example to compute overlap ratios, detect
// configuration drift, or drive merge decisions.
//
// Example usage:
//
//	left, _ := loader.Load("staging.env")
//	right, _ := loader.Load("production.env")
//	result := comparator.Compare(left, right)
//	fmt.Printf("Overlap: %.0f%%\n", comparator.OverlapRatio(result)*100)
package comparator
