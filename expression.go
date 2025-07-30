package baranka

type Expression struct {
	template string
	args     []any
}

func NewExpression(template string, args ...any) Expression {
	return Expression{
		template: template,
		args:     args,
	}
}
