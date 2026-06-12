package baranka

import (
	"fmt"
	"strings"
)

// Expression is a SQL fragment with its own placeholders, e.g. "POINT(%s %s)".
// Each %s in the template is replaced with a placeholder ($1, ?, ...) and the
// corresponding argument is collected for the driver. Construct it with
// NewExpression. Literal % characters are not supported in the template.
type Expression struct {
	template string
	args     []any
}

// NewExpression creates an Expression. It panics if the number of %s verbs in
// the template does not match the number of arguments, or if the template
// contains any other % verb — both would silently produce invalid SQL.
func NewExpression(template string, args ...any) Expression {
	placeholders := strings.Count(template, "%s")
	if placeholders != len(args) {
		panic(fmt.Sprintf("baranka: expression template %q has %d %%s placeholders but %d arguments", template, placeholders, len(args)))
	}
	if strings.Count(template, "%") != placeholders {
		panic(fmt.Sprintf("baranka: expression template %q must use only %%s placeholders", template))
	}

	return Expression{
		template: template,
		args:     args,
	}
}
