package twaz

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunScan scans paths for Tailwind class order violations.
func RunScan(paths []string, options ScanOptions, out io.Writer) ScanResult {
	extensions := options.Extensions
	if len(extensions) == 0 {
		extensions = defaultExtensions
	}

	scanPaths := paths
	if len(scanPaths) == 0 {
		scanPaths = []string{"."}
	}

	rootDir, _ := filepath.Abs(scanPaths[0])
	files := collectFiles(scanPaths, extensions)

	if options.Fix {
		fixedCount := 0
		for _, file := range files {
			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			fixedCount += applyFixes(file, string(content), ExtractClassStrings(string(content)))
		}

		if out != nil {
			fmt.Fprintf(out, "Scanned %d files\n", len(files))
			fmt.Fprintf(out, "Fixed %d class string%s\n\n", fixedCount, plural(fixedCount))
		}

		if fixedCount > 0 {
			if out != nil {
				fmt.Fprintln(out, "Re-checking after fix...")
			fmt.Fprintln(out)
			}
			remaining := scanForViolations(files, rootDir)
			return ScanResult{FileCount: len(files), Violations: remaining, FixedCount: fixedCount}
		}

		return ScanResult{FileCount: len(files), Violations: nil, FixedCount: fixedCount}
	}

	violations := scanForViolations(files, rootDir)
	return ScanResult{FileCount: len(files), Violations: violations, FixedCount: 0}
}

func scanForViolations(filePaths []string, rootDir string) []FileViolation {
	var violations []FileViolation

	for _, file := range filePaths {
		contentBytes, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		content := string(contentBytes)
		rel, err := filepath.Rel(rootDir, file)
		if err != nil {
			rel = file
		}
		rel = filepath.ToSlash(rel)

		for _, match := range ExtractClassStrings(content) {
			classViolations := CheckClassString(match.Value)
			if len(classViolations) == 0 {
				continue
			}
			line := strings.Count(content[:match.Index], "\n") + 1
			violations = append(violations, FileViolation{
				File:       rel,
				Line:       line,
				Value:      match.Value,
				Violations: classViolations,
			})
		}
	}

	return violations
}

func applyFixes(file string, content string, matches []ClassMatch) int {
	var toApply []fixReplacement

	for _, match := range matches {
		if len(CheckClassString(match.Value)) == 0 {
			continue
		}
		fixed := SortClassString(match.Value)
		if fixed == match.Value {
			continue
		}
		toApply = append(toApply, fixReplacement{
			index:       match.Index,
			length:      match.Length,
			replacement: strings.Replace(match.Full, match.Value, fixed, 1),
		})
	}

	if len(toApply) == 0 {
		return 0
	}

	sort.Slice(toApply, func(i, j int) bool {
		return toApply[i].index > toApply[j].index
	})

	updated := content
	for _, item := range toApply {
		updated = updated[:item.index] + item.replacement + updated[item.index+item.length:]
	}

	if err := os.WriteFile(file, []byte(updated), 0o644); err != nil {
		return 0
	}

	return len(toApply)
}

func collectFiles(paths []string, extensions []string) []string {
	seen := make(map[string]struct{})
	var files []string

	for _, inputPath := range paths {
		resolved, err := filepath.Abs(inputPath)
		if err != nil {
			continue
		}
		info, err := os.Stat(resolved)
		if err != nil {
			continue
		}

		if info.IsDir() {
			walk(resolved, extensions, seen, &files)
		} else if hasExtension(resolved, extensions) {
			if _, ok := seen[resolved]; !ok {
				seen[resolved] = struct{}{}
				files = append(files, resolved)
			}
		}
	}

	sort.Strings(files)
	return files
}

func walk(dir string, extensions []string, seen map[string]struct{}, files *[]string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		full := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if _, ignored := ignoredDirectories[entry.Name()]; ignored {
				continue
			}
			walk(full, extensions, seen, files)
			continue
		}
		if hasExtension(full, extensions) {
			if _, ok := seen[full]; !ok {
				seen[full] = struct{}{}
				*files = append(*files, full)
			}
		}
	}
}

func hasExtension(path string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}

func plural(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
