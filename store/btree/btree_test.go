package btree

import (
	"bytes"
	//"fmt"
	"simplex.sh/store/cas"
	"simplex.sh/store/memory"
	"testing"
)

func TestBTree(t *testing.T) {
	s := memory.New()

	hello, err := cas.Encode(s, "Hello", 0)
	if err != nil {
		t.Fatal("expected err to be nil")
	}

	world, err := cas.Encode(s, "world", 0)
	if err != nil {
		t.Fatal("expected err to be nil")
	}

	moon, err := cas.Encode(s, "moon", 0)
	if err != nil {
		t.Fatal("expected err to be nil")
	}

	tree := New(s)

	// test Set (new)
	prev, err := tree.Set(cas.Collate("Hello"), hello, world)
	if prev != nil {
		t.Fatal("expected prev to be the zero cas.Addr")
	}
	if err != nil {
		t.Fatal("expected err to be nil")
	}

	key, elt, err := tree.Get(cas.Collate("Hello"))
	if bytes.Compare(key, hello) != 0 {
		t.Fatalf("expected key to be the addr of hello (`%v`) but was (`%+v`)", hello, key)
	}
	if bytes.Compare(elt, world) != 0 {
		t.Fatalf("expected elt to be the addr of world (`%v`) but was (`%+v`)", world, elt)
	}
	if err != nil {
		t.Fatal("expected err to be nil")
	}

	// test Set (update)
	prev, err = tree.Set(cas.Collate("Hello"), hello, moon)
	if bytes.Compare(prev, world) != 0 {
		t.Fatalf("expected prev to be the addr of world (`%v`) but was (`%+v`)", world, prev)
	}
	if err != nil {
		t.Fatal("expected err to be nil")
	}

	key, elt, err = tree.Get(cas.Collate("Hello"))
	if bytes.Compare(key, hello) != 0 {
		t.Fatalf("expected key to be the addr of hello (`%v`) but was (`%+v`)", hello, key)
	}
	if bytes.Compare(elt, moon) != 0 {
		t.Fatalf("expected elt to be the addr of moon (`%v`) but was (`%+v`)", moon, elt)
	}
	if err != nil {
		t.Fatal("expected err to be nil")
	}

	// test Del
	prev_key, prev_elt, err := tree.Del(cas.Collate("Hello"))
	if bytes.Compare(prev_key, hello) != 0 {
		t.Fatalf("expected prev to be the addr of hello (`%v`) but was (`%+v`)", hello, prev_key)
	}
	if bytes.Compare(prev_elt, moon) != 0 {
		t.Fatalf("expected prev to be the addr of moon (`%v`) but was (`%+v`)", moon, prev_elt)
	}
	if err != nil {
		t.Fatal("expected err to be nil")
	}

	key, elt, err = tree.Get(cas.Collate("Hello"))
	if bytes.Compare(key, nil) != 0 {
		t.Fatalf("expected key to be the addr of hello (`%v`) but was (`%+v`)", hello, key)
	}
	if bytes.Compare(elt, nil) != 0 {
		t.Fatalf("expected elt to be the addr of moon (`%v`) but was (`%+v`)", moon, elt)
	}
	if err != nil {
		t.Fatal("expected err to be nil")
	}
}

func TestBTreeLarge(t *testing.T) {
	s := memory.New()

	foo, err := cas.Encode(s, "foo", 0)
	if err != nil {
		t.Fatal("expected err to be nil")
	}

	tree := New(s)

	C := 3000

	//defer func() { fmt.Printf("T: %+v\n", tree) }()

	for i := 0; i < C; i++ {
		bar_x, err := cas.Encode(s, i, 0)
		if err != nil {
			t.Fatal("expected err to be nil")
		}

		_, err = tree.Set(cas.Collate(i), bar_x, foo)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < C; i++ {
		bar_x, err := cas.Encode(s, i, 0)
		if err != nil {
			t.Fatal("expected err to be nil")
		}

		key, elt, err := tree.Get(cas.Collate(i))
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(key, bar_x) != 0 {
			t.Fatalf("expected key to be the addr of bar_x (`%v`) but was (`%+v`)", bar_x, key)
		}
		if bytes.Compare(elt, foo) != 0 {
			t.Fatalf("expected elt to be the addr of foo (`%v`) but was (`%+v`)", foo, elt)
		}
	}

	if tree.Len != uint64(C) {
		t.Fatalf("expected tree.Len() to be %d but was %d", C, tree.Len)
	}
}

func TestBTree_search_ref(t *testing.T) {
	var (
		n       *node_t
		key_idx int
		ref_idx int
		ref     *ref_t
	)

	{ // root node (not a leaf)
		n = &node_t{
			Type: root_node_type,
			CollatedKeys: [][]byte{
				[]byte("3"),
				[]byte("5"),
			},
			Children: []*ref_t{
				{},
				{},
				{},
			},
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("1"))
		if key_idx != -1 {
			t.Fatal("expected key_idx to be -1 (%d)", key_idx)
		}
		if ref_idx != 0 {
			t.Fatal("expected ref_idx to be 0 (%d)", ref_idx)
		}
		if ref != n.Children[0] {
			t.Fatal("expected ref to be n.Children[0]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("2"))
		if key_idx != -1 {
			t.Fatal("expected key_idx to be -1 (%d)", key_idx)
		}
		if ref_idx != 0 {
			t.Fatal("expected ref_idx to be 0 (%d)", ref_idx)
		}
		if ref != n.Children[0] {
			t.Fatal("expected ref to be n.Children[0]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("3"))
		if key_idx != 0 {
			t.Fatal("expected key_idx to be 0 (%d)", key_idx)
		}
		if ref_idx != 1 {
			t.Fatal("expected ref_idx to be 1 (%d)", ref_idx)
		}
		if ref != n.Children[1] {
			t.Fatal("expected ref to be n.Children[1]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("4"))
		if key_idx != 0 {
			t.Fatal("expected key_idx to be 0 (%d)", key_idx)
		}
		if ref_idx != 1 {
			t.Fatal("expected ref_idx to be 1 (%d)", ref_idx)
		}
		if ref != n.Children[1] {
			t.Fatal("expected ref to be n.Children[1]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("5"))
		if key_idx != 1 {
			t.Fatal("expected key_idx to be 1 (%d)", key_idx)
		}
		if ref_idx != 2 {
			t.Fatal("expected ref_idx to be 2 (%d)", ref_idx)
		}
		if ref != n.Children[2] {
			t.Fatal("expected ref to be n.Children[2]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("6"))
		if key_idx != 1 {
			t.Fatal("expected key_idx to be 1 (%d)", key_idx)
		}
		if ref_idx != 2 {
			t.Fatal("expected ref_idx to be 2 (%d)", ref_idx)
		}
		if ref != n.Children[2] {
			t.Fatal("expected ref to be n.Children[2]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("51"))
		if key_idx != 1 {
			t.Fatal("expected key_idx to be 1 (%d)", key_idx)
		}
		if ref_idx != 2 {
			t.Fatal("expected ref_idx to be 2 (%d)", ref_idx)
		}
		if ref != n.Children[2] {
			t.Fatal("expected ref to be n.Children[2]")
		}
	}

	{ // root node (is a leaf)
		n = &node_t{
			Type: root_node_type | leaf_node_type,
			CollatedKeys: [][]byte{
				[]byte("3"),
				[]byte("5"),
			},
			Children: []*ref_t{
				{},
				{},
			},
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("1"))
		if key_idx != 0 {
			t.Fatal("expected key_idx to be 0 (%d)", key_idx)
		}
		if ref_idx != 0 {
			t.Fatal("expected ref_idx to be 0 (%d)", ref_idx)
		}
		if ref != nil {
			t.Fatal("expected ref to be nil")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("2"))
		if key_idx != 0 {
			t.Fatal("expected key_idx to be 0 (%d)", key_idx)
		}
		if ref_idx != 0 {
			t.Fatal("expected ref_idx to be 0 (%d)", ref_idx)
		}
		if ref != nil {
			t.Fatal("expected ref to be nil")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("3"))
		if key_idx != 0 {
			t.Fatal("expected key_idx to be 0 (%d)", key_idx)
		}
		if ref_idx != 0 {
			t.Fatal("expected ref_idx to be 0 (%d)", ref_idx)
		}
		if ref != n.Children[0] {
			t.Fatal("expected ref to be n.Children[0]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("4"))
		if key_idx != 1 {
			t.Fatal("expected key_idx to be 1 (%d)", key_idx)
		}
		if ref_idx != 1 {
			t.Fatal("expected ref_idx to be 1 (%d)", ref_idx)
		}
		if ref != nil {
			t.Fatal("expected ref to be nil")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("5"))
		if key_idx != 1 {
			t.Fatal("expected key_idx to be 1 (%d)", key_idx)
		}
		if ref_idx != 1 {
			t.Fatal("expected ref_idx to be 1 (%d)", ref_idx)
		}
		if ref != n.Children[1] {
			t.Fatal("expected ref to be n.Children[1]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("6"))
		if key_idx != 2 {
			t.Fatal("expected key_idx to be 2 (%d)", key_idx)
		}
		if ref_idx != 2 {
			t.Fatal("expected ref_idx to be 2 (%d)", ref_idx)
		}
		if ref != nil {
			t.Fatal("expected ref to be nil")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("51"))
		if key_idx != 2 {
			t.Fatal("expected key_idx to be 2 (%d)", key_idx)
		}
		if ref_idx != 2 {
			t.Fatal("expected ref_idx to be 2 (%d)", ref_idx)
		}
		if ref != nil {
			t.Fatal("expected ref to be nil")
		}
	}

	{ // leaf node
		n = &node_t{
			Type: leaf_node_type,
			CollatedKeys: [][]byte{
				[]byte("3"),
				[]byte("5"),
			},
			Children: []*ref_t{
				{},
				{},
			},
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("1"))
		if key_idx != 0 {
			t.Fatal("expected key_idx to be 0 (%d)", key_idx)
		}
		if ref_idx != 0 {
			t.Fatal("expected ref_idx to be 0 (%d)", ref_idx)
		}
		if ref != nil {
			t.Fatal("expected ref to be nil")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("2"))
		if key_idx != 0 {
			t.Fatal("expected key_idx to be 0 (%d)", key_idx)
		}
		if ref_idx != 0 {
			t.Fatal("expected ref_idx to be 0 (%d)", ref_idx)
		}
		if ref != nil {
			t.Fatal("expected ref to be nil")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("3"))
		if key_idx != 0 {
			t.Fatal("expected key_idx to be 0 (%d)", key_idx)
		}
		if ref_idx != 0 {
			t.Fatal("expected ref_idx to be 0 (%d)", ref_idx)
		}
		if ref != n.Children[0] {
			t.Fatal("expected ref to be n.Children[0]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("4"))
		if key_idx != 1 {
			t.Fatal("expected key_idx to be 1 (%d)", key_idx)
		}
		if ref_idx != 1 {
			t.Fatal("expected ref_idx to be 1 (%d)", ref_idx)
		}
		if ref != nil {
			t.Fatal("expected ref to be nil")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("5"))
		if key_idx != 1 {
			t.Fatal("expected key_idx to be 1 (%d)", key_idx)
		}
		if ref_idx != 1 {
			t.Fatal("expected ref_idx to be 1 (%d)", ref_idx)
		}
		if ref != n.Children[1] {
			t.Fatal("expected ref to be n.Children[1]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("6"))
		if key_idx != 2 {
			t.Fatal("expected key_idx to be 2 (%d)", key_idx)
		}
		if ref_idx != 2 {
			t.Fatal("expected ref_idx to be 2 (%d)", ref_idx)
		}
		if ref != nil {
			t.Fatal("expected ref to be nil")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("51"))
		if key_idx != 2 {
			t.Fatal("expected key_idx to be 2 (%d)", key_idx)
		}
		if ref_idx != 2 {
			t.Fatal("expected ref_idx to be 2 (%d)", ref_idx)
		}
		if ref != nil {
			t.Fatal("expected ref to be nil")
		}
	}

	{ // inner node
		n = &node_t{
			Type: inner_node_type,
			CollatedKeys: [][]byte{
				[]byte("3"),
				[]byte("5"),
			},
			Children: []*ref_t{
				{},
				{},
				{},
			},
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("1"))
		if key_idx != -1 {
			t.Fatal("expected key_idx to be -1 (%d)", key_idx)
		}
		if ref_idx != 0 {
			t.Fatal("expected ref_idx to be 0 (%d)", ref_idx)
		}
		if ref != n.Children[0] {
			t.Fatal("expected ref to be n.Children[0]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("2"))
		if key_idx != -1 {
			t.Fatal("expected key_idx to be -1 (%d)", key_idx)
		}
		if ref_idx != 0 {
			t.Fatal("expected ref_idx to be 0 (%d)", ref_idx)
		}
		if ref != n.Children[0] {
			t.Fatal("expected ref to be n.Children[0]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("3"))
		if key_idx != 0 {
			t.Fatal("expected key_idx to be 0 (%d)", key_idx)
		}
		if ref_idx != 1 {
			t.Fatal("expected ref_idx to be 1 (%d)", ref_idx)
		}
		if ref != n.Children[1] {
			t.Fatal("expected ref to be n.Children[1]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("4"))
		if key_idx != 0 {
			t.Fatal("expected key_idx to be 0 (%d)", key_idx)
		}
		if ref_idx != 1 {
			t.Fatal("expected ref_idx to be 1 (%d)", ref_idx)
		}
		if ref != n.Children[1] {
			t.Fatal("expected ref to be n.Children[1]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("5"))
		if key_idx != 1 {
			t.Fatal("expected key_idx to be 1 (%d)", key_idx)
		}
		if ref_idx != 2 {
			t.Fatal("expected ref_idx to be 2 (%d)", ref_idx)
		}
		if ref != n.Children[2] {
			t.Fatal("expected ref to be n.Children[2]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("6"))
		if key_idx != 1 {
			t.Fatal("expected key_idx to be 1 (%d)", key_idx)
		}
		if ref_idx != 2 {
			t.Fatal("expected ref_idx to be 2 (%d)", ref_idx)
		}
		if ref != n.Children[2] {
			t.Fatal("expected ref to be n.Children[2]")
		}

		key_idx, ref_idx, ref = n.search_ref([]byte("51"))
		if key_idx != 1 {
			t.Fatal("expected key_idx to be 1 (%d)", key_idx)
		}
		if ref_idx != 2 {
			t.Fatal("expected ref_idx to be 2 (%d)", ref_idx)
		}
		if ref != n.Children[2] {
			t.Fatal("expected ref to be n.Children[2]")
		}
	}
}
