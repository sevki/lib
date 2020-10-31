package cron

import (
	"flag"
	"strings"
	"testing"
)

func TestCron(t *testing.T) {
	tests := []struct {
		exp string
	}{
		{
			exp: `15 3 * * 1-5`,
		},
		{
			exp: `*/15 0 1,15 * 1-5`,
		},
	}
	for _, test := range tests {
		t.Run(test.exp, func(t *testing.T) {
			fset := NewFlagSet("test", flag.PanicOnError)
			fset.Parse(strings.Fields(test.exp)...)
		})
	}

}

func TestLexCronExpression(t *testing.T) {
	tests := []struct {
		expression string
		vals       []int
	}{
		{
			"*",
			[]int{firstlast},
		},
		{
			"1",
			[]int{1},
		},
		{
			"1,2",
			[]int{1, 2},
		},
		{
			"1,2-3",
			[]int{1, 2, rng, 3},
		},
		{
			"1,2-4",
			[]int{1, 2, rng, 4},
		},
		{
			"1,2-4,6",
			[]int{1, 2, rng, 4, 6},
		},
	}
	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			vals := parseCronexp(test.expression)
			fail := func() {
				t.Logf("was expecting %d but got %d instead", test.vals, vals)
				t.Fail()
			}
			if len(vals) != len(test.vals) {
				fail()
				return
			}
			for i, u := range vals {
				if u != test.vals[i] {
					fail()
				}
			}
		})
	}
}
