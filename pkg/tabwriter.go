package pkg

import (
	"io"
	"text/tabwriter"
)

const (
	tabwriterMinWidth = 0
	tabwriterPadding  = 3
	tabwriterWidth    = 2
	tabwriterPadChar  = ' '
)

// NewTabWriter creates a tabwriter with default configurations to align
// input text into tab-spaced columns.
func NewTabWriter(out io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(out, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, 0)
}
