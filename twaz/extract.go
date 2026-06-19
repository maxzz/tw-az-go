package twaz

import "fmt"

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

			key := fmt.Sprintf("%d:%d", loc[0], loc[1]-loc[0])
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}

			results = append(results, ClassMatch{
				Value:  content[capturedStart:capturedEnd],
				Index:  loc[0],
				Length: loc[1] - loc[0],
				Full:   content[loc[0]:loc[1]],
			})
		}
	}

	return results
}
