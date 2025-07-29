package baranka

import (
	"strings"
)

// option defines a configuration function for Baranka.
type option func(*Baranka)

// getOptions applies a list of options to a new Baranka instance.
func getOptions(opts []option) *Baranka {
	cfg := New()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// WithLength pre-allocates the argument slice with the given length.
func WithLength(length int) option {
	return func(b *Baranka) {
		b.totalLength = length
		b.args = make([]any, 0, length)
	}
}

// PlaceholderFormat specifies the placeholder style for SQL queries.
type PlaceholderFormat byte

const (
	// PlaceholderFormatDollar uses $1, $2, ... (PostgreSQL style).
	PlaceholderFormatDollar PlaceholderFormat = iota
	// PlaceholderFormatQuestionMark uses ?, ?, ... (MySQL/SQLite style).
	PlaceholderFormatQuestionMark
)

// WithPlaceholderFormat sets the placeholder format for the Baranka instance.
func WithPlaceholderFormat(format PlaceholderFormat) option {
	return func(b *Baranka) {
		switch format {
		case PlaceholderFormatQuestionMark:
			b.placeholderFormat = format
		default:
			b.placeholderFormat = PlaceholderFormatDollar
		}
	}
}

// WithIncludeTemplate sets the template for value blocks (default: "(%s)").
func WithIncludeTemplate(template string) option {
	return func(b *Baranka) {
		if !strings.Contains(template, "%s") {
			template = "(%s)" // fallback to default
		}
		b.template = template
	}
}
