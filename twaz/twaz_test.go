package twaz

import "testing"

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
	if Classify("group/accordion-trigger") != 2 {
		t.Fatalf("group/accordion-trigger group = %d, want 2", Classify("group/accordion-trigger"))
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
