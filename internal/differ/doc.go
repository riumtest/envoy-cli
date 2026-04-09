// Package differ provides functionality for comparing environment variable files
// and generating human-readable diffs.
//
// The package supports:
//   - Comparing two environment variable maps
//   - Identifying added, removed, and changed variables
//   - Formatting differences in multiple output formats (text, JSON)
//   - Secret masking for sensitive values
//
// Example usage:
//
//	source := map[string]string{"KEY1": "value1"}
//	target := map[string]string{"KEY1": "value2", "KEY2": "new"}
//
//	result := differ.Compare(source, target)
//	formatter := &differ.TextFormatter{MaskSecrets: true}
//	fmt.Println(formatter.Format(result))
//
// Output:
//
//	Found 2 difference(s):
//
//	+ KEY2=new
//	~ KEY1
//	  - value1
//	  + value2
//
//	1 added, 0 removed, 1 changed
package differ
