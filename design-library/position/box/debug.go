package box

import (
	"fmt"
	"strings"
)

// String implements the Stringer interface for pretty printing a Box
func (b *Box) String() string {
	// return b.prettyPrint(0, map[*Box]bool{})

	return ""
}

// prettyPrint returns a formatted string representation of the Box with indentation and colors
func (b *Box) prettyPrint(level int, visited map[*Box]bool) string {
	if visited[b] {
		return fmt.Sprintf("\033[33m%s└─ Box %p (already printed)\033[0m", strings.Repeat("  ", level), b)
	}
	visited[b] = true

	indent := strings.Repeat("  ", level)
	nextIndent := indent + "  "

	// Color constants for better readability
	cyan := "\033[36m"
	green := "\033[32m"
	blue := "\033[34m"
	magenta := "\033[35m"
	reset := "\033[0m"

	// Build the output string
	var sb strings.Builder

	// Header with box address
	sb.WriteString(fmt.Sprintf("%sBox %p%s {\n", cyan, b, reset))

	// Position and dimensions
	sb.WriteString(fmt.Sprintf("%s%sPosition: (%d, %d)%s\n", nextIndent, green, b.x, b.y, reset))
	sb.WriteString(fmt.Sprintf("%s%sDimensions: %dx%d%s\n", nextIndent, green, b.width, b.height, reset))

	// Calculation info
	sb.WriteString(fmt.Sprintf("%s%sCalculation Steps: %d%s\n", nextIndent, blue, len(b.calculationSteps), reset))
	sb.WriteString(fmt.Sprintf("%s%sRecalculation Count: %d%s\n", nextIndent, blue, b.recalculationCount, reset))

	// Dependents section
	sb.WriteString(fmt.Sprintf("%s%sDependents (%d):%s", nextIndent, magenta, len(b.dependents), reset))

	if len(b.dependents) == 0 {
		sb.WriteString(" none\n")
	} else {
		sb.WriteString("\n")
		for i, dependent := range b.dependents {
			prefix := "├─ "
			if i == len(b.dependents)-1 {
				prefix = "└─ "
			}
			sb.WriteString(fmt.Sprintf("%s  %s%s\n", nextIndent, prefix, dependent.prettyPrint(level+2, visited)))
		}
	}

	sb.WriteString(fmt.Sprintf("%s}", indent))
	return sb.String()
}
