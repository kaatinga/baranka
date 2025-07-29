package baranka

import (
	"strings"
	"testing"
)

func TestBaranka_Add_DollarPlaceholders(t *testing.T) {
	b := getOptions([]option{
		WithIncludeTemplate("(%s)"),
		WithPlaceholderFormat(PlaceholderFormatDollar),
	})

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
	b := getOptions([]option{
		WithIncludeTemplate("(%s)"),
		WithPlaceholderFormat(PlaceholderFormatQuestionMark),
	})

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
	b := NewBaranka()
	if len(b.Args()) != 0 {
		t.Errorf("expected no args, got %d", len(b.Args()))
	}
	if b.Values() != "" {
		t.Errorf("expected empty values, got %q", b.Values())
	}
}
