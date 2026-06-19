//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"

	console "twaz/internal/8-console"
	"twaz/twaz"
)

func main() {
	args := twaz.ParseArgs(os.Args[1:])

	if args.Help {
		printUsage()
	}

	console.PrintVersion(twaz.Version)

	paths, err := resolvePaths(args.Paths)
	if err != nil {
		console.PrintError(err)
	}

	operation := console.OperationScan
	if args.Fix {
		operation = console.OperationFix
	}

	result := twaz.RunScan(paths, twaz.ScanOptions{Fix: args.Fix})
	targetFolder := ""
	if args.Fix {
		targetFolder = twaz.TargetFolderLabel(paths)
	}
	console.PrintScanReport(result, operation, targetFolder, !args.Fix)

	exitCode := 0
	if len(result.Violations) > 0 {
		exitCode = 1
	}
	console.WaitAndExit(exitCode)
}

func printUsage() {
	console.PrintUsage(console.UsageHelp{
		Message: "Check and fix Tailwind CSS utility class order in JSX/TSX files.",
		Syntax:  "twaz [options] [paths...]",
		Options: []console.UsageOption{
			{
				Flag:        "--check, -c",
				Description: "Report violations only; do not reorder classes (default: off; fix is enabled)",
			},
			{
				Flag:        "--fix, -f",
				Description: "Reorder classes automatically in place (default: on)",
			},
			{
				Flag:        "--help, -h",
				Description: "Show this help message (default: off)",
			},
		},
		Args: []console.UsageArg{
			{
				Label: "paths",
				Value: "Files or directories to scan (default: current directory)",
			},
		},
		Examples: []string{
			"twaz src",
			"twaz --check src",
			"twaz src/App.tsx",
		},
	})
}

func resolvePaths(paths []string) ([]string, error) {
	if len(paths) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("resolve working directory: %w", err)
		}
		return []string{cwd}, nil
	}

	resolved := make([]string, len(paths))
	for i, p := range paths {
		abs, err := filepath.Abs(p)
		if err != nil {
			return nil, fmt.Errorf("resolve path %q: %w", p, err)
		}
		resolved[i] = abs
	}
	return resolved, nil
}
