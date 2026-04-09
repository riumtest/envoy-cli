// Package merger implements multi-file .env merging for envoy-cli.
//
// It supports combining entries from two or more parsed environment files
// into a single unified slice of key-value pairs. When the same key appears
// in more than one source file the caller chooses one of three conflict
// resolution strategies:
//
//   - StrategyFirst – keep the value from the earliest file in the merge
//     order (useful for "base overrides" patterns).
//   - StrategyLast  – overwrite with the value from the latest file
//     (useful for "layered environment" patterns).
//   - StrategyError – abort and return an error so the caller can surface
//     the conflict to the user.
//
// All conflicts are recorded in Result.Conflicts regardless of the chosen
// strategy, allowing callers to report or log them as needed.
package merger
