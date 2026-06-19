//go:build windows

package console

import (
	"fmt"
	"os"
	"sort"

	"twaz/twaz"
)

// PrintScanReport writes a colored scan/fix summary and violation details.
func PrintScanReport(result twaz.ScanResult, operation, targetFolder string, showFixHint bool) {
	if operation == OperationFix && targetFolder != "" {
		fmt.Fprintln(os.Stdout, targetFolder)
		fmt.Fprintln(os.Stdout)
	}

	fmt.Fprintf(os.Stdout, "%s%s%s\n", ColorCyan, operation, ColorReset)
	fmt.Fprintln(os.Stdout)

	if operation == OperationFix {
		fmt.Fprintf(os.Stdout, "Scanned %s%d%s files\n", ColorGray, result.FileCount, ColorReset)
		fmt.Fprintf(
			os.Stdout,
			"Fixed %s%d%s class string%s\n",
			ColorGray,
			result.FixedCount,
			ColorReset,
			pluralSuffix(result.FixedCount),
		)
		if result.FixedCount > 0 {
			fmt.Fprintln(os.Stdout)
			fmt.Fprintf(os.Stdout, "%sRe-checking after fix...%s\n", ColorDim, ColorReset)
			fmt.Fprintln(os.Stdout)
		}
	}

	printViolations(result.Violations, result.FileCount, operation != OperationFix)

	if showFixHint && len(result.Violations) > 0 {
		fmt.Fprintf(os.Stdout, "%sRun without --check to reorder classes automatically.%s\n", ColorYellow, ColorReset)
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
		fmt.Fprintf(os.Stdout, "Scanned %s%d%s files\n", ColorGray, fileCount, ColorReset)
		fmt.Fprintf(
			os.Stdout,
			"Found %s%d%s class strings with order violations in %s%d%s files\n\n",
			ColorGray,
			len(violations),
			ColorReset,
			ColorGray,
			len(byFile),
			ColorReset,
		)
	} else if len(violations) > 0 {
		fmt.Fprintf(
			os.Stdout,
			"Found %s%d%s class strings with order violations in %s%d%s files\n\n",
			ColorGray,
			len(violations),
			ColorReset,
			ColorGray,
			len(byFile),
			ColorReset,
		)
	}

	for _, file := range files {
		fmt.Fprintln(os.Stdout, file)
		items := byFile[file]
		for _, item := range items {
			preview := item.Value
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			fmt.Fprintf(os.Stdout, "  %sL%d:%s %s\n", ColorGray, item.Line, ColorReset, preview)
			for _, v := range item.Violations {
				fmt.Fprintf(
					os.Stdout,
					"    - %s\"%s\"%s %s(%s) appears after %s%s\n",
					ColorYellow,
					v.Token,
					ColorReset,
					ColorDim,
					v.Group,
					v.After,
					ColorReset,
				)
			}
		}
		fmt.Fprintln(os.Stdout)
	}

	if len(violations) == 0 && operationSuccessMessage(includeScanSummary) {
		fmt.Fprintf(os.Stdout, "%sNo class order violations found.%s\n", ColorGreen, ColorReset)
	}
}

func operationSuccessMessage(includeScanSummary bool) bool {
	return includeScanSummary
}

func pluralSuffix(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
