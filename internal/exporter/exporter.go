package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/masker"
)

// Format represents the output format for exported env entries.
type Format string

const (
	FormatDotEnv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatExport Format = "export"
)

// Options controls how the export is rendered.
type Options struct {
	Format  Format
	MaskSensitive bool
	Masker  *masker.Masker
}

// Write serialises env entries to w using the given options.
func Write(w io.Writer, entries []envfile.Entry, opts Options) error {
	switch opts.Format {
	case FormatJSON:
		return writeJSON(w, entries, opts)
	case FormatExport:
		return writeExport(w, entries, opts)
	default:
		return writeDotEnv(w, entries, opts)
	}
}

func applyMask(opts Options, key, value string) string {
	if opts.MaskSensitive && opts.Masker != nil && opts.Masker.IsSensitive(key) {
		return opts.Masker.Mask(value)
	}
	return value
}

func writeDotEnv(w io.Writer, entries []envfile.Entry, opts Options) error {
	for _, e := range entries {
		val := applyMask(opts, e.Key, e.Value)
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, val); err != nil {
			return err
		}
	}
	return nil
}

func writeExport(w io.Writer, entries []envfile.Entry, opts Options) error {
	for _, e := range entries {
		val := applyMask(opts, e.Key, e.Value)
		escaped := strings.ReplaceAll(val, "'", "'\\''")
		if _, err := fmt.Fprintf(w, "export %s='%s'\n", e.Key, escaped); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, entries []envfile.Entry, opts Options) error {
	m := make(map[string]string, len(entries))
	keys := make([]string, 0, len(entries))
	for _, e := range entries {
		m[e.Key] = applyMask(opts, e.Key, e.Value)
		keys = append(keys, e.Key)
	}
	sort.Strings(keys)
	ordered := make([]map[string]string, 0, len(keys))
	for _, k := range keys {
		ordered = append(ordered, map[string]string{"key": k, "value": m[k]})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(ordered)
}
