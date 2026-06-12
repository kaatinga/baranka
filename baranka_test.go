package baranka

import (
	"fmt"
	"strings"
	"testing"
)

func TestBaranka_Add_DollarPlaceholders(t *testing.T) {
	b := New(
		WithIncludeTemplate("(%s)"),
		WithPlaceholderFormat(PlaceholderFormatDollar),
	)

	b.Add(1, "foo")
	b.Add(2, "bar")

	expectedArgs := []any{1, "foo", 2, "bar"}
	if len(b.Args()) != 4 {
		t.Fatalf("expected 4 args, got %d", len(b.Args()))
	}
	for i, v := range expectedArgs {
		if b.Args()[i] != v {
			t.Errorf("arg %d: expected %v, got %v", i, v, b.Args()[i])
		}
	}

	values := b.Values()
	if !strings.Contains(values, "($1,$2)") || !strings.Contains(values, "($3,$4)") {
		t.Errorf("unexpected values: %s", values)
	}
}

func TestBaranka_Add_QuestionMarkPlaceholders(t *testing.T) {
	b := New(
		WithIncludeTemplate("(%s)"),
		WithPlaceholderFormat(PlaceholderFormatQuestionMark),
	)

	b.Add(1, 2)
	b.Add(3, 4)

	expectedArgs := []any{1, 2, 3, 4}
	if len(b.Args()) != 4 {
		t.Fatalf("expected 4 args, got %d", len(b.Args()))
	}
	for i, v := range expectedArgs {
		if b.Args()[i] != v {
			t.Errorf("arg %d: expected %v, got %v", i, v, b.Args()[i])
		}
	}

	values := b.Values()
	if !strings.Contains(values, "(?,?)") {
		t.Errorf("unexpected values: %s", values)
	}
}

func TestBaranka_Empty(t *testing.T) {
	b := New()
	if len(b.Args()) != 0 {
		t.Errorf("expected no args, got %d", len(b.Args()))
	}
	if b.Values() != "" {
		t.Errorf("expected empty values, got %q", b.Values())
	}
}

func TestBaranka_Add_WithExpression(t *testing.T) {
	b := New(WithBlocksLength(2), WithPlaceholderFormat(PlaceholderFormatDollar))

	type item struct {
		X float32
		Y float32
	}
	points := []item{
		{10.1, 20.2},
		{11.1, 21.2},
	}

	for _, it := range points {
		b.Add(
			Expression{
				template: "ST_SetSRID(ST_Point(%s, %s)",
				args:     []any{it.X, it.Y},
			},
		)
	}

	values := b.Values()
	if !strings.Contains(values, "ST_SetSRID(ST_Point($1, $2)") || !strings.Contains(values, "ST_SetSRID(ST_Point($3, $4)") {
		t.Errorf("unexpected values: %s", values)
	}

	expectedArgs := []any{
		float32(10.1), float32(20.2),
		float32(11.1), float32(21.2),
	}
	args := b.Args()
	if len(args) != len(expectedArgs) {
		t.Fatalf("expected %d args, got %d", len(expectedArgs), len(args))
	}
	for i, v := range expectedArgs {
		if args[i] != v {
			t.Errorf("arg %d: expected %v, got %v", i, v, args[i])
		}
	}

	// Example of how the query could look
	queryTemplate := `
		INSERT INTO test_points (geom)
		VALUES %s
	`
	_ = fmt.Sprintf(queryTemplate, b.Values())
}

func TestBaranka_Reset(t *testing.T) {
	b := New()
	b.Add(1, "foo")
	b.Add(2, "bar")

	b.Reset()

	if len(b.Args()) != 0 {
		t.Errorf("expected no args after reset, got %d", len(b.Args()))
	}
	if b.Values() != "" {
		t.Errorf("expected empty values after reset, got %q", b.Values())
	}

	b.Add(3, "baz")
	if values := b.Values(); values != "($1,$2)" {
		t.Errorf("expected placeholders to restart at $1, got %q", values)
	}
}

func TestBaranka_Add_NestedExpression(t *testing.T) {
	b := New()

	b.Add(NewExpression("ST_SetSRID(%s, %s)", NewExpression("ST_Point(%s, %s)", 10.1, 20.2), 4326))

	values := b.Values()
	if values != "(ST_SetSRID(ST_Point($1, $2), $3))" {
		t.Errorf("unexpected values: %s", values)
	}

	expectedArgs := []any{10.1, 20.2, 4326}
	args := b.Args()
	if len(args) != len(expectedArgs) {
		t.Fatalf("expected %d args, got %d", len(expectedArgs), len(args))
	}
	for i, v := range expectedArgs {
		if args[i] != v {
			t.Errorf("arg %d: expected %v, got %v", i, v, args[i])
		}
	}
}

func TestBaranka_Add_SingleExpression(t *testing.T) {
	b := New()
	b.Add(NewExpression("lower(%s)", "FOO"))

	if values := b.Values(); values != "(lower($1))" {
		t.Errorf("unexpected values: %s", values)
	}
	args := b.Args()
	if len(args) != 1 || args[0] != "FOO" {
		t.Errorf("expected [FOO], got %v", args)
	}
}

func TestNewExpression_Panics(t *testing.T) {
	expectPanic := func(name string, f func()) {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Error("expected panic")
				}
			}()
			f()
		})
	}

	expectPanic("too few args", func() { NewExpression("POINT(%s %s)", 1.0) })
	expectPanic("too many args", func() { NewExpression("POINT(%s)", 1.0, 2.0) })
	expectPanic("other verb", func() { NewExpression("POINT(%d)", 1) })
}

func TestWithIncludeTemplate_Invalid(t *testing.T) {
	for _, template := range []string{"()", "(%s, %s)", "(%d)"} {
		b := New(WithIncludeTemplate(template))
		b.Add(1)
		if values := b.Values(); values != "($1)" {
			t.Errorf("template %q: expected default template to be kept, got %q", template, values)
		}
	}

	b := New(WithIncludeTemplate("ROW(%s)"))
	b.Add(1)
	if values := b.Values(); values != "ROW($1)" {
		t.Errorf("expected custom template to apply, got %q", values)
	}
}

func BenchmarkExtractArgs_Original(b *testing.B) {
	// Create test data with 200 blocks worth of args
	args := make([]any, 0, 200)
	for i := range 200 {
		if i%3 == 0 {
			args = append(args, Expression{
				template: "ST_Point(%s, %s)",
				args:     []any{float64(i), float64(i + 1)},
			})
		} else {
			args = append(args, i, "value"+string(rune(i)))
		}
	}

	for b.Loop() {
		// Original single-pass version
		extractedArgs := make([]any, 0, len(args))
		for _, arg := range args {
			switch typedArg := arg.(type) {
			case Expression:
				extractedArgs = append(extractedArgs, typedArg.args...)
			default:
				extractedArgs = append(extractedArgs, typedArg)
			}
		}
		_ = extractedArgs
	}
}

func BenchmarkExtractArgs_TwoPass(b *testing.B) {
	args := make([]any, 0, 200)
	for i := range 200 {
		if i%3 == 0 {
			args = append(args, Expression{
				template: "ST_Point(%s, %s)",
				args:     []any{float64(i), float64(i + 1)},
			})
		} else {
			args = append(args, i, "value"+string(rune(i)))
		}
	}

	for b.Loop() {
		// Two-pass version with capacity estimation
		capacity := 0
		for _, arg := range args {
			switch typedArg := arg.(type) {
			case Expression:
				capacity += len(typedArg.args)
			default:
				capacity++
			}
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
		_ = extractedArgs
	}
}

func BenchmarkValues_StringsJoin(b *testing.B) {
	b.ReportAllocs()
	// Create Baranka with 200 blocks
	baranka := New()
	for i := range 200 {
		baranka.Add(i, "value"+string(rune(i)))
	}

	for b.Loop() {
		result := strings.Join(baranka.blocks, ",\n")
		_ = result
	}
}

func BenchmarkValues_StringsBuilder(b *testing.B) {
	b.ReportAllocs()
	// Create Baranka with 200 blocks
	baranka := New()
	for i := range 200 {
		baranka.Add(i, "value"+string(rune(i)))
	}

	for b.Loop() {
		baranka.Values()
	}
}
