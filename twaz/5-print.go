package twaz

import (
	"fmt"
	"io"
	"sort"
)

// PrintViolations writes a human-readable violation report.
func PrintViolations(out io.Writer, violations []FileViolation, fileCount int) {
	byFile := make(map[string][]FileViolation)
	for _, violation := range violations {
		byFile[violation.File] = append(byFile[violation.File], violation)
	}

	files := make([]string, 0, len(byFile))
	for file := range byFile {
		files = append(files, file)
	}
	sort.Strings(files)

	fmt.Fprintf(out, "Scanned %d files\n", fileCount)
	fmt.Fprintf(out, "Found %d class strings with order violations in %d files\n\n", len(violations), len(byFile))

	for _, file := range files {
		items := byFile[file]
		fmt.Fprintln(out, file)
		for _, item := range items {
			preview := item.Value
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			fmt.Fprintf(out, "  L%d: %s\n", item.Line, preview)
			for _, v := range item.Violations {
				fmt.Fprintf(out, "    - \"%s\" (%s) appears after %s\n", v.Token, v.Group, v.After)
			}
		}
		fmt.Fprintln(out)
	}
}
