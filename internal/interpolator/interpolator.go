package interpolator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/envoy-cli/internal/envfile"
)

// varPattern matches ${VAR_NAME} and $VAR_NAME style references.
var varPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// Result holds the output of an interpolation pass.
type Result struct {
	Entries    []envfile.Entry
	Unresolved []string // keys whose values contained unresolvable references
}

// Interpolate expands variable references inside entry values using the
// provided entries as the source of truth. References to unknown variables
// are left as-is and the key is recorded in Result.Unresolved.
func Interpolate(entries []envfile.Entry) Result {
	lookup := buildLookup(entries)
	result := Result{}

	for _, e := range entries {
		expanded, unresolved := expand(e.Value, lookup)
		result.Entries = append(result.Entries, envfile.Entry{Key: e.Key, Value: expanded})
		if len(unresolved) > 0 {
			result.Unresolved = append(result.Unresolved, e.Key)
		}
	}
	return result
}

// InterpolateWithEnv expands variable references using both the provided
// entries and an additional environment map (e.g. os.Environ parsed into a map).
func InterpolateWithEnv(entries []envfile.Entry, env map[string]string) Result {
	lookup := buildLookup(entries)
	for k, v := range env {
		if _, exists := lookup[k]; !exists {
			lookup[k] = v
		}
	}
	return interpolateWithLookup(entries, lookup)
}

func interpolateWithLookup(entries []envfile.Entry, lookup map[string]string) Result {
	result := Result{}
	for _, e := range entries {
		expanded, unresolved := expand(e.Value, lookup)
		result.Entries = append(result.Entries, envfile.Entry{Key: e.Key, Value: expanded})
		if len(unresolved) > 0 {
			result.Unresolved = append(result.Unresolved, e.Key)
		}
	}
	return result
}

func expand(value string, lookup map[string]string) (string, []string) {
	var unresolved []string
	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		name := extractName(match)
		if v, ok := lookup[name]; ok {
			return v
		}
		unresolved = append(unresolved, name)
		return match
	})
	return result, unresolved
}

func extractName(match string) string {
	match = strings.TrimPrefix(match, "$")
	match = strings.TrimPrefix(match, "{")
	match = strings.TrimSuffix(match, "}")
	return match
}

func buildLookup(entries []envfile.Entry) map[string]string {
	lookup := make(map[string]string, len(entries))
	for _, e := range entries {
		lookup[e.Key] = e.Value
	}
	return lookup
}

// Summary returns a human-readable summary of the interpolation result.
func Summary(r Result) string {
	if len(r.Unresolved) == 0 {
		return fmt.Sprintf("interpolated %d entries, all references resolved", len(r.Entries))
	}
	return fmt.Sprintf("interpolated %d entries, %d unresolved: %s",
		len(r.Entries), len(r.Unresolved), strings.Join(r.Unresolved, ", "))
}
