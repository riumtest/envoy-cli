// Package cloner implements env entry cloning between two sets of parsed entries.
//
// It supports:
//   - Full or partial key cloning via an allowlist
//   - Overwrite control to protect existing destination keys
//   - A Result summary reporting cloned, skipped, and conflicted counts
//
// Example usage:
//
//	res, err := cloner.Clone(dst, src, cloner.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Cloned %d keys\n", res.Cloned)
package cloner
