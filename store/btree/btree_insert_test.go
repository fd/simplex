package btree

import (
	"fmt"
)

func ExampleBtree_insert_leaf_first() {
	n := leaf_node()

	n.insert_ref([]byte("A"), ref("A"), 5, nil)

	fmt.Printf("%+v\n", n)

	// Output:
	// [NODE 1:1]
	// - k: 41
	//   v: `41` => `41`
}

func ExampleBtree_insert_leaf_next() {
	n := leaf_node("A")

	n.insert_ref([]byte("B"), ref("B"), 5, nil)

	fmt.Printf("%+v\n", n)

	// Output:
	// [NODE 2:2]
	// - k: 41
	//   v: `41` => `41`
	// - k: 42
	//   v: `42` => `42`
}

func ExampleBtree_insert_leaf_last() {
	n := leaf_node("A", "B", "C")

	n.insert_ref([]byte("D"), ref("D"), 5, nil)

	fmt.Printf("%+v\n", n)

	// Output:
	// [NODE 4:4]
	// - k: 41
	//   v: `41` => `41`
	// - k: 42
	//   v: `42` => `42`
	// - k: 43
	//   v: `43` => `43`
	// - k: 44
	//   v: `44` => `44`
}

func ExampleBtree_insert_leaf_at_first_idx() {
	n := leaf_node("B", "C", "D")

	n.insert_ref([]byte("A"), ref("A"), 5, nil)

	fmt.Printf("%+v\n", n)

	// Output:
	// [NODE 4:4]
	// - k: 41
	//   v: `41` => `41`
	// - k: 42
	//   v: `42` => `42`
	// - k: 43
	//   v: `43` => `43`
	// - k: 44
	//   v: `44` => `44`
}

func ExampleBtree_insert_leaf_at_last_idx() {
	n := leaf_node("A", "B", "D")

	n.insert_ref([]byte("C"), ref("C"), 5, nil)

	fmt.Printf("%+v\n", n)

	// Output:
	// [NODE 4:4]
	// - k: 41
	//   v: `41` => `41`
	// - k: 42
	//   v: `42` => `42`
	// - k: 43
	//   v: `43` => `43`
	// - k: 44
	//   v: `44` => `44`
}

func ExampleBtree_insert_leaf_middle() {
	n := leaf_node("A", "B", "C")

	n.insert_ref([]byte("D"), ref("D"), 5, nil)

	fmt.Printf("%+v\n", n)

	// Output:
	// [NODE 4:4]
	// - k: 41
	//   v: `41` => `41`
	// - k: 42
	//   v: `42` => `42`
	// - k: 43
	//   v: `43` => `43`
	// - k: 44
	//   v: `44` => `44`
}

func ExampleBtree_insert_inner_no_split() {
	n := inner_node(
		leaf_node("A", "B"),
		leaf_node("C", "D"),
		leaf_node("E", "F"),
	)

	n.insert_ref([]byte("G"), ref("G"), 5, nil)

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
	//   - k: 44
	//     v: `44` => `44`
	// - k: 45
	//   v: [NODE 3:3]
	//   - k: 45
	//     v: `45` => `45`
	//   - k: 46
	//     v: `46` => `46`
	//   - k: 47
	//     v: `47` => `47`
}

func ExampleBtree_insert_inner_last() {
	n := inner_node(
		leaf_node("A", "B"),
		leaf_node("C", "D", "E", "F"),
	)

	n.insert_ref([]byte("G"), ref("G"), 5, nil)

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
	//   - k: 44
	//     v: `44` => `44`
	// - k: 45
	//   v: [NODE 3:3]
	//   - k: 45
	//     v: `45` => `45`
	//   - k: 46
	//     v: `46` => `46`
	//   - k: 47
	//     v: `47` => `47`
}
