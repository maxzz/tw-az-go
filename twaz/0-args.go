package twaz

import (
	"fmt"
	"os"
	"path/filepath"
)

// ParsedArgs holds CLI arguments.
type ParsedArgs struct {
	Fix   bool
	Help  bool
	Paths []string
}

// Help returns the CLI help text.
func Help() string {
	return fmt.Sprintf(`twaz — check and fix Tailwind CSS class order (version %s)

Usage:
  twaz [options] [paths...]

Options:
  --check, -c  Report violations only; do not reorder classes
  --fix, -f    Reorder classes automatically in place (default)
  --help, -h   Show this help message

Arguments:
  paths       Files or directories to scan (default: current directory)

Examples:
  twaz src
  twaz --check src
  twaz src/App.tsx
`, Version)
}

// ParseArgs parses CLI arguments. Fix mode is enabled by default.
func ParseArgs(argv []string) ParsedArgs {
	var paths []string
	fix := true
	var help bool

	for _, arg := range argv {
		switch arg {
		case "--check", "-c", "--no-fix":
			fix = false
		case "--fix", "-f":
			fix = true
		case "--help", "-h":
			help = true
		default:
			if len(arg) == 0 || arg[0] == '-' {
				continue
			}
			paths = append(paths, arg)
		}
	}

	return ParsedArgs{Fix: fix, Help: help, Paths: paths}
}

// TargetFolderLabel returns the display name of the primary scan folder.
func TargetFolderLabel(paths []string) string {
	if len(paths) == 0 {
		return "."
	}

	p := paths[0]
	info, err := os.Stat(p)
	if err != nil {
		return filepath.Base(p)
	}
	if info.IsDir() {
		return filepath.Base(p)
	}
	return filepath.Base(filepath.Dir(p))
}
