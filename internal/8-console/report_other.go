//go:build !windows

package console

import (
	"fmt"
	"sort"

	"twaz/twaz"
)

// PrintScanReport writes a plain scan/fix summary and violation details.
func PrintScanReport(result twaz.ScanResult, operation, targetFolder string, showFixHint bool) {
	if operation == OperationFix && targetFolder != "" {
		fmt.Println(targetFolder)
		fmt.Println()
	}

	fmt.Println(operation)
	fmt.Println()

	if operation == OperationFix {
		fmt.Printf("Scanned %d files\n", result.FileCount)
		fmt.Printf("Fixed %d class string%s\n", result.FixedCount, pluralSuffix(result.FixedCount))
		if result.FixedCount > 0 {
			fmt.Println()
			fmt.Println("Re-checking after fix...")
			fmt.Println()
		}
	}

	printViolations(result.Violations, result.FileCount, operation != OperationFix)

	if showFixHint && len(result.Violations) > 0 {
		fmt.Println("Run without --check to reorder classes automatically.")
	}
}

func printViolations(violations []twaz.FileViolation, fileCount int, includeScanSummary bool) {
	byFile := make(map[string][]twaz.FileViolation)
	for _, violation := range violations {
		byFile[violation.File] = append(byFile[violation.File], violation)
	}

	files := make([]string, 0, len(byFile))
	for file := range byFile {
		files = append(files, file)
	}
	sort.Strings(files)

	if includeScanSummary {
		fmt.Printf("Scanned %d files\n", fileCount)
		fmt.Printf("Found %d class strings with order violations in %d files\n\n", len(violations), len(byFile))
	} else if len(violations) > 0 {
		fmt.Printf("Found %d class strings with order violations in %d files\n\n", len(violations), len(byFile))
	}

	for _, file := range files {
		fmt.Println(file)
		items := byFile[file]
		for _, item := range items {
			preview := item.Value
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			fmt.Printf("  L%d: %s\n", item.Line, preview)
			for _, v := range item.Violations {
				fmt.Printf("    - \"%s\" (%s) appears after %s\n", v.Token, v.Group, v.After)
			}
		}
		fmt.Println()
	}

	if len(violations) == 0 && includeScanSummary {
		fmt.Println("No class order violations found.")
	}
}

func pluralSuffix(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
