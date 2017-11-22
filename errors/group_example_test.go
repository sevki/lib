package errors

import (
	"fmt"
)

func ExampleNewGroup() {
	g := NewGroup("group1").Add(
		NewGroup("group2").Add(
			NewGroup("group3").Add(New("err")).Add(New("err2")).Add(New(`this is a multiline error
no but really`)),
		).New("err3"),
	)

	fmt.Println(g.Error())
	// Output:
	// group1: group2: group3: err
	//                         err2
	//                         this is a multiline error
	//                       ↪ no but really
	//group1: group2: err3
}
