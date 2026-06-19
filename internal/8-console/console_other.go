//go:build !windows

package console

import (
	"fmt"
	"os"
)

const (
	ColorRed    = ""
	ColorGreen  = ""
	ColorYellow = ""
	ColorGray   = ""
	ColorCyan   = ""
	ColorDim    = ""
	ColorReset  = ""
)

const (
	ProgramName        = "twaz"
	ProgramDescription = "Check and fix Tailwind CSS class order in JSX/TSX files."
)

// PrintVersion writes the program name, description, and version at startup.
func PrintVersion(version string) {
	fmt.Printf("%s — %s (version %s)\n\n", ProgramName, ProgramDescription, version)
}

// PrintError writes err to stderr and exits.
func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

// WaitForAnyKey is a no-op on non-Windows platforms.
func WaitForAnyKey() {}

// WaitAndExit exits immediately without waiting for a key.
func WaitAndExit(code int) {
	os.Exit(code)
}
