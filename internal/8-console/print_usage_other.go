//go:build !windows

package console

import (
	"fmt"
	"os"
)

// PrintUsage shows help text and exits.
func PrintUsage(help UsageHelp) {
	fmt.Println("Usage:")
	fmt.Printf("  %s\n\n", help.Syntax)

	if len(help.Options) > 0 {
		fmt.Println("Options:")
		for _, opt := range help.Options {
			fmt.Printf("  %s  %s\n", opt.Flag, opt.Description)
		}
		fmt.Println()
	}

	if len(help.Args) > 0 {
		fmt.Println("Arguments:")
		for _, arg := range help.Args {
			fmt.Printf("  %s: %s\n", arg.Label, arg.Value)
		}
		fmt.Println()
	}

	if len(help.Examples) > 0 {
		fmt.Println("Examples:")
		for _, example := range help.Examples {
			fmt.Printf("  %s\n", example)
		}
		fmt.Println()
	}

	if help.Message != "" {
		fmt.Println(help.Message)
		fmt.Println()
	}

	os.Exit(0)
}
