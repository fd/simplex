package btree

import (
	"bytes"
	"github.com/fd/simplex/data/storage"
	"sort"
)

const B = 512

// 512 * ((2 * 20) + 1 + 4 + 20)

const (
	root_node_type node_type_t = 1 << iota
	inner_node_type
	leaf_node_type
)

const (
	key_is_set ref_flags_t = 1 << iota
	key_is_sha
	elt_is_set
	elt_is_sha
	ref_is_val
	ref_is_nod
)

type (
	node_type_t byte
	ref_flags_t byte

	tree_t struct {
		root *node_t
	}

	node_t struct {
		Type node_type_t

		CollatedKeys [][]byte
		Children     []*ref_t

		parent *node_t
		ref    *ref_t
		store  *storage.S
	}

	ref_t struct {
		Flags ref_flags_t
		Len   uint64
		Key   []byte
		Elt   []byte

		cache interface{}
	}
)

func (t *tree_t) insert(collated_key []byte, key, elt interface{}) {
	ref, err := build_ref(key, elt)
	if err != nil {
		return err
	}

	new_root = t.root.search(collated_key).insert_ref(collated_key, ref)
	if new_root != nil {
		t.root = new_root
	}
}

func (n *node_t) Len() uint64 {
	var l uint64
	for _, ref := range n.Children {
		l += ref.Len
	}
	return l
}

func (n *node_t) has_too_many_children() bool {
	if n.Type&leaf_node_type > 0 {
		return len(n.Children) > (B - 1)
	}
	return len(n.Children) > B
}

func (n *node_t) has_too_few_children() bool {
	if n.Type&root_node_type > 0 && n.Type&leaf_node_type > 0 {
		return len(n.Children) < 1
	}

	if n.Type&root_node_type > 0 {
		return len(n.Children) < 2
	}

	if n.Type&inner_node_type > 0 {
		return len(n.Children) < ((B / 2) + 1)
	}

	if n.Type&leaf_node_type > 0 {
		return len(n.Children) < (B / 2)
	}

	panic("not reached")
}

func (n *node_t) insert_ref(collated_key []byte, ref *ref_t) (new_root *node_t) {
	// add record
	{
		// find the idx to insert ref at
		idx := sort.Search(len(n.CollatedKeys), func(i int) bool {
			return bytes.Compare(n.CollatedKeys[i], collated_key) > 0
		})

		{
			collated_keys := n.CollatedKeys
			if cap(collated_keys) < (B) {
				collated_keys = make([][]byte, len(n.CollatedKeys), B)
				copy(collated_keys, n.CollatedKeys)
			}
			copy(collated_keys[idx+1:], collated_keys[idx:])
			collated_keys[idx] = collated_key
			n.CollatedKeys = collated_keys
		}

		{
			children := n.Children
			if cap(children) < (B + 1) {
				children = make([]*ref_t, len(n.Children), B+1)
				copy(children, n.Children)
			}
			copy(children[idx+1:], children[idx:])
			children[idx] = ref
			n.Children = children
		}
	}

	if n.has_too_many_children() {
		// split
		if n.Type&root_node_type > 0 {
			if n.Type&leaf_node_type > 0 {
				n.Type = leaf_node_type
			} else {
				n.Type = inner_node_type
			}

			new_root = &node_t{
				Type:         root_node_type,
				CollatedKeys: make([][]byte, 0, B),
				Children:     make([]*ref_t, 0, B+1),
				store:        n.store,
			}

			left_ref := &ref_t{
				Flags: elt_is_set | elt_is_sha,
				Len:   n.Len(),
				cache: n,
			}
			n.ref = left_ref
			n.parent = new_root

			new_root.Children[0] = left_ref
		}

		split_idx := len(n.CollatedKeys) / 2
		right := &node_t{
			Type:         n.Type,
			CollatedKeys: make([][]byte, 0, B),
			Children:     make([]*ref_t, 0, B+1),
			store:        n.store,
			parent:       n.parent,
		}

		{
			slice := n.CollatedKeys[split_idx:]
			copy(right.CollatedKeys[:len(slice)], slice)
			n.CollatedKeys = n.CollatedKeys[:split_idx]
		}

		{
			slice := n.Children[split_idx+1:]
			copy(right.Children[:len(slice)], slice)
			n.Children = n.Children[:split_idx+1]

			for _, ref := range n.Children {
				if node, ok := ref.cache.(*node_t); ok && node != nil {
					node.parent = right
				}
			}
		}

		n.ref.Len = n.Len()

		right_ref := &ref_t{
			Flags: elt_is_set | elt_is_sha,
			Len:   right.Len(),
			cache: right,
		}
		right.ref = right_ref

		new_root = n.parent.insert_ref(right.CollatedKeys[0], right_ref)
	}

	return
}

func (n *node_t) search(key []byte) *node_t {
	var (
		idx int
		ref *ref_t
	)

	if n.Type&leaf_node_type > 0 {
		return n
	}

	if len(n.CollatedKeys) == 0 {
		return n
	}

	// find ref idx
	idx = sort.Search(len(n.CollatedKeys), func(i int) bool {
		return bytes.Compare(n.CollatedKeys[i], key) > 0
	})

	ref = n.Children[idx]

	return ref.load_node(n.store, n).search(key)
}

func (ref *ref_t) load_node(store *storage.S, parent *node_t) *node_t {
	if ref.Flags&elt_is_set == 0 {
		panic("corrupt btree ref")
	}

	if ref.Flags&elt_is_sha == 0 {
		panic("corrupt btree ref")
	}

	if node, ok := ref.cache.(*node_t); ok && node != nil {
		return node
	}

	var (
		sha  = storage.SHA{}
		node = &node_t{}
	)

	copy([]byte(sha[:]), ref.Elt)
	if !store.Get(sha, &node) {
		panic("corrupt btree ref")
	}

	node.parent = parent
	node.ref = ref
	node.store = store
	ref.cache = node
	return node
}

func build_ref(store *storage.S, key, elt interface{}) (*ref_t, err) {
	var (
		key_buf bytes.Buffer
		elt_buf bytes.Buffer
		key_sha cas.SHA
		elt_sha cas.SHA
		ref     = &ref_t{Len: 1, Flags: ref_is_val}
	)

	key_sha, err := cas.NewEncoder(&key_buf).Encode(key)
	if err != nil {
		return nil, err
	}

	elt_sha, err := cas.NewEncoder(&elt_buf).Encode(elt)
	if err != nil {
		return nil, err
	}

	// Set the key
	if key_buf.Len() > 40 {
		err := store.SetObject(key_sha, key_buf.Bytes())
		if err != nil {
			return nil, err
		}
		ref.Flags |= key_is_set | key_is_sha
		ref.Key = []byte(key_sha[:])

	} else {
		ref.Flags |= key_is_set
		ref.Key = key_buf.Bytes()

	}

	// Set the elt
	if elt_buf.Len() > 40 {
		err := store.SetObject(elt_sha, elt_buf.Bytes())
		if err != nil {
			return nil, err
		}
		ref.Flags |= elt_is_set | elt_is_sha
		ref.Elt = []byte(elt_sha[:])

	} else {
		ref.Flags |= elt_is_set
		ref.Elt = elt_buf.Bytes()

	}

	return ref, nil
}
