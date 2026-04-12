// Package auditor analyses a set of .env entries and produces a structured
// audit Report containing Findings at three severity levels:
//
//   - error   – something is definitively wrong (e.g. empty key)
//   - warning – something looks suspicious (e.g. empty value, short secret)
//   - info    – an observation worth noting (e.g. TODO marker, underscore prefix)
//
// Usage:
//
//	report := auditor.Audit(entries, maskerInstance)
//	for _, f := range report.Findings {
//		fmt.Printf("[%s] %s: %s\n", f.Severity, f.Key, f.Message)
//	}
package auditor
