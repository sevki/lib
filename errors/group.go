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

var New = errors.New

// NewGroup given a prefix will return an Group with the given prefix.
// Group formats errors like so;
//
//  {prefix}: {err msg}
func NewGroup(prefix string) *Group {
	return &Group{
		prefix: prefix,
	}
}

// Add adds a given error to the Group
func (e *Group) Add(err error) *Group {
	if err == nil {
		return e
	}

	e.mu.Lock()
	e.errs = append(e.errs, err)
	e.mu.Unlock()
	return e
}

// New creates a new error and adds the new error to the Group
func (e *Group) New(err string) *Group {
	e.mu.Lock()
	e.errs = append(e.errs, errors.New(err))
	e.mu.Unlock()
	return e
}

// Errored returns true if an error has been added and false if
// no errors have been added.
func (e *Group) Errored() bool { return len(e.errs) > 0 }

// Error implements the error interface
func (e *Group) Error() string {
	buf := bytes.NewBuffer(nil)

	e.printError(buf, 0)
	return buf.String()
}

func (e *Group) printError(w io.Writer, level int) {
	padding := strings.Repeat("\t", level)
	fmt.Fprintf(w, "%s%s:\n", padding, e.prefix)
	for _, err := range e.errs {
		switch x := err.(type) {
		case *Group:
			x.printError(w, level+1)
		case error:
			fmt.Fprintf(w, "%s%s\n", strings.Repeat("\t", level+1), err.Error())
		default:
		}
	}
}
