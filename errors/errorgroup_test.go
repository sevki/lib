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
			err:  New("group1").Add(errors.New("err-1")),
		},
		{
			name: "twoLevels",
			err: New("group1").Add(
				New("group2").Add(errors.New("err-2")),
			),
		},
		{
			name: "treeLevels",
			err: New("group1").Add(
				New("group2").Add(
					New("group3").Add(errors.New("err-3")),
				),
			).Add(New("group4").Add(
				New("group5").Add(
					New("group6").New("err-4"),
				),
			)),
		},
		{
			name: "nil",
			err: New("group1").Add(err).Add(
				New("group2").Add(
					New("group3").Add(errors.New("err-3")),
				),
			).Add(New("group4").Add(
				New("group5").Add(
					New("group6").New("err-4"),
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
