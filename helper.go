package baranka

import (
	"fmt"
	"strconv"
	"strings"
)

type Baranka struct {
	count             int
	totalLength       int
	args              []any
	placeholderFormat PlaceholderFormat
	template          string
	blocks            []string
}

// NewBaranka creates a new Baranka helper with default settings.
func NewBaranka() *Baranka {
	return &Baranka{
		count:    1,
		template: "(%s)",
	}
}

// Add appends a new block of arguments and generates corresponding placeholders.
func (b *Baranka) Add(newArgs ...any) {
	if len(newArgs) == 0 {
		return
	}

	b.args = append(b.args, newArgs...)
	if b.blocks == nil {
		b.blocks = make([]string, 0, b.totalLength/len(newArgs))
	}
	b.blocks = append(b.blocks, fmt.Sprintf(b.template, strings.Join(b.getPlaceholders(len(newArgs)), ",")))
}

func (b *Baranka) getPlaceholders(increment int) []string {
	var placeholders = make([]string, 0, increment)
	for i := 0; i < increment; i++ {
		placeholders = append(placeholders, b.getPlaceholder())
	}
	return placeholders
}

func (b *Baranka) getPlaceholder() string {
	defer func() {
		b.count++
	}()

	switch b.placeholderFormat {
	case PlaceholderFormatQuestionMark:
		return "?"
	default:
		return "$" + strconv.Itoa(b.count)
	}
}

// Args returns the collected arguments in order.
func (b *Baranka) Args() []any {
	return b.args
}

// Values returns the SQL value blocks as a string, e.g. ($1,$2),\n($3,$4).
func (b *Baranka) Values() string {
	return strings.Join(b.blocks, ",\n")
}
