package source

import "testing"

var sourceTests = []struct {
	name    string
	message string
	errors  []Error
}{
	{
		"single message",
		"a.go:123:6: expected '(', found 'IDENT' decodeLoggedInUserRequest (and 2 more errors)",
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
		"single message no column",
		"a.go:123: expected '(', found 'IDENT' decodeLoggedInUserRequest (and 2 more errors)",
		[]Error{
			{
				"a.go",
				123,
				-1,
				"expected '(', found 'IDENT' decodeLoggedInUserRequest (and 2 more errors)",
			},
		},
	},
	{
		"multiple messages",
		`a.go:123:12: expected '(', found 'IDENT' decodeLoggedInUserRequest (and 2 more errors)
    b.go:123:1122: expected '(', found 'IDENT' decodeLoggedInUserRequest (and 2 more errors)`,
		[]Error{
			{
				"a.go",
				123,
				12,
				"expected '(', found 'IDENT' decodeLoggedInUserRequest (and 2 more errors)",
			},
			{
				"b.go",
				123,
				1122,
				"expected '(', found 'IDENT' decodeLoggedInUserRequest (and 2 more errors)",
			},
		},
	},
}

func TestParse(t *testing.T) {
	for _, test := range sourceTests {
		t.Run(test.name, func(t *testing.T) {
			errs := ParseSourceErrors(test.message)
			if len(errs) != len(test.errors) {
				t.Logf("was expecting %d errors got %d instead", len(test.errors), len(errs))
				t.Fail()
			}
			for i, actual := range errs {
				expected := test.errors[i]
				if actual.File != expected.File {
					t.Logf("was expecting file %q got %q instead", expected.File, actual.File)
					t.Fail()
				}
				if actual.Line != expected.Line {
					t.Logf("was expecting line %d got %d instead", expected.Line, actual.Line)
					t.Fail()
				}
				if actual.Column != expected.Column {
					t.Logf("was expecting column %d got %d instead", expected.Column, actual.Column)
					t.Fail()
				}
				if actual.Message != expected.Message {
					t.Logf("was expecting message %q got %q instead", expected.Message, actual.Message)
					t.Fail()
				}
			}
		})
	}
}
