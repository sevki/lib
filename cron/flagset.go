package cron

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// NewFlagSet returns a new cron flagset
func NewFlagSet(name string, errorhandling flag.ErrorHandling) *FlagSet {
	return &FlagSet{
		flag.FlagSet{},
		os.Args[1:], //
		Parser{},
	}
}

// FlagSet is a cron FlagSet
type FlagSet struct {
	flag.FlagSet
	args []string
	p    Parser
}

// Parse parses the strings in to a cronexpr flagset
func (f *FlagSet) Parse(args ...string) error {
	if len(args) <= 0 {
		args = f.args
	} else {
		f.args = args
	}
	if err := f.p.Parse(args); err != nil {
		return err
	}
	return f.FlagSet.Parse(args[5:])
}

// CronTab2 is pretty minimal, it prints a table but doen't explain stuff
func (f *FlagSet) CronTab2() string {
	const padding = 3
	colNames := []string{
		"minute",
		"hour",
		"day of the month",
		"month",
		"day of the week",
		"args",
	}
	buf := bytes.Buffer{}
	w := tabwriter.NewWriter(&buf, 0, 0, padding, 0x20, tabwriter.StripEscape)
	rows := f.p.intervals
	rows = append(rows, f.Args())
	for col, ints := range rows {
		fmt.Fprintf(w, "%s\t%s\t\n", colNames[col], ints)
	}
	w.Flush()
	return buf.String()
}

// CronTab prints a very nice explanatinon of what things do
func (f *FlagSet) CronTab() string {
	const padding = 1
	colNames := []string{
		"minutes",
		"hours",
		"days of the month",
		"months",
		"days of the week",
		"args",
	}
	buf := bytes.Buffer{}
	w := tabwriter.NewWriter(&buf, 0, 0, padding, ' ', tabwriter.StripEscape)
	rows := f.p.intervals
	rows = append(rows, f.Args())
	for col := range rows {
		fmt.Fprintf(w, "%s%s\n",
			strings.Repeat("│\t", col),
			"┌")
	}
	fmt.Fprintf(w, "%s", strings.Join(f.args, "\t"))
	w.Flush()
	headers := strings.Split(buf.String(), "\n")
	buf.Reset()
	w = tabwriter.NewWriter(&buf, 0, 0, padding, '-', 0)
	for col, vals := range rows {
		colName := colNames[col]
		pad := 17 - len(colName)
		fmt.Fprintf(w, "%s\t%s %s\n", headers[col], colName+strings.Repeat(" ", pad), vals)
	}

	w.Flush()
	table := strings.ReplaceAll(buf.String(), "-", "─") // ─ is greater than 127 there for not a byte
	return table + headers[len(headers)-1]
}
