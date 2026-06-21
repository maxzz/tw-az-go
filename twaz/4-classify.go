package twaz

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	reEnd           = regexp.MustCompile(`^(?:cursor-|pointer-events)`)
	reZIndex        = regexp.MustCompile(`^z-`)
	rePosition      = regexp.MustCompile(`^(?:relative|absolute|fixed|sticky|static)`)
	rePositionOff   = regexp.MustCompile(`^(?:inset-|top-|right-|bottom-|left-)`)
	reSelfGroup     = regexp.MustCompile(`^self-`)
	reGroupSlash    = regexp.MustCompile(`^group/`)
	reTransition    = regexp.MustCompile(`^(?:transition|duration|animate)`)
	reOverflow      = regexp.MustCompile(`^overflow-`)
	reGrid          = regexp.MustCompile(`^grid(?:-|$)`)
	reInlineGrid    = regexp.MustCompile(`^inline-grid`)
	reFlex          = regexp.MustCompile(`^flex(?:-|$)`)
	reInlineFlex    = regexp.MustCompile(`^inline-flex`)
	reGap           = regexp.MustCompile(`^gap-`)
	reItems         = regexp.MustCompile(`^items-`)
	reJustify       = regexp.MustCompile(`^justify-`)
	reContent       = regexp.MustCompile(`^content-`)
	rePlace         = regexp.MustCompile(`^place-`)
	reOrder         = regexp.MustCompile(`^order-`)
	reCol           = regexp.MustCompile(`^col-`)
	reRow           = regexp.MustCompile(`^row-`)
	reSpace         = regexp.MustCompile(`^space-[xy]-`)
	reList          = regexp.MustCompile(`^list-`)
	reShrinkGrow    = regexp.MustCompile(`^(?:shrink|grow)$`)
	reShrink        = regexp.MustCompile(`^shrink-`)
	reGrow          = regexp.MustCompile(`^grow-`)
	reBasis         = regexp.MustCompile(`^basis-`)
	reSelect        = regexp.MustCompile(`^select-`)
	reWhitespace    = regexp.MustCompile(`^whitespace-`)
	reSpacing       = regexp.MustCompile(`^(?:m-|mx-|my-|mt-|mr-|mb-|ml-|p-|px-|py-|pt-|pr-|pb-|pl-)`)
	reDimensions    = regexp.MustCompile(`^(?:w-|h-|min-w-|max-w-|min-h-|max-h-|size-|aspect-)`)
	reFont          = regexp.MustCompile(`^font-`)
	reBorder        = regexp.MustCompile(`^(?:border|outline-|ring-|divide-)`)
	reRounded       = regexp.MustCompile(`^rounded`)
	reShadow        = regexp.MustCompile(`^shadow`)
	reTextColor     = regexp.MustCompile(`^text-`)
	reBgFill        = regexp.MustCompile(`^(?:bg-|fill-|stroke-|from-|to-|via-|opacity-|accent-|caret-|decoration-)`)
	reGroupNamed    = regexp.MustCompile(`^group/[\w-]+$`)
)

// Classify returns the sort group for a Tailwind utility token, or -1 if unknown.
func Classify(token string) int {
	base := baseToken(token)
	variant := hasVariantPrefix(token)

	if reEnd.MatchString(base) || base == "pointer-events-none" || reZIndex.MatchString(base) {
		return endGroup
	}

	if rePosition.MatchString(base) {
		if variant {
			return variantGroup
		}
		return 0
	}

	if rePositionOff.MatchString(base) {
		if variant {
			return variantGroup
		}
		return 1
	}

	if reSelfGroup.MatchString(base) || base == "group" || reGroupSlash.MatchString(base) {
		return 2
	}

	if reTransition.MatchString(base) {
		return transitionGroup
	}

	if base == "truncate" || base == "text-ellipsis" || reOverflow.MatchString(base) || reSelect.MatchString(base) {
		return truncateOverflowGrp
	}

	if reGrid.MatchString(base) || base == "grid" || reInlineGrid.MatchString(base) ||
		reFlex.MatchString(base) || base == "flex" || reInlineFlex.MatchString(base) ||
		reGap.MatchString(base) || reItems.MatchString(base) || reJustify.MatchString(base) ||
		reContent.MatchString(base) || rePlace.MatchString(base) || reOrder.MatchString(base) ||
		reCol.MatchString(base) || reRow.MatchString(base) || reSpace.MatchString(base) ||
		reList.MatchString(base) {
		if variant {
			return variantGroup
		}
		return childrenGroup
	}

	if reShrinkGrow.MatchString(base) || reShrink.MatchString(base) || reGrow.MatchString(base) ||
		reBasis.MatchString(base) || reWhitespace.MatchString(base) ||
		base == "compress-zero" {
		return 3
	}

	if !variant && reSpacing.MatchString(base) {
		return 4
	}

	if reDimensions.MatchString(base) {
		return 5
	}

	if base == "block" || base == "inline" || base == "hidden" || base == "visible" || base == "isolate" {
		if variant {
			return variantGroup
		}
		return 6
	}

	if _, ok := textSizes[base]; ok {
		return 7
	}

	if reFont.MatchString(base) {
		return 8
	}

	if reBorder.MatchString(base) && !reRounded.MatchString(base) {
		if variant {
			return variantGroup
		}
		return borderGroup
	}

	if reRounded.MatchString(base) {
		if variant {
			return variantGroup
		}
		return roundingGroup
	}

	if reShadow.MatchString(base) {
		if variant {
			return variantGroup
		}
		return shadowGroup
	}

	if !variant && reTextColor.MatchString(base) {
		if _, isSize := textSizes[base]; !isSize {
			return 9
		}
	}

	if !variant && reBgFill.MatchString(base) {
		return 10
	}

	if variant {
		return variantGroup
	}

	return -1
}

func baseToken(token string) string {
	if reGroupNamed.MatchString(token) {
		return token
	}
	parts := strings.Split(token, ":")
	return parts[len(parts)-1]
}

func hasVariantPrefix(token string) bool {
	if reGroupNamed.MatchString(token) {
		return false
	}
	return strings.Contains(token, ":")
}

func groupName(group int) string {
	if group >= 0 && group < len(groupNames) {
		return groupNames[group]
	}
	return "group " + strconv.Itoa(group)
}
