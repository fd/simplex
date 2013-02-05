package btree

import (
	"encoding/hex"
	"fmt"
)

func ExampleBtree_split_leaf() {
	n := leaf_node("A", "B", "C", "D", "E", "F")

	right_key, right_ref := split(n, 5)

	fmt.Printf("%+v\n", n)
	fmt.Printf("%+v\n", hex.EncodeToString(right_key))
	fmt.Printf("%+v\n", right_ref.cache.(*node_t))

	// Output:
	// [NODE 3:3]
	// - k: 41
	//   v: `41` => `41`
	// - k: 42
	//   v: `42` => `42`
	// - k: 43
	//   v: `43` => `43`
	// 44
	// [NODE 3:3]
	// - k: 44
	//   v: `44` => `44`
	// - k: 45
	//   v: `45` => `45`
	// - k: 46
	//   v: `46` => `46`
}

func ExampleBtree_split_inner() {
	n := inner_node(
		leaf_node("A"),
		leaf_node("B"),
		leaf_node("C"),
		leaf_node("D"),
		leaf_node("E"),
		leaf_node("F"),
	)

	right_key, right_ref := split(n, 5)

	fmt.Printf("%+v\n", n)
	fmt.Printf("%+v\n", hex.EncodeToString(right_key))
	fmt.Printf("%+v\n", right_ref.cache.(*node_t))

	// Output:
	// [NODE 3:3]
	// - k: [BEFORE]
	//   v: [NODE 1:1]
	//   - k: 41
	//     v: `41` => `41`
	// - k: 42
	//   v: [NODE 1:1]
	//   - k: 42
	//     v: `42` => `42`
	// - k: 43
	//   v: [NODE 1:1]
	//   - k: 43
	//     v: `43` => `43`
	// 44
	// [NODE 3:3]
	// - k: [BEFORE]
	//   v: [NODE 1:1]
	//   - k: 44
	//     v: `44` => `44`
	// - k: 45
	//   v: [NODE 1:1]
	//   - k: 45
	//     v: `45` => `45`
	// - k: 46
	//   v: [NODE 1:1]
	//   - k: 46
	//     v: `46` => `46`
}
