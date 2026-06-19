//go:build windows

package console

import (
	"bufio"
	"fmt"
	"os"
	"unsafe"

	ascii "twaz/internal/8-result-ascii"

	"golang.org/x/sys/windows"
)

const (
	enableVirtualTerminalProcessing = 0x0004
)

var (
	kernel32           = windows.NewLazySystemDLL("kernel32.dll")
	procGetStdHandle   = kernel32.NewProc("GetStdHandle")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
	procGetch          = windows.NewLazySystemDLL("msvcrt.dll").NewProc("_getch")
)

const stdOutputHandle = ^uintptr(10) // -11 as uintptr
const stdInputHandle = ^uintptr(9)   // -10 as uintptr

const (
	ColorRed    = "\x1b[31m"
	ColorGreen  = "\x1b[32m"
	ColorYellow = "\x1b[33m"
	ColorGray   = "\x1b[90m"
	ColorCyan   = "\x1b[36m"
	ColorDim    = "\x1b[2m\x1b[90m"
	ColorReset  = "\x1b[0m"
)

const (
	ProgramName        = "twaz"
	ProgramDescription = "Check and fix Tailwind CSS class order in JSX/TSX files."
)

func init() {
	enableColors()
}

func enableColors() {
	handle, _, _ := procGetStdHandle.Call(stdOutputHandle)
	if handle == 0 {
		return
	}

	var mode uint32
	r1, _, _ := procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	if r1 == 0 {
		return
	}

	procSetConsoleMode.Call(handle, uintptr(mode|enableVirtualTerminalProcessing))
}

// PrintVersion writes the program name, description, and version at startup.
func PrintVersion(version string) {
	fmt.Printf("%s — %s", ProgramName, ProgramDescription)
	fmt.Printf(" %s(version %s)%s\n\n", ColorGray, version, ColorReset)
}

// PrintError writes Inspector Gadget art and err in red, then waits for a key press before exiting.
func PrintError(err error) {
	fmt.Print(ColorRed)
	fmt.Print(ascii.InspectorGadget())
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	fmt.Print(ColorReset)
	fmt.Print("Press any key to close...")
	waitForKey()
	os.Exit(1)
}

// WaitForAnyKey blocks until the user presses any key.
func WaitForAnyKey() {
	waitForKey()
}

// WaitAndExit waits for any key, then exits with the given code.
func WaitAndExit(code int) {
	fmt.Fprint(os.Stdout, "\nPress any key to close the window")
	waitForKey()
	os.Exit(code)
}

func waitForKey() {
	handle, _, _ := procGetStdHandle.Call(stdInputHandle)
	if handle == 0 {
		waitForKeyFallback()
		return
	}

	var mode uint32
	r1, _, _ := procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	if r1 == 0 {
		waitForKeyFallback()
		return
	}

	procGetch.Call()
}

func waitForKeyFallback() {
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadByte()
}
