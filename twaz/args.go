package twaz

import "fmt"

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
  --fix, -f   Reorder classes automatically in place
  --help, -h  Show this help message

Arguments:
  paths       Files or directories to scan (default: current directory)

Examples:
  twaz src
  twaz --fix src/components
  twaz src/App.tsx
`, Version)
}

// ParseArgs parses CLI arguments.
func ParseArgs(argv []string) ParsedArgs {
	var paths []string
	var fix, help bool

	for _, arg := range argv {
		switch arg {
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
