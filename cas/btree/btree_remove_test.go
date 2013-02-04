package btree

import (
	"fmt"
)

func ExampleBtree_remove_leaf_first() {
	n := leaf_node("A", "B", "C", "D")

	n.remove_ref([]byte("A"), 5)

	fmt.Printf("%+v\n", n)

	// Output:
	// [NODE 3:3]
	// - k: 42
	//   v: `42` => `42`
	// - k: 43
	//   v: `43` => `43`
	// - k: 44
	//   v: `44` => `44`
}

func ExampleBtree_remove_leaf_last() {
	n := leaf_node("A", "B", "C", "D")

	n.remove_ref([]byte("D"), 5)

	fmt.Printf("%+v\n", n)

	// Output:
	// [NODE 3:3]
	// - k: 41
	//   v: `41` => `41`
	// - k: 42
	//   v: `42` => `42`
	// - k: 43
	//   v: `43` => `43`
}

func ExampleBtree_remove_leaf_middle() {
	n := leaf_node("A", "B", "C", "D")

	n.remove_ref([]byte("B"), 5)

	fmt.Printf("%+v\n", n)

	// Output:
	// [NODE 3:3]
	// - k: 41
	//   v: `41` => `41`
	// - k: 43
	//   v: `43` => `43`
	// - k: 44
	//   v: `44` => `44`
}

func ExampleBtree_remove_inner_borrow_left() {
	n := inner_node(
		leaf_node("A", "B", "C"),
		leaf_node("D", "E"),
		leaf_node("F", "G", "H"),
	)

	n.remove_ref([]byte("D"), 5)

	fmt.Printf("%+v\n", n)

	// Output:
	// [NODE 3:7]
	// - k: [BEFORE]
	//   v: [NODE 2:2]
	//   - k: 41
	//     v: `41` => `41`
	//   - k: 42
	//     v: `42` => `42`
	// - k: 43
	//   v: [NODE 2:2]
	//   - k: 43
	//     v: `43` => `43`
	//   - k: 45
	//     v: `45` => `45`
	// - k: 46
	//   v: [NODE 3:3]
	//   - k: 46
	//     v: `46` => `46`
	//   - k: 47
	//     v: `47` => `47`
	//   - k: 48
	//     v: `48` => `48`
}

func ExampleBtree_remove_inner_borrow_right() {
	n := inner_node(
		leaf_node("D", "E"),
		leaf_node("F", "G", "H"),
	)

	n.remove_ref([]byte("D"), 5)

	fmt.Printf("%+v\n", n)

	// Output:
	// [NODE 2:4]
	// - k: [BEFORE]
	//   v: [NODE 2:2]
	//   - k: 45
	//     v: `45` => `45`
	//   - k: 46
	//     v: `46` => `46`
	// - k: 47
	//   v: [NODE 2:2]
	//   - k: 47
	//     v: `47` => `47`
	//   - k: 48
	//     v: `48` => `48`
}

func ExampleBtree_remove_inner_merge_left() {
	n := inner_node(
		leaf_node("D", "E"),
		leaf_node("F", "G"),
	)

	n.remove_ref([]byte("F"), 5)

	fmt.Printf("%+v\n", n)

	// Output:
	// [NODE 1:3]
	// - k: [BEFORE]
	//   v: [NODE 3:3]
	//   - k: 44
	//     v: `44` => `44`
	//   - k: 45
	//     v: `45` => `45`
	//   - k: 47
	//     v: `47` => `47`
}

func ExampleBtree_remove_inner_merge_right() {
	n := inner_node(
		leaf_node("D", "E"),
		leaf_node("F", "G"),
	)

	n.remove_ref([]byte("D"), 5)

	fmt.Printf("%+v\n", n)

	// Output:
	// [NODE 1:3]
	// - k: [BEFORE]
	//   v: [NODE 3:3]
	//   - k: 45
	//     v: `45` => `45`
	//   - k: 46
	//     v: `46` => `46`
	//   - k: 47
	//     v: `47` => `47`
}
