package twaz

import (
	"strings"
	"testing"
)

func TestParseArgsFixDefault(t *testing.T) {
	args := ParseArgs([]string{"src"})
	if !args.Fix {
		t.Fatal("expected fix enabled by default")
	}
}

func TestParseArgsCheckDisablesFix(t *testing.T) {
	args := ParseArgs([]string{"--check", "src"})
	if args.Fix {
		t.Fatal("expected fix disabled with --check")
	}
}

func TestTargetFolderLabel(t *testing.T) {
	if got := TargetFolderLabel([]string{"testdata"}); got != "testdata" {
		t.Fatalf("dir label = %q, want testdata", got)
	}
}

	func TestSortClassString(t *testing.T) {
	input := "bg-muted text-sm absolute top-0"
	want := "absolute top-0 text-sm bg-muted"
	if got := SortClassString(input); got != want {
		t.Fatalf("SortClassString() = %q, want %q", got, want)
	}
}

func TestCheckClassString(t *testing.T) {
	violations := CheckClassString("bg-muted text-sm absolute")
	if len(violations) < 2 {
		t.Fatalf("expected at least 2 violations, got %d", len(violations))
	}
	if violations[0].Token != "text-sm" {
		t.Fatalf("first token = %q, want text-sm", violations[0].Token)
	}
}

func TestClassifyTextSizeVsColor(t *testing.T) {
	if Classify("text-sm") != 7 {
		t.Fatalf("text-sm group = %d, want 7", Classify("text-sm"))
	}
	if Classify("text-red-500") != 9 {
		t.Fatalf("text-red-500 group = %d, want 9", Classify("text-red-500"))
	}
}

func TestClassifyGroupNamed(t *testing.T) {
	if Classify("group/accordion-trigger") != 0 {
		t.Fatalf("group/accordion-trigger group = %d, want 0", Classify("group/accordion-trigger"))
	}
}

func TestExtractClassStrings(t *testing.T) {
	content := `<div className="bg-muted text-sm absolute" />`
	matches := ExtractClassStrings(content)
	if len(matches) != 1 {
		t.Fatalf("matches = %d, want 1", len(matches))
	}
	if matches[0].Value != "bg-muted text-sm absolute" {
		t.Fatalf("value = %q", matches[0].Value)
	}
}

func TestExtractClassStringsCn(t *testing.T) {
	content := `cn("px-4 h-9 text-sm")`
	matches := ExtractClassStrings(content)
	if len(matches) != 1 || matches[0].Value != "px-4 h-9 text-sm" {
		t.Fatalf("unexpected matches: %+v", matches)
	}
}

func TestExtractTemplateLiteralWithInterpolation(t *testing.T) {
	content := "className={`h-8 ${isInvalid ? 'border-red-500 focus-visible:ring-red-500' : ''}`}"
	matches := ExtractClassStrings(content)

	var values []string
	for _, m := range matches {
		values = append(values, m.Value)
	}

	expectContains(t, values, "h-8")
	expectContains(t, values, "border-red-500 focus-visible:ring-red-500")
}

func TestExtractTemplateLiteralMultipleInterpolations(t *testing.T) {
	content := "className={`border-0 border-l ${isRegex ? 'text-primary bg-primary/10' : 'hover:bg-primary/5'} ${isInvalid ? 'border-l-red-500' : ''}`}"
	matches := ExtractClassStrings(content)

	var values []string
	for _, m := range matches {
		values = append(values, m.Value)
	}

	expectContains(t, values, "border-0 border-l")
	expectContains(t, values, "text-primary bg-primary/10")
	expectContains(t, values, "hover:bg-primary/5")
	expectContains(t, values, "border-l-red-500")
}

func TestExtractTemplateLiteralStaticBeforeAndAfter(t *testing.T) {
	content := "className={`pr-8 h-8 ${expr} mt-4 px-5`}"
	matches := ExtractClassStrings(content)

	var values []string
	for _, m := range matches {
		values = append(values, m.Value)
	}

	expectContains(t, values, "pr-8 h-8")
	expectContains(t, values, "mt-4 px-5")
}

func TestExtractTemplateLiteralNoInterpolation(t *testing.T) {
	content := "className={`h-8 pr-8 text-sm`}"
	matches := ExtractClassStrings(content)
	if len(matches) != 1 || matches[0].Value != "h-8 pr-8 text-sm" {
		t.Fatalf("expected single match 'h-8 pr-8 text-sm', got %+v", matches)
	}
}

func TestExtractMultiArgCn(t *testing.T) {
	content := `cn("px-4 h-9", "text-sm font-medium")`
	matches := ExtractClassStrings(content)

	var values []string
	for _, m := range matches {
		values = append(values, m.Value)
	}

	expectContains(t, values, "px-4 h-9")
	expectContains(t, values, "text-sm font-medium")
}

func TestExtractMultiArgClassNames(t *testing.T) {
	content := `classNames("mx-5 mt-1 text-xs text-muted-foreground", className)`
	matches := ExtractClassStrings(content)

	var values []string
	for _, m := range matches {
		values = append(values, m.Value)
	}

	expectContains(t, values, "mx-5 mt-1 text-xs text-muted-foreground")
}

func TestExtractMultiArgCnThreeStrings(t *testing.T) {
	content := `cn("bg-muted", "text-sm absolute", "border rounded-md")`
	matches := ExtractClassStrings(content)

	var values []string
	for _, m := range matches {
		values = append(values, m.Value)
	}

	expectContains(t, values, "bg-muted")
	expectContains(t, values, "text-sm absolute")
	expectContains(t, values, "border rounded-md")
}

func TestExtractInterpolationDoesNotProduceFalseTokens(t *testing.T) {
	content := "className={`h-8 ${flag ? 'border-red-500' : ''}`}"
	matches := ExtractClassStrings(content)

	for _, m := range matches {
		if strings.Contains(m.Value, "${") || strings.Contains(m.Value, "?") {
			t.Fatalf("value contains interpolation syntax: %q", m.Value)
		}
	}
}

func expectContains(t *testing.T, values []string, want string) {
	t.Helper()
	for _, v := range values {
		if v == want {
			return
		}
	}
	t.Fatalf("expected values to contain %q, got %v", want, values)
}
