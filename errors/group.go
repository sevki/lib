// Copyright 2018 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errors provides grouped errors.
package errors // import "sevki.org/x/errors"

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

	g.printError(buf, []string{})
	x := buf.String()
	if strings.HasSuffix(x, "\n") {
		return x[:len(x)-1]
	}
	return x
}

func (g *Group) printError(w io.Writer, prefixes []string) {
	for i, err := range g.errs {
		padding := strings.Join(append(prefixes, g.prefix), ": ")
		spacePadding := strings.Repeat(" ", len(padding))

		switch x := err.(type) {
		case *Group:
			x.printError(w, append(prefixes, g.prefix))
		case error:
			a := strings.Split(err.Error(), "\n")
			for j, line := range a {
				sep := ":"
				p := padding
				if i == 0 {
					// no op
				} else {
					switch g.errs[i-1].(type) {
					case *Group:
						// no op
					default:
						sep = " "
						p = spacePadding
					}
				}
				if j > 0 {
					sep = "⮎"
					p = spacePadding
				}
				fmt.Fprintf(w, "%s%s %s\n", p, sep, line)
			}
		}
	}
}
