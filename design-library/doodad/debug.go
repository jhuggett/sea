package doodad

import (
	"fmt"
	"strings"
)

func (c *Children) PrettyPrint(indent int) string {
	var builder strings.Builder
	prefix := strings.Repeat("  ", indent)

	// ANSI color codes
	colorReset := "\033[0m"
	colorBlue := "\033[34m"
	colorGreen := "\033[32m"
	colorYellow := "\033[33m"
	colorCyan := "\033[36m"
	colorMagenta := "\033[35m"
	colorRed := "\033[31m"

	// Print parent information at the root level
	if indent == 0 {
		if c.Parent != nil {
			parentType := fmt.Sprintf("%T", c.Parent)

			// Get parent layout information (same format as children)
			layout := c.Parent.Layout()
			layoutInfo := ""
			if layout != nil {
				layoutInfo = fmt.Sprintf("%spos(%d,%d) size(%d×%d)%s %slayout(%p)%s",
					colorMagenta, layout.X(), layout.Y(), layout.Width(), layout.Height(), colorReset,
					colorRed, layout, colorReset)
			}

			builder.WriteString(fmt.Sprintf("%sParent: %s%s%s %s\n",
				colorBlue, colorGreen, parentType, colorReset, layoutInfo))
		} else {
			builder.WriteString(fmt.Sprintf("%sParent: %s<no parent>%s\n",
				colorBlue, colorRed, colorReset))
		}
	}

	for i, doodad := range c.Doodads {
		typeStr := fmt.Sprintf("%T", doodad)

		// Get layout information
		layout := doodad.Layout()
		layoutInfo := ""
		if layout != nil {
			layoutInfo = fmt.Sprintf("%spos(%d,%d) size(%d×%d)%s %slayout(%p)%s",
				colorMagenta, layout.X(), layout.Y(), layout.Width(), layout.Height(), colorReset,
				colorRed, layout, colorReset)
		}

		// Add horizontal pipes for padding
		indentStr := prefix
		if indent > 0 {
			indentStr = strings.Repeat(colorBlue+"│ "+colorReset, indent-1) + colorBlue + "├─" + colorReset
		}

		builder.WriteString(fmt.Sprintf("%s %s%d%s: %s%s%s %s\n",
			indentStr,
			colorYellow, i, colorReset,
			colorGreen, typeStr, colorReset,
			layoutInfo))

		// Recursively print children
		if doodad.Children() != nil && len(doodad.Children().Doodads) > 0 {
			builder.WriteString(doodad.Children().PrettyPrint(indent + 1))
		}
	}

	// For root level, add a header
	if indent == 0 {
		result := builder.String()
		result = strings.TrimSuffix(result, "\n")
		fmt.Printf("%s┌── Doodad Tree ───%s\n%s\n%s└─────────────────%s\n",
			colorCyan, colorReset,
			result,
			colorCyan, colorReset)
	}

	return builder.String()
}
