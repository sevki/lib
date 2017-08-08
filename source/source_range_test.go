package source

import "testing"

var sourceRangeTests = []struct {
	name      string
	lines     []int
	tolerance int
	errors    []Error
}{
	{
		"zero tolerance",
		[]int{123},
		0,
		[]Error{
			{
				"a.go",
				123,
				6,
				"expected '(', found 'IDENT' decodeLoggedInUserRequest (and 2 more errors)",
			},
		},
	},
	{
		"1 tolerance",
		[]int{122, 123, 124},
		1,
		[]Error{
			{
				"a.go",
				123,
				6,
				"expected '(', found 'IDENT' decodeLoggedInUserRequest (and 2 more errors)",
			},
		},
	},
	{
		"starting at 1",
		[]int{1, 2},
		1,
		[]Error{
			{
				"a.go",
				1,
				6,
				"expected '(', found 'IDENT' decodeLoggedInUserRequest (and 2 more errors)",
			},
		},
	},
}

func TestGetRanges(t *testing.T) {
	for _, test := range sourceRangeTests {
		t.Run(test.name, func(t *testing.T) {
			lines := GetRangesFromErrors(test.errors, test.tolerance)
			if len(lines) != len(test.lines) {
				t.Logf("was expecting %d got %d instead", test.lines, lines)
				t.FailNow()
			}
			for i, actual := range test.lines {
				expected := lines[i]
				if actual != expected {
					t.Logf("was expecting %d got %d instead", test.lines, lines)
					t.Fail()
				}
			}
		})
	}
}
