package errors

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
)

type levelederror struct {
	level int
	err   error
}

// Group defines a set errors
type Group struct {
	prefix string

	mu   sync.Mutex
	errs []error
}

// New creates a new error
var New = errors.New

// NewGroup given a prefix will return an Group with the given prefix.
// Group formats errors like so;
//
//  {prefix}: {err msg}
func NewGroup(prefix string) *Group {
	return &Group{
		prefix: prefix,
		errs:   make([]error, 0),
	}
}

// Add adds a given error to the Group
func (g *Group) Add(err error) *Group {
	if err != nil {
		g.mu.Lock()
		g.errs = append(g.errs, err)
		g.mu.Unlock()
	}
	return g
}

// Newf creates a new error with formatting and adds the new error to the Group
func (g *Group) Newf(format string, a ...interface{}) *Group {
	return g.Add(fmt.Errorf(format, a...))
}

// New creates a new error and adds the new error to the Group
func (g *Group) New(s string) *Group {
	return g.Add(errors.New(s))
}

// Errored returns true if an error has been added and false if
// no errors have been added.
func (g *Group) Errored() bool {
	if g == nil {
		panic("group is nil")
	}
	if len(g.errs) < 1 {
		return false
	}
	for _, child := range g.errs {
		switch err := child.(type) {
		case *Group:
			if err.Errored() {
				return true
			}
		default:
			return true
		}
	}
	return false
}

// Error implements the error interface
func (g *Group) Error() string {
	buf := bytes.NewBuffer(nil)

	g.printError(buf, 0)

	return buf.String()
}

func (g *Group) printError(w io.Writer, level int) {
	padding := strings.Repeat("\t", level)
	for i, err := range g.errs {
		if i < 1 {
			fmt.Fprintf(w, "%s%s:\n", padding, g.prefix)
		}
		switch x := err.(type) {
		case *Group:
			x.printError(w, level+1)
		case error:
			fmt.Fprintf(w, "%s%s\n", strings.Repeat("\t", level+1), err.Error())
		default:
		}
	}
}
