// Package interpolator expands variable references within .env file values.
//
// It supports two reference styles:
//
//	${VAR_NAME}   — brace-delimited reference
//	$VAR_NAME     — bare dollar reference
//
// References are resolved using the entries themselves (self-referential
// expansion) and, optionally, an additional environment map supplied by the
// caller (e.g. the host process environment).
//
// Entries are processed in declaration order. Chained references that depend
// on values produced by earlier expansions in the same pass are not
// automatically re-expanded; callers that need multi-pass expansion should
// call Interpolate repeatedly until Result.Unresolved is empty.
package interpolator
