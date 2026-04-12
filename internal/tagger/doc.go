// Package tagger assigns user-defined tags to .env entries based on
// configurable prefix rules. Tagged entries can then be filtered,
// exported, or processed independently per tag group.
//
// Example usage:
//
//	opts := tagger.DefaultOptions()
//	opts.Rules = map[string]string{
//	    "database": "DB_",
//	    "cloud":    "AWS_",
//	}
//	tagged := tagger.Tag(entries, opts)
//	dbEntries := tagger.FilterByTag(tagged, "database")
package tagger
