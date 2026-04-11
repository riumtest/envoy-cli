// Package flattener provides utilities for namespacing .env keys by adding
// or removing a common prefix.
//
// Flattening is useful when merging variables from multiple services into a
// single file: each service's keys are prefixed with a namespace so they
// remain distinct. Unflattening reverses the process, stripping a known
// prefix to restore the original key names.
//
// Example:
//
//	DB_HOST=localhost  →  flatten("DB")  →  DB_HOST=localhost
//	DB_HOST=localhost  →  unflatten("DB") →  HOST=localhost
package flattener
