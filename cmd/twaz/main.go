package main

import (
	"fmt"
	"os"
	"path/filepath"

	"twaz/twaz"
)

func main() {
	args := twaz.ParseArgs(os.Args[1:])

	if args.Help {
		fmt.Print(twaz.Help())
		os.Exit(0)
	}

	paths := args.Paths
	if len(paths) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "twaz: %v\n", err)
			os.Exit(1)
		}
		paths = []string{cwd}
	} else {
		resolved := make([]string, len(paths))
		for i, p := range paths {
			abs, err := filepath.Abs(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "twaz: %v\n", err)
				os.Exit(1)
			}
			resolved[i] = abs
		}
		paths = resolved
	}

	result := twaz.RunScan(paths, twaz.ScanOptions{Fix: args.Fix}, os.Stdout)

	if args.Fix {
		if len(result.Violations) > 0 {
			twaz.PrintViolations(os.Stdout, result.Violations, result.FileCount)
			os.Exit(1)
		}
		os.Exit(0)
	}

	twaz.PrintViolations(os.Stdout, result.Violations, result.FileCount)

	if len(result.Violations) > 0 {
		fmt.Println("Run with --fix to reorder classes automatically.")
		fmt.Println()
		os.Exit(1)
	}

	os.Exit(0)
}
