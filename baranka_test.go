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

func BenchmarkExtractArgs_Original(b *testing.B) {
	// Create test data with 200 blocks worth of args
	args := make([]any, 0, 200)
	for i := 0; i < 200; i++ {
		if i%3 == 0 {
			args = append(args, Expression{
				template: "ST_Point(%s, %s)",
				args:     []any{float64(i), float64(i + 1)},
			})
		} else {
			args = append(args, i, "value"+string(rune(i)))
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
	for i := 0; i < 200; i++ {
		if i%3 == 0 {
			args = append(args, Expression{
				template: "ST_Point(%s, %s)",
				args:     []any{float64(i), float64(i + 1)},
			})
		} else {
			args = append(args, i, "value"+string(rune(i)))
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
	for i := 0; i < 200; i++ {
		baranka.Add(i, "value"+string(rune(i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := strings.Join(baranka.blocks, ",\n")
		_ = result
	}
}

func BenchmarkValues_StringsBuilder(b *testing.B) {
	b.ReportAllocs()
	// Create Baranka with 200 blocks
	baranka := New()
	for i := 0; i < 200; i++ {
		baranka.Add(i, "value"+string(rune(i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		baranka.Values()
	}
}
