package baranka

import (
	"fmt"
	"strconv"
	"strings"
)

type Baranka struct {
	count             int
	expectedBlocks    int
	args              []any
	placeholderFormat PlaceholderFormat
	template          string
	blocks            []string
}

// New creates a new Baranka.
func New(opts ...option) *Baranka {
	b := Baranka{
		count:    1,
		template: "(%s)",
	}

	b.applyOptions(opts)

	return &b
}

// Add appends a new block of arguments and generates corresponding placeholders.
func (b *Baranka) Add(newArgs ...any) {
	if len(newArgs) == 0 {
		return
	}

	extractedArgs := extractArgs(newArgs)

	if b.args == nil {
		b.args = make([]any, 0, b.expectedBlocks*len(extractedArgs))
	}
	b.args = append(b.args, extractedArgs...)

	if b.blocks == nil {
		b.blocks = make([]string, 0, b.expectedBlocks)
	}
	b.blocks = append(b.blocks, fmt.Sprintf(b.template, strings.Join(b.getPlaceholders(newArgs), ",")))
}

func (b *Baranka) Reset() {
	b.count = 1
	b.args = nil
	b.blocks = nil
}

func extractArgs(args []any) []any {
	capacity := 0
	for _, arg := range args {
		switch typedArg := arg.(type) {
		case Expression:
			capacity += len(typedArg.args)
		default:
			capacity++
		}
	}

	if capacity == len(args) {
		return args
	}

	extractedArgs := make([]any, 0, capacity)
	for _, arg := range args {
		switch typedArg := arg.(type) {
		case Expression:
			extractedArgs = append(extractedArgs, typedArg.args...)
		default:
			extractedArgs = append(extractedArgs, typedArg)
		}
	}

	return extractedArgs
}

func (b *Baranka) getPlaceholders(args []any) []string {
	var placeholders = make([]string, 0, len(args))
	for _, arg := range args {
		placeholders = append(placeholders, b.getPlaceholder(arg))
	}
	return placeholders
}

func (b *Baranka) getPlaceholder(arg any) string {
	switch typedArg := arg.(type) {
	case Expression:
		return fmt.Sprintf(typedArg.template, toAny(b.getPlaceholders(typedArg.args))...)
	}

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

func toAny[T any](slice []T) []any {
	result := make([]any, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

// Args returns the collected arguments in order.
func (b *Baranka) Args() []any {
	return b.args
}

// Values returns the SQL value blocks as a string, e.g. ($1,$2),\n($3,$4).
func (b *Baranka) Values() string {
	if len(b.blocks) == 0 {
		return ""
	}

	estimatedSize := 0
	for i := range b.blocks {
		estimatedSize += len(b.blocks[i])
	}
	estimatedSize += (len(b.blocks) - 1) * 2 // 2 for the commas and newlines

	var builder strings.Builder
	builder.Grow(estimatedSize)

	builder.WriteString(b.blocks[0])
	for i := 1; i < len(b.blocks); i++ {
		builder.WriteString(",\n")
		builder.WriteString(b.blocks[i])
	}

	return builder.String()
}
