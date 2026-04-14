// Package tokenizer provides lexical tokenization of .env file entries.
//
// It converts a slice of envfile.Entry values into a flat sequence of
// typed Token values, each carrying a kind (key, equals, value, comment,
// or blank) and the original raw text.
//
// Tokens are useful for downstream features such as syntax highlighting,
// structural analysis, and custom formatters that need to distinguish
// between the different parts of an env file entry.
//
// Example:
//
//	tokens := tokenizer.Tokenize(entries)
//	keys   := tokenizer.FilterByKind(tokens, tokenizer.TokenKey)
package tokenizer
