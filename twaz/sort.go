package twaz

import (
	"sort"
	"strings"
)

// CheckClassString returns order violations for a class string.
func CheckClassString(value string) []OrderViolation {
	tokens := splitTokens(value)
	var violations []OrderViolation
	maxGroup := -1

	for _, token := range tokens {
		group := Classify(token)
		if group < 0 {
			continue
		}
		if group < maxGroup {
			violations = append(violations, OrderViolation{
				Token: token,
				Group: groupName(group),
				After: groupName(maxGroup),
			})
		}
		if group > maxGroup {
			maxGroup = group
		}
	}

	return violations
}

// SortClassString reorders utility classes according to twaz rules.
func SortClassString(value string) string {
	tokens := splitTokens(value)
	if len(tokens) < 2 {
		return value
	}

	type tokenMeta struct {
		token     string
		sortGroup int
		index     int
	}

	withMeta := make([]tokenMeta, len(tokens))
	for i, token := range tokens {
		group := Classify(token)
		sortGroup := group
		if sortGroup < 0 {
			sortGroup = unknownGroup
		}
		withMeta[i] = tokenMeta{token: token, sortGroup: sortGroup, index: i}
	}

	sort.SliceStable(withMeta, func(i, j int) bool {
		a, b := withMeta[i], withMeta[j]
		if a.sortGroup != b.sortGroup {
			return a.sortGroup < b.sortGroup
		}
		return a.index < b.index
	})

	sorted := make([]string, len(withMeta))
	for i, item := range withMeta {
		sorted[i] = item.token
	}
	return strings.Join(sorted, " ")
}

func splitTokens(value string) []string {
	raw := strings.Fields(value)
	if len(raw) == 0 {
		return nil
	}
	return raw
}
