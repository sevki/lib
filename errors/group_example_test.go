package errors

import (
	"fmt"
)

func ExampleNewGroup() {
	g := NewGroup("group1").Add(
		NewGroup("group2").Add(
			NewGroup("group3").Add(New("err")),
		))

	fmt.Println(g.Error())
	// Output:
	// group1:
	//	group2:
	//		group3:
	//			err
}
