// Package classifier categorises .env entries into logical groups
// based on their key names and value patterns (e.g. credentials, URLs,
// feature flags, infrastructure settings).
package classifier

import (
	"regexp"
	"strings"

	"github.com/envoy-cli/envoy/internal/envfile"
)

// Category represents a logical classification for an env entry.
type Category string

const (
	CategoryCredential  Category = "credential"
	CategoryURL         Category = "url"
	CategoryFeatureFlag Category = "feature_flag"
	CategoryInfra       Category = "infra"
	CategoryApp         Category = "app"
	CategoryUnknown     Category = "unknown"
)

// Result holds the classification outcome for a single entry.
type Result struct {
	Key      string
	Value    string
	Category Category
	Reason   string
}

// Summary holds aggregate counts per category.
type Summary struct {
	Total       int
	ByCategory  map[Category]int
	Uncategorised int
}

var (
	credentialPattern  = regexp.MustCompile(`(?i)(secret|password|passwd|token|api_key|apikey|auth|credential|private_key|passphrase)`)
	urlPattern         = regexp.MustCompile(`(?i)(url|uri|endpoint|host|addr|address|dsn|connection)`)
	featureFlagPattern = regexp.MustCompile(`(?i)(enable|disable|feature|flag|toggle|active)`)
	infraPattern       = regexp.MustCompile(`(?i)(port|host|region|zone|cluster|namespace|env|environment|stage|tier|replica|shard)`)
	boolValuePattern   = regexp.MustCompile(`(?i)^(true|false|yes|no|1|0|on|off)$`)
	urlValuePattern    = regexp.MustCompile(`(?i)^(https?|grpc|amqp|redis|postgres|mysql|mongodb|ftp)://`)
)

// Classify assigns a Category to each entry and returns the results.
func Classify(entries []envfile.Entry) []Result {
	results := make([]Result, 0, len(entries))
	for _, e := range entries {
		results = append(results, classifyEntry(e))
	}
	return results
}

// Summarise returns aggregate counts across all classified results.
func Summarise(results []Result) Summary {
	s := Summary{
		Total:      len(results),
		ByCategory: make(map[Category]int),
	}
	for _, r := range results {
		s.ByCategory[r.Category]++
		if r.Category == CategoryUnknown {
			s.UncategorisedCount()
		}
	}
	// count uncategorised
	s.Uncategorised = s.ByCategory[CategoryUnknown]
	return s
}

// ToMap converts a slice of Results into a map keyed by Category.
func ToMap(results []Result) map[Category][]Result {
	m := make(map[Category][]Result)
	for _, r := range results {
		m[r.Category] = append(m[r.Category], r)
	}
	return m
}

func classifyEntry(e envfile.Entry) Result {
	key := strings.ToUpper(e.Key)
	val := e.Value

	if credentialPattern.MatchString(key) {
		return Result{Key: e.Key, Value: val, Category: CategoryCredential, Reason: "key matches credential pattern"}
	}

	if urlValuePattern.MatchString(val) {
		return Result{Key: e.Key, Value: val, Category: CategoryURL, Reason: "value is a URL"}
	}

	if urlPattern.MatchString(key) {
		return Result{Key: e.Key, Value: val, Category: CategoryURL, Reason: "key matches URL pattern"}
	}

	if featureFlagPattern.MatchString(key) && boolValuePattern.MatchString(val) {
		return Result{Key: e.Key, Value: val, Category: CategoryFeatureFlag, Reason: "key matches feature flag pattern with boolean value"}
	}

	if infraPattern.MatchString(key) {
		return Result{Key: e.Key, Value: val, Category: CategoryInfra, Reason: "key matches infrastructure pattern"}
	}

	return Result{Key: e.Key, Value: val, Category: CategoryUnknown, Reason: "no pattern matched"}
}
