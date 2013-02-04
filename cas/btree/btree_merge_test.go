package btree

import (
	"fmt"
)

func ExampleBtree_merge_leaf() {
	left := leaf_node("A", "B")
	right := leaf_node("X", "Y")

	merge(left, right, []byte("M"), 5)

	fmt.Printf("%+v\n", left)

	// Output:
	// [NODE 4:4]
	// - k: 41
	//   v: `41` => `41`
	// - k: 42
	//   v: `42` => `42`
	// - k: 58
	//   v: `58` => `58`
	// - k: 59
	//   v: `59` => `59`
}

func ExampleBtree_merge_inner() {
	left := inner_node(
		leaf_node("A"),
		leaf_node("B"),
	)
	right := inner_node(
		leaf_node("X"),
		leaf_node("Y"),
	)

	merge(left, right, []byte("M"), 5)

	fmt.Printf("%+v\n", left)

	// Output:
	// [NODE 4:4]
	// - k: [BEFORE]
	//   v: [NODE 1:1]
	//   - k: 41
	//     v: `41` => `41`
	// - k: 42
	//   v: [NODE 1:1]
	//   - k: 42
	//     v: `42` => `42`
	// - k: 4d
	//   v: [NODE 1:1]
	//   - k: 58
	//     v: `58` => `58`
	// - k: 59
	//   v: [NODE 1:1]
	//   - k: 59
	//     v: `59` => `59`
}
