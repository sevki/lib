// Copyright 2019 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cron implements Crontab specification as defined in
//
// The Open Group Base Specifications Issue 7, 2018 edition
// IEEE Std 1003.1-2017 (Revision of IEEE Std 1003.1-2008)
//
// https://pubs.opengroup.org/onlinepubs/9699919799/utilities/crontab.html
//
// In addition to the standard, this package also implements the */15 syntax
// I am thinking about removing this as the step syntax doesn't make that
// much sense to me but the implementation is only an additional stack
// so it's ok to stay for now.
package cron

import (
	"fmt"
	"strconv"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	eof       = -1
	step      = -2
	firstlast = -3
	rng       = -4
	etx       = rune(3)
)

type stateFn func(*Parser) stateFn

// Parser is a cronexp parser.
type Parser struct {
	state     stateFn
	arg       string
	intervals []interface{}
	err       error
}

// Parse parses args. These are usually os.Args
func (p *Parser) Parse(args []string) error {
	for p.state = parseMinute; p.state != nil; {
		p.arg, args = args[0], args[1:]
		p.state = p.state(p)
	}
	return p.err
}

func parseCronexp(s string) (vals []int) {
	for {
		advance := func(size int) { s = s[size:] }
		next := func() rune {
			r, size := utf8.DecodeRuneInString(s)
			if size == 0 {
				return etx
			}
			return r
		}
		switch next() {
		case ',':
			advance(1)
		case '*':
			vals = append(vals, firstlast)
			advance(1)
		case '-':
			advance(1)
			vals = append(vals, rng)
		case '/':
			advance(1)
			vals = append(vals, step)
		case etx:
			return
		default:
			u, s, _ := readNumber(s)
			// err dropped here
			// TODO(sevki): propogate these errors up to the parser
			advance(s)
			vals = append(vals, u)
		}
	}
}

func readNumber(s string) (int, int, error) {
	buf := ""
	for r, size := utf8.DecodeRuneInString(s); unicode.IsDigit(r); r, size = utf8.DecodeRuneInString(s) {
		buf += s[:size]
		s = s[size:]
	}
	u, err := strconv.Atoi(buf)
	if err != nil {
		return -1, -1, err
	}
	return u, len(buf), nil
}

type intStack []int

func (s *intStack) push(x int) {
	*s = append(*s, x)
}

func (s *intStack) pop() int {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

func evalCronexp(s string, start, end int, transform func(int) interface{}) []interface{} {
	tokens := parseCronexp(s)
	intervals := intStack{}
	next := func() int {
		if len(tokens) == 0 {
			return eof
		}
		i := tokens[0]
		tokens = tokens[1:]
		return i
	}
	fill := func(from, to int) {
		for i := from; i <= to; i++ {
			intervals.push(i)
		}
	}
	for t := next(); t != eof; t = next() {
		switch t {
		case firstlast:
			fill(start, end)
		case rng:
			fill(intervals.pop(), next())
		case step:
			// https://en.wikipedia.org/wiki/Cron#Non-standard_characters
			// this isn't standard and makes no sense
			// TODO(sevki): delete '/'
			nth := next() // next token should be the step value
			nintervals := []int{}
			for i := 0; i < len(intervals); i += nth {
				nintervals = append(nintervals, intervals[i])
			}
			intervals = nintervals
		default:
			intervals.push(t)
		}
	}
	f := []interface{}{}

	for _, v := range intervals {
		f = append(f, transform(v))
	}
	return f
}

func parseMinute(p *Parser) stateFn {
	p.intervals = append(p.intervals, evalCronexp(p.arg, 0, 59, func(i int) interface{} {
		return time.Duration(i) * time.Minute
	}))
	return parseHour
}
func parseHour(p *Parser) stateFn {
	p.intervals = append(p.intervals, evalCronexp(p.arg, 0, 59, func(i int) interface{} {
		return time.Duration(i) * time.Hour
	}))
	return parseDayOfMonth
}

type dayofthemonth int

func (d dayofthemonth) String() string { j := int(d); return fmt.Sprint(j) }

func parseDayOfMonth(p *Parser) stateFn {
	p.intervals = append(p.intervals, evalCronexp(p.arg, 1, 31, func(i int) interface{} {
		return dayofthemonth(i)
	}))
	return parseMonth
}
func parseMonth(p *Parser) stateFn {
	p.intervals = append(p.intervals, evalCronexp(p.arg, 1, 12, func(i int) interface{} {
		return time.Month(i)
	}))
	return parseDayOfTheWeek
}

func parseDayOfTheWeek(p *Parser) stateFn {
	p.intervals = append(p.intervals, evalCronexp(p.arg, 0, 6, func(i int) interface{} {
		return time.Weekday(i)
	}))
	return nil
}
