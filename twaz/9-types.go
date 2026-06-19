package twaz

// ScanOptions configures a directory scan.
type ScanOptions struct {
	Extensions []string
	Fix        bool
}

// ClassMatch is a class string occurrence in source text.
type ClassMatch struct {
	Value  string
	Index  int
	Length int
	Full   string
}

// OrderViolation describes a single out-of-order utility class.
type OrderViolation struct {
	Token string
	Group string
	After string
}

// FileViolation groups violations for one class string in a file.
type FileViolation struct {
	File       string
	Line       int
	Value      string
	Violations []OrderViolation
}

// ScanResult is returned by RunScan.
type ScanResult struct {
	FileCount  int
	Violations []FileViolation
	FixedCount int
}

type fixReplacement struct {
	index        int
	length       int
	replacement  string
}
