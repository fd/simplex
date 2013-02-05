package btree

import (
	"encoding/hex"
	"fmt"
	"github.com/fd/simplex/cas"
)

func ExampleBtree_borrow_leaf() {
	left := leaf_node("A", "B", "C", "D")
	right := leaf_node("X", "Y")

	placehoder, _ := borrow(right, left, true, []byte("M"), 5)

	fmt.Printf("%+v\n", left)
	fmt.Printf("%+v\n", hex.EncodeToString(placehoder))
	fmt.Printf("%+v\n", right)

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
	// - k: 58
	//   v: `58` => `58`
	// - k: 59
	//   v: `59` => `59`
}

func ExampleBtree_borrow_leaf_reverse() {
	left := leaf_node("A", "B")
	right := leaf_node("V", "W", "X", "Y")

	placehoder, _ := borrow(left, right, false, []byte("M"), 5)

	fmt.Printf("%+v\n", left)
	fmt.Printf("%+v\n", hex.EncodeToString(placehoder))
	fmt.Printf("%+v\n", right)

	// Output:
	// [NODE 3:3]
	// - k: 41
	//   v: `41` => `41`
	// - k: 42
	//   v: `42` => `42`
	// - k: 56
	//   v: `56` => `56`
	// 57
	// [NODE 3:3]
	// - k: 57
	//   v: `57` => `57`
	// - k: 58
	//   v: `58` => `58`
	// - k: 59
	//   v: `59` => `59`
}

func ExampleBtree_borrow_inner() {
	left := inner_node(
		leaf_node("A"),
		leaf_node("B"),
		leaf_node("C"),
		leaf_node("D"),
	)
	right := inner_node(
		leaf_node("X"),
		leaf_node("Y"),
	)

	placehoder, _ := borrow(right, left, true, []byte("M"), 5)

	fmt.Printf("%+v\n", left)
	fmt.Printf("%+v\n", hex.EncodeToString(placehoder))
	fmt.Printf("%+v\n", right)

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
	// - k: 4d
	//   v: [NODE 1:1]
	//   - k: 58
	//     v: `58` => `58`
	// - k: 59
	//   v: [NODE 1:1]
	//   - k: 59
	//     v: `59` => `59`
}

func ExampleBtree_borrow_inner_reverse() {
	left := inner_node(
		leaf_node("A"),
		leaf_node("B"),
	)
	right := inner_node(
		leaf_node("V"),
		leaf_node("W"),
		leaf_node("X"),
		leaf_node("Y"),
	)

	placehoder, _ := borrow(left, right, false, []byte("M"), 5)

	fmt.Printf("%+v\n", left)
	fmt.Printf("%+v\n", hex.EncodeToString(placehoder))
	fmt.Printf("%+v\n", right)

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
	// - k: 4d
	//   v: [NODE 1:1]
	//   - k: 56
	//     v: `56` => `56`
	// 57
	// [NODE 3:3]
	// - k: [BEFORE]
	//   v: [NODE 1:1]
	//   - k: 57
	//     v: `57` => `57`
	// - k: 58
	//   v: [NODE 1:1]
	//   - k: 58
	//     v: `58` => `58`
	// - k: 59
	//   v: [NODE 1:1]
	//   - k: 59
	//     v: `59` => `59`
}

func inner_node(nodes ...*node_t) *node_t {
	l := 0
	if len(nodes) >= 1 {
		l = len(nodes) - 1
	}

	keys := make([][]byte, l, B)
	refs := make([]*ref_t, len(nodes), B+1)

	for i, node := range nodes {
		if i > 0 {
			keys[i-1] = smallest_key(node)
		}
		refs[i] = inner_ref(node)
	}

	return &node_t{
		Type:         inner_node_type,
		CollatedKeys: keys,
		Children:     refs,
	}
}

func leaf_node(names ...string) *node_t {
	keys := make([][]byte, len(names), B)
	refs := make([]*ref_t, len(names), B)

	for i, name := range names {
		keys[i] = []byte(name)
		refs[i] = ref(name)
	}

	return &node_t{
		Type:         leaf_node_type,
		CollatedKeys: keys,
		Children:     refs,
	}
}

func ref(name string) *ref_t {
	return &ref_t{
		Flags: key_is_set | elt_is_set | ref_is_val,
		Len:   1,
		Key:   cas.Addr(name),
		Elt:   cas.Addr(name),
	}
}

func smallest_key(node *node_t) []byte {
	if node.Type&leaf_node_type > 0 {
		return node.CollatedKeys[0]
	}
	return smallest_key(node.Children[0].cache.(*node_t))
}

func inner_ref(node *node_t) *ref_t {
	r := &ref_t{
		Flags: elt_is_set | ref_is_nod,
		Len:   node.Len(),
		cache: node,
	}

	node.ref = r

	return r
}
