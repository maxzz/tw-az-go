package twaz

import "regexp"

const (
	variantGroup        = 11
	transitionGroup     = 12
	borderGroup         = 13
	roundingGroup       = 14
	shadowGroup         = 15
	truncateOverflowGrp = 16
	childrenGroup       = 17
	endGroup            = 18
	unknownGroup        = 999
)

var defaultExtensions = []string{".tsx", ".jsx"}

var groupNames = []string{
	"position anchor",
	"position offsets",
	"self & group",
	"element",
	"margin & padding",
	"width & height",
	"display",
	"text size",
	"font",
	"text color",
	"background & fill color",
	"variant modifiers",
	"transition",
	"border",
	"rounding",
	"shadow",
	"truncate & overflow",
	"children",
	"end",
}

var textSizes = map[string]struct{}{
	"text-xs":  {},
	"text-sm":  {},
	"text-base": {},
	"text-lg":  {},
	"text-xl":  {},
	"text-2xl": {},
	"text-3xl": {},
	"text-4xl": {},
	"text-5xl": {},
	"text-6xl": {},
	"text-7xl": {},
	"text-8xl": {},
	"text-9xl": {},
}

var classPatterns = []*regexp.Regexp{
	regexp.MustCompile(`className\s*=\s*"([^"]+)"`),
	regexp.MustCompile(`className\s*=\s*'([^']+)'`),
	regexp.MustCompile(`className\s*=\s*\{` + "`" + `([^` + "`" + `]+)` + "`" + `\}`),
	regexp.MustCompile(`className\s*=\s*\{\s*["'` + "`" + `]([^"'` + "`" + `]+)["'` + "`" + `]\s*\}`),
	regexp.MustCompile(`class\s*=\s*"([^"]+)"`),
	regexp.MustCompile(`cn\(\s*["'` + "`" + `]([^"'` + "`" + `]+)["'` + "`" + `]`),
	regexp.MustCompile(`classNames\(\s*["'` + "`" + `]([^"'` + "`" + `]+)["'` + "`" + `]`),
}

var reFuncCall = regexp.MustCompile(`(?:cn|classNames)\(`)

var ignoredDirectories = map[string]struct{}{
	"node_modules": {},
	"dist":         {},
	".git":         {},
}
