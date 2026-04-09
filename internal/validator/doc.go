// Package validator provides static analysis for .env files.
//
// It checks for common mistakes such as:
//   - Duplicate keys (warning)
//   - Empty values (warning)
//   - Empty or whitespace-containing keys (error)
//
// Usage:
//
//	entries := []validator.Entry{
//		{Line: 1, Key: "APP_ENV", Value: "production"},
//		{Line: 2, Key: "PORT", Value: ""},
//	}
//	result := validator.Validate(entries)
//	if result.HasErrors() {
//		fmt.Println("validation failed")
//	}
//	for _, issue := range result.Issues {
//		fmt.Println(issue)
//	}
package validator
