// Package inspector analyses a collection of parsed env entries and
// produces a structured Report containing statistics such as the total
// number of keys, how many are empty, how many appear sensitive (via
// the masker package), and simple value-type heuristics (URLs, booleans,
// and numeric strings).
//
// Example usage:
//
//	m := masker.New()
//	entries := []inspector.Entry{
//		{Key: "DATABASE_URL", Value: "postgres://localhost/mydb"},
//		{Key: "DEBUG",        Value: "true"},
//		{Key: "API_SECRET",   Value: "s3cr3t"},
//	}
//	report := inspector.Inspect(entries, m)
//	fmt.Printf("Total: %d, Sensitive: %d\n", report.Total, report.Sensitive)
package inspector
