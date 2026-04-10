// Package patcher applies declarative patches to a slice of env entries.
//
// Supported operations:
//
//   - set:    create or overwrite a key with a given value
//   - delete: remove a key if it exists
//   - rename: change the name of an existing key
//
// Example usage:
//
//	patches := []patcher.Patch{
//		{Op: patcher.OpSet,    Key: "APP_ENV",  Value: "staging"},
//		{Op: patcher.OpDelete, Key: "DEBUG"},
//		{Op: patcher.OpRename, Key: "DB_URL",   NewKey: "DATABASE_URL"},
//	}
//	result, err := patcher.Apply(entries, patches)
package patcher
