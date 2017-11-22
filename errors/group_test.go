package errors

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var (
	update = flag.Bool("u", false, "-u")
)

func TestNilReciever(t *testing.T) {
	var g *Group
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()
	g.Errored()
}

func TestNewF(t *testing.T) {

	g := NewGroup("group1").Newf("%d", 1)

	s := g.Error()
	x := "group1: 1"
	if s != x {
		t.Logf("%q != %q", g.Error(), x)
		t.Fail()
	}

}

func TestErrored(t *testing.T) {
	a := NewGroup("a")
	b := NewGroup("b")
	c := NewGroup("c")
	d := NewGroup("d")
	e := NewGroup("e")
	f := NewGroup("f")
	g := NewGroup("g")

	// a = {error}
	// b = {}
	// c = {a}
	// d = {b}
	// e = {a,b}
	// f = {d}
	// g = {c}

	a = a.New("hello")
	if !a.Errored() {
		t.Fail()
	}
	if b.Errored() {
		t.Fail()
	}
	c = c.Add(a)
	if !c.Errored() {
		t.Fail()
	}
	d = d.Add(b)
	if d.Errored() {
		t.Fail()
	}
	e = e.Add(a).Add(b)
	if !e.Errored() {
		t.Fail()
	}
	f = f.Add(d)
	if f.Errored() {
		t.Fail()
	}
	g = g.Add(c)
	if !g.Errored() {
		t.Fail()
	}
}

func TestNil(t *testing.T) {
	g := NewGroup("prefix")
	if g.Errored() {
		t.Fail()
	}
	g.Add(nil)
	if g.Errored() {
		t.Fail()
	}
}
func TestErrorGroups(t *testing.T) {
	flag.Parse()
	var err error
	tests := []struct {
		name   string
		err    *Group
		output string
	}{
		{
			name: "single",
			err:  NewGroup("group1").Add(errors.New("err-1")),
		},
		{
			name: "twoLevels",
			err: NewGroup("group1").Add(
				NewGroup("group2").Add(errors.New("err-2")),
			),
		},
		{
			name: "treeLevels",
			err: NewGroup("group1").Add(
				NewGroup("group2").Add(
					NewGroup("group3").Add(errors.New("err-3")),
				),
			).Add(NewGroup("group4").Add(
				NewGroup("group5").Add(
					NewGroup("group6").New("err-4"),
				),
			)),
		},
		{
			name: "nil",
			err: NewGroup("group1").Add(err).Add(
				NewGroup("group2").Add(
					NewGroup("group3").Add(errors.New("err-3")),
				),
			).Add(NewGroup("group4").Add(
				NewGroup("group5").Add(
					NewGroup("group6").New("err-4"),
				),
			)),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gold := "./testdata/" + test.name + ".gold"

			if *update {
				var goldFile *os.File

				goldFile, err := os.Create(gold)
				if err != nil {
					t.Fatal(err)
				}
				io.WriteString(goldFile, test.err.Error())

				goldFile.Close()
				return
			}
			bytz, err := ioutil.ReadFile(gold)
			if err != nil {
				t.Fatal(err)
			}
			dmp := diffmatchpatch.New()

			diffs := dmp.DiffMain(string(bytz), test.err.Error(), false)

			if bytes.Compare(bytz, []byte(test.err.Error())) > 0 {
				t.Log("\n", dmp.DiffPrettyText(diffs))
				t.Fail()
			}
		})
	}
}
