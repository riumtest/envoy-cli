// Package linker provides cross-file environment variable linking.
//
// It allows a destination .env file to pull values from one or more source
// files according to a set of declarative rules. Each rule specifies the
// source file path, the key to read, and the target key to write into the
// destination entry list.
//
// This is useful for environment promotion workflows where a staging or
// production file needs to inherit specific secrets from a vault-managed
// source without duplicating them by hand.
package linker
