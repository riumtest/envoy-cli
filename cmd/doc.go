// Package cmd provides the CLI commands for envoy-cli.
//
// Commands:
//
//	envoy diff <base> <target> [flags]
//		Compare two .env files and print a human-readable or JSON diff.
//
// Flags for diff:
//
//	-m, --mask              Mask sensitive values (passwords, secrets, tokens)
//	-f, --format string     Output format: "text" (default) or "json"
//	-p, --pattern strings   Additional key patterns to treat as sensitive
//
// Example usage:
//
//	# Basic diff
//	envoy diff .env.staging .env.production
//
//	# Diff with masked secrets in JSON format
//	envoy diff --mask --format json .env.staging .env.production
//
//	# Diff with a custom sensitive pattern
//	envoy diff --mask --pattern "STRIPE_" .env.local .env.prod
package cmd
