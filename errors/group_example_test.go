package errors

import (
	"fmt"
)

func ExampleNewGroup() {
	g := NewGroup("group1").Add(
		NewGroup("group2").Add(
			NewGroup("group3").Add(New("group before this")).Add(New("space before this")).Add(New(`space before this
carriage return before this`)),
		).New(`group before this
carriage return before this
carriage return before this
carriage return before this `),
	)

	fmt.Println(g.Error())
	// Output:
	// group1: group2: group3: group before this
	//                         space before this
	//                         space before this
	//                       ⮎ carriage return before this
	// group1: group2: group before this
	//               ⮎ carriage return before this
	//               ⮎ carriage return before this
	//               ⮎ carriage return before this

}
