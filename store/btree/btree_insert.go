package btree

import (
	"simplex.sh/store/cas"
)

func (n *node_t) insert_ref(collated_key []byte, next_ref *ref_t, order int, store cas.GetterSetter) (prev *ref_t, new_root *node_t, err error) {
	prev_ref, right_key, right_ref, err := n.insert_ref_inner(collated_key, next_ref, order, store)
	if err != nil {
		return nil, nil, err
	}

	prev = prev_ref
	new_root = n

	// do we need a new root
	if right_key != nil && right_ref != nil {

		var (
			left_node  = n
			right_node = right_ref.cache.(*node_t)
		)

		new_root = &node_t{
			Type:         root_node_type,
			CollatedKeys: make([][]byte, 1, order),
			Children:     make([]*ref_t, 2, order+1),
			changed:      true,
		}

		left_ref := &ref_t{
			Flags: elt_is_set | ref_is_nod,
			Len:   left_node.Len(),
			cache: left_node,
		}

		left_node.ref = left_ref

		if left_node.Type&leaf_node_type > 0 {
			left_node.Type = leaf_node_type
		} else {
			left_node.Type = inner_node_type
		}

		if right_node.Type&leaf_node_type > 0 {
			right_node.Type = leaf_node_type
		} else {
			right_node.Type = inner_node_type
		}

		new_root.CollatedKeys[0] = right_key
		new_root.Children[0] = left_ref
		new_root.Children[1] = right_ref

	}

	return
}

func (n *node_t) insert_ref_inner(collated_key []byte, next_ref *ref_t, order int, store cas.GetterSetter) (prev *ref_t, right_key []byte, right_ref *ref_t, err error) {

	// for leaf nodes
	if n.Type&leaf_node_type > 0 {
		prev = n.insert_ref_into_leaf(collated_key, next_ref, order)
	}

	// for inner nodes
	if n.Type&leaf_node_type == 0 {
		prev, err = n.insert_ref_into_inner(collated_key, next_ref, order, store)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	right_key, right_ref = split(n, order)

	// recalculate n Len
	if n.ref != nil {
		n.ref.Len = n.Len()
	}

	return
}

func (n *node_t) insert_ref_into_inner(collated_key []byte, next_ref *ref_t, order int, store cas.GetterSetter) (prev *ref_t, err error) {
	_, _, ref := n.search_ref(collated_key)

	if ref == nil {
		panic("ref should never be nil for inner nodes")
	}

	child_node, err := ref.load_node(store)
	if err != nil {
		return nil, err
	}

	ch_prev, ch_right_key, ch_right_ref, err := child_node.insert_ref_inner(collated_key, next_ref, order, store)
	if err != nil {
		return nil, err
	}

	// mark as changed
	n.changed = true

	// pass prev ref
	prev = ch_prev

	// a split occured we must insert a new ref
	if ch_right_key != nil && ch_right_ref != nil {
		key_idx, ref_idx, _ := n.search_ref(ch_right_key)

		n.insert_collated_key(key_idx+1, ch_right_key, order)
		n.insert_child_ref(ref_idx+1, ch_right_ref, order)
	}

	return prev, nil
}

func (n *node_t) insert_ref_into_leaf(collated_key []byte, next_ref *ref_t, order int) *ref_t {
	var (
		prev *ref_t
	)

	key_idx, ref_idx, ref := n.search_ref(collated_key)

	// mark as changed
	n.changed = true

	// found exact match
	if ref != nil {
		prev = ref
		n.Children[ref_idx] = next_ref
	}

	// must insert
	if ref == nil {
		n.insert_collated_key(key_idx, collated_key, order)
		n.insert_child_ref(ref_idx, next_ref, order)
	}

	return prev
}

func (n *node_t) insert_collated_key(key_idx int, collated_key []byte, order int) {
	collated_keys := resize_collated_keys(n, n.CollatedKeys, order)

	if l := len(collated_keys); key_idx >= l {
		collated_keys = collated_keys[:l+1]
		collated_keys[l] = collated_key

	} else {

		// move other refs
		collated_keys = collated_keys[:l+1]
		copy(collated_keys[key_idx+1:], collated_keys[key_idx:])

		// set new key
		collated_keys[key_idx] = collated_key

	}

	n.CollatedKeys = collated_keys
}

func (n *node_t) insert_child_ref(ref_idx int, ref *ref_t, order int) {
	children := resize_children(n, n.Children, order)

	if l := len(children); ref_idx >= l {
		children = children[:l+1]
		children[l] = ref

	} else {

		// move other refs
		children = children[:l+1]
		copy(children[ref_idx+1:], children[ref_idx:])

		// set new ref
		children[ref_idx] = ref

	}

	n.Children = children
}

func resize_collated_keys(n *node_t, c [][]byte, order int) [][]byte {
	capacity := order

	if cap(c) < capacity {
		d := make([][]byte, len(c), capacity)
		copy(d, c)
		return d
	}

	return c
}

func resize_children(n *node_t, c []*ref_t, order int) []*ref_t {
	capacity := order + 1

	if n.Type&leaf_node_type > 0 {
		capacity -= 1
	}

	if cap(c) < capacity {
		d := make([]*ref_t, len(c), capacity)
		copy(d, c)
		return d
	}

	return c
}
