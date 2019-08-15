package pkg

import (
	"fmt"
	"io"
	"sort"

	"github.com/fatih/color"
)

var cliWarnings *Warnings

func init() {
	// Create a new CLI-wide Warnings instance which can be used by
	// any component of the CLI to collect warnings.
	cliWarnings = NewWarnings()
}

// Warnings contains a collection of warnings and their sources.
//
// This is used to collect warnings during the envsnap render. The source
// of the warning is the key, and the string which describes the warning
// itself is kept in the value. A warning source may have multiple warnings
// associated with it.
type Warnings struct {
	Warnings map[string][]string
}

// NewWarnings creates a new instance of a Warnings struct, used to accumulate
// and print warnings.
func NewWarnings() *Warnings {
	return &Warnings{
		Warnings: make(map[string][]string),
	}
}

// Add a new warning.
//
// The source should describe which component/area the warning was generated from.
// The message should be the warning itself, describing what is wrong. The message
// may be a format string, in which case the format components may be passed along
// as well.
func (w *Warnings) Add(src string, msg string, a ...interface{}) {
	_, exists := w.Warnings[src]
	if !exists {
		w.Warnings[src] = []string{fmt.Sprintf(msg, a...)}
	} else {
		w.Warnings[src] = append(w.Warnings[src], fmt.Sprintf(msg, a...))
	}
}

// Clear all collected warnings.
func (w *Warnings) Clear() {
	w.Warnings = make(map[string][]string)
}

// HasWarnings checks to see whether the Warnings instance has any tracked warnings.
func (w *Warnings) HasWarnings() bool {
	return len(w.Warnings) > 0
}

// Print out the tracked warnings.
//
// The printed message is formatted with color and in a tabular view. Warnings
// are sorted alphabetically first by source, then by warning message.
func (w *Warnings) Print(writer io.Writer) {
	if w.HasWarnings() {
		yellow := color.New(color.FgYellow)
		yellow.Fprintln(writer, "------------------------------")
		yellow.Fprintf(writer, "warnings: %d\n", len(w.Warnings))

		tw := NewTabWriter(writer)
		defer tw.Flush()
		fmt.Fprintln(tw)

		var srcs []string
		for src := range w.Warnings {
			srcs = append(srcs, src)
		}
		sort.Strings(srcs)

		for _, src := range srcs {
			warnings := w.Warnings[src]
			sort.Strings(warnings)

			for _, warning := range warnings {
				fmt.Fprintln(tw, yellow.Sprintf("[%s]\t%s", src, warning))
			}
		}
		fmt.Fprintln(tw)
	}
}
