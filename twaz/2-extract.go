package twaz

import (
	"fmt"
	"strings"
)

// ExtractClassStrings finds Tailwind class strings in JSX/TSX source.
func ExtractClassStrings(content string) []ClassMatch {
	var results []ClassMatch
	seen := make(map[string]struct{})

	for _, pattern := range classPatterns {
		locs := pattern.FindAllStringSubmatchIndex(content, -1)
		for _, loc := range locs {
			if len(loc) < 4 {
				continue
			}
			capturedStart, capturedEnd := loc[2], loc[3]
			if capturedStart < 0 || capturedEnd < 0 {
				continue
			}

			value := content[capturedStart:capturedEnd]

			if strings.Contains(value, "${") {
				subMatches := parseTemplateLiteral(content, capturedStart, capturedEnd)
				for _, sm := range subMatches {
					key := fmt.Sprintf("%d:%d", sm.Index, sm.Length)
					if _, ok := seen[key]; !ok {
						seen[key] = struct{}{}
						results = append(results, sm)
					}
				}
				continue
			}

			key := fmt.Sprintf("%d:%d", loc[0], loc[1]-loc[0])
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}

			results = append(results, ClassMatch{
				Value:  value,
				Index:  loc[0],
				Length: loc[1] - loc[0],
				Full:   content[loc[0]:loc[1]],
			})
		}
	}

	extractMultiArgCalls(content, seen, &results)

	return results
}

// parseTemplateLiteral extracts class strings from a template literal containing
// ${...} interpolation. Returns ClassMatch entries for each static segment and
// for string literals found inside interpolation expressions.
func parseTemplateLiteral(content string, capturedStart, capturedEnd int) []ClassMatch {
	var results []ClassMatch
	tmpl := content[capturedStart:capturedEnd]

	i := 0
	for i < len(tmpl) {
		dollarIdx := strings.Index(tmpl[i:], "${")
		if dollarIdx < 0 {
			addStaticSegment(content, tmpl[i:], capturedStart+i, &results)
			break
		}

		addStaticSegment(content, tmpl[i:i+dollarIdx], capturedStart+i, &results)

		exprStart := i + dollarIdx + 2
		exprEnd := findClosingBrace(tmpl, exprStart)
		if exprEnd < 0 {
			break
		}

		expr := tmpl[exprStart:exprEnd]
		absExprStart := capturedStart + exprStart
		extractStringLiteralsFromExpr(content, absExprStart, expr, &results)

		i = exprEnd + 1
	}

	return results
}

// addStaticSegment adds a ClassMatch for a static template literal segment
// if it contains at least one class-like token.
func addStaticSegment(content string, segment string, absStart int, results *[]ClassMatch) {
	trimmed := strings.TrimSpace(segment)
	if trimmed == "" || !containsClassToken(trimmed) {
		return
	}
	*results = append(*results, ClassMatch{
		Value:  trimmed,
		Index:  absStart,
		Length: len(segment),
		Full:   segment,
	})
}

// findClosingBrace finds the matching } for an interpolation expression,
// handling nested braces and skipping over quoted strings.
func findClosingBrace(tmpl string, start int) int {
	depth := 1
	i := start
	for i < len(tmpl) && depth > 0 {
		ch := tmpl[i]
		switch {
		case ch == '{':
			depth++
			i++
		case ch == '}':
			depth--
			if depth == 0 {
				return i
			}
			i++
		case ch == '\'' || ch == '"':
			i = skipQuotedString(tmpl, i)
		case ch == '`':
			i = skipTemplateLiteral(tmpl, i)
		default:
			i++
		}
	}
	return -1
}

// skipQuotedString advances past a single or double quoted string.
func skipQuotedString(s string, start int) int {
	quote := s[start]
	i := start + 1
	for i < len(s) && s[i] != quote {
		if s[i] == '\\' {
			i++
		}
		i++
	}
	if i < len(s) {
		i++
	}
	return i
}

// skipTemplateLiteral advances past a nested template literal (backtick string).
func skipTemplateLiteral(s string, start int) int {
	i := start + 1
	for i < len(s) && s[i] != '`' {
		if s[i] == '\\' {
			i++
		} else if i+1 < len(s) && s[i] == '$' && s[i+1] == '{' {
			i += 2
			depth := 1
			for i < len(s) && depth > 0 {
				if s[i] == '{' {
					depth++
				} else if s[i] == '}' {
					depth--
				}
				i++
			}
			continue
		}
		i++
	}
	if i < len(s) {
		i++
	}
	return i
}

// extractStringLiteralsFromExpr finds quoted string literals within an interpolation
// expression and adds ClassMatch entries for those that contain class-like tokens.
func extractStringLiteralsFromExpr(content string, absExprStart int, expr string, results *[]ClassMatch) {
	i := 0
	for i < len(expr) {
		ch := expr[i]
		if ch != '\'' && ch != '"' {
			i++
			continue
		}

		quote := ch
		strStart := i + 1
		j := strStart
		for j < len(expr) && expr[j] != quote {
			if expr[j] == '\\' {
				j++
			}
			j++
		}
		if j >= len(expr) {
			i++
			continue
		}

		strContent := expr[strStart:j]
		trimmed := strings.TrimSpace(strContent)
		if trimmed != "" && containsClassToken(trimmed) {
			absStart := absExprStart + i
			full := expr[i : j+1]
			*results = append(*results, ClassMatch{
				Value:  trimmed,
				Index:  absStart,
				Length: len(full),
				Full:   full,
			})
		}
		i = j + 1
	}
}

// extractMultiArgCalls finds cn() and classNames() calls and extracts string literal
// arguments beyond the first one (which is already handled by classPatterns).
func extractMultiArgCalls(content string, seen map[string]struct{}, results *[]ClassMatch) {
	for _, loc := range reFuncCall.FindAllStringIndex(content, -1) {
		argsStart := loc[1]
		argsEnd := findMatchingParen(content, argsStart)
		if argsEnd < 0 {
			continue
		}

		argsContent := content[argsStart:argsEnd]
		extractFuncStringArgs(content, argsStart, argsContent, seen, results)
	}
}

// findMatchingParen finds the closing ) starting after the opening (.
func findMatchingParen(content string, start int) int {
	depth := 1
	i := start
	for i < len(content) && depth > 0 {
		ch := content[i]
		switch {
		case ch == '(':
			depth++
			i++
		case ch == ')':
			depth--
			if depth == 0 {
				return i
			}
			i++
		case ch == '\'' || ch == '"':
			i = skipQuotedString(content, i)
		case ch == '`':
			i = skipTemplateLiteral(content, i)
		default:
			i++
		}
	}
	return -1
}

// extractFuncStringArgs extracts string literal arguments from a function call's
// argument list, skipping the first argument (already handled by regex patterns).
func extractFuncStringArgs(content string, argsAbsStart int, argsContent string, seen map[string]struct{}, results *[]ClassMatch) {
	argIndex := 0
	i := 0
	for i < len(argsContent) {
		ch := argsContent[i]

		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			i++
			continue
		}

		if ch == ',' {
			argIndex++
			i++
			continue
		}

		if ch == '\'' || ch == '"' {
			quote := ch
			strStart := i + 1
			j := strStart
			for j < len(argsContent) && argsContent[j] != quote {
				if argsContent[j] == '\\' {
					j++
				}
				j++
			}
			if j >= len(argsContent) {
				break
			}

			if argIndex > 0 {
				strContent := argsContent[strStart:j]
				trimmed := strings.TrimSpace(strContent)
				if trimmed != "" && containsClassToken(trimmed) {
					absStart := argsAbsStart + i
					full := argsContent[i : j+1]
					key := fmt.Sprintf("%d:%d", absStart, len(full))
					if _, ok := seen[key]; !ok {
						seen[key] = struct{}{}
						*results = append(*results, ClassMatch{
							Value:  trimmed,
							Index:  absStart,
							Length: len(full),
							Full:   full,
						})
					}
				}
			}
			i = j + 1
			continue
		}

		if ch == '`' {
			i = skipTemplateLiteral(argsContent, i) - 0
			continue
		}

		if ch == '(' {
			end := findMatchingParen(argsContent, i+1)
			if end >= 0 {
				i = end + 1
			} else {
				i++
			}
			continue
		}

		for i < len(argsContent) && argsContent[i] != ',' && argsContent[i] != ')' {
			i++
		}
	}
}

// containsClassToken checks if a string contains at least one token that
// looks like a Tailwind CSS class name.
func containsClassToken(s string) bool {
	for _, token := range strings.Fields(s) {
		if isClassLikeToken(token) {
			return true
		}
	}
	return false
}

// isClassLikeToken checks if a token looks like a CSS/Tailwind class name.
func isClassLikeToken(token string) bool {
	if len(token) == 0 {
		return false
	}
	for _, ch := range token {
		if !isClassChar(ch) {
			return false
		}
	}
	return true
}

func isClassChar(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') ||
		ch == '-' || ch == '/' || ch == '[' || ch == ']' || ch == ':' || ch == '.' || ch == '_' || ch == '!'
}
