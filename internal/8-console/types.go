package console

// UsageOption describes one CLI flag.
type UsageOption struct {
	Flag        string
	Description string
}

// UsageArg labels one token from the command line.
type UsageArg struct {
	Label string
	Value string
}

// UsageHelp is the content shown for incorrect user input.
type UsageHelp struct {
	Message  string
	Syntax   string
	Options  []UsageOption
	Args     []UsageArg
	Examples []string
}

const (
	OperationScan = "Check class order"
	OperationFix  = "Fix class order"
)
