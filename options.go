package baranka

import (
	"strings"
)

// option defines a configuration function for Baranka.
type option func(*Baranka)

// applyOptions applies a list of options to a Baranka instance.
func (b *Baranka) applyOptions(opts []option) {
	for _, opt := range opts {
		opt(b)
	}
}

// WithBlocksLength sets the expected number of blocks for the Baranka instance.
func WithBlocksLength(length int) option {
	return func(b *Baranka) {
		b.expectedBlocks = length
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
