package btree

func (n *node_t) insert_ref(collated_key []byte, next_ref *ref_t) (prev *ref_t, new_root *node_t, err error) {
	prev_ref, split_key, right_ref, err := n.insert_ref_inner(collated_key, next_ref)
	if err != nil {
		return nil, nil, err
	}

	prev = prev_ref
	new_root = n

	// do we need a new root
	if split_key != nil && right_ref != nil {

		var (
			left_node  = n
			right_node = right_ref.cache.(*node_t)
		)

		new_root = &node_t{
			Type:         root_node_type,
			CollatedKeys: make([][]byte, 1, B),
			Children:     make([]*ref_t, 2, B+1),
			store:        left_node.store,
			parent:       nil,
			changed:      true,
		}

		left_ref := &ref_t{
			Flags: elt_is_set | ref_is_nod,
			Len:   left_node.Len(),
			cache: left_node,
		}

		left_node.ref = left_ref
		left_node.parent = new_root

		right_node.parent = new_root

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

		new_root.CollatedKeys[0] = split_key
		new_root.Children[0] = left_ref
		new_root.Children[1] = right_ref

	}

	return
}

func (n *node_t) insert_ref_inner(collated_key []byte, next_ref *ref_t) (prev *ref_t, split_key []byte, right_ref *ref_t, err error) {

	key_idx, ref_idx, ref := n.search_ref(collated_key)

	// for leaf nodes
	if n.Type&leaf_node_type > 0 {

		// mark as changed
		n.changed = true

		// found exact match
		if ref != nil {
			prev = ref
			n.Children[ref_idx] = next_ref
		}

		// must insert
		if ref == nil {
			n.insert_collated_key(key_idx, collated_key)
			n.insert_child_ref(ref_idx, next_ref)
		}

	}

	// for inner nodes
	if n.Type&leaf_node_type == 0 {

		// found bucket
		if ref != nil {
			child_node, err := ref.load_node(n.store, n)
			if err != nil {
				return nil, nil, nil, err
			}

			ch_prev, ch_split_key, ch_right_ref, err := child_node.insert_ref_inner(collated_key, next_ref)
			if err != nil {
				return nil, nil, nil, err
			}

			// mark as changed
			n.changed = true

			// pass prev ref
			prev = ch_prev

			// a split occured we must insert a new ref
			if ch_split_key != nil && ch_right_ref != nil {
				key_idx, ref_idx, _ := n.search_ref(ch_split_key)

				n.insert_collated_key(key_idx+1, ch_split_key)
				n.insert_child_ref(ref_idx+1, ch_right_ref)
			}

		}

		// fatal: no bucket found
		if ref == nil {
			panic("ref should never be nil for inner nodes")
		}

	}

	// do we need to split this node
	if n.has_too_many_children() {

		// find the split idx :: floor(len(keys) / 2)
		split_idx := len(n.CollatedKeys) / 2

		right_node := &node_t{
			Type:         n.Type,
			CollatedKeys: make([][]byte, 0, B),
			Children:     make([]*ref_t, 0, B+1),
			store:        n.store,
			parent:       n.parent,
			changed:      true,
		}

		{ // move keys
			slice := n.CollatedKeys[split_idx:]
			right_node.CollatedKeys = right_node.CollatedKeys[:len(slice)]
			copy(right_node.CollatedKeys, slice)
			n.CollatedKeys = n.CollatedKeys[:split_idx]
		}

		{ // move refs
			if n.Type&leaf_node_type == 0 {
				split_idx += 1
			}
			slice := n.Children[split_idx:]
			right_node.Children = right_node.Children[:len(slice)]
			copy(right_node.Children, slice)
			n.Children = n.Children[:split_idx]

			for _, ref := range right_node.Children {
				if ref == nil {
					continue
				}
				if node, ok := ref.cache.(*node_t); ok && node != nil {
					node.parent = right_node
				}
			}
		}

		right_ref = &ref_t{
			Flags: elt_is_set | ref_is_nod,
			Len:   right_node.Len(),
			cache: right_node,
		}
		right_node.ref = right_ref

		split_key = right_node.CollatedKeys[0]

		if n.Type&leaf_node_type == 0 {
			right_node.CollatedKeys = right_node.CollatedKeys[1:]
		}
	}

	// recalculate n Len
	if n.ref != nil {
		n.ref.Len = n.Len()
	}

	return
}

func (n *node_t) insert_collated_key(key_idx int, collated_key []byte) {
	collated_keys := n.CollatedKeys

	// make sure we have enough space (B - 1) + 1 (1 extra for insert/split)
	capacity := B
	if cap(collated_keys) < capacity {
		collated_keys = make([][]byte, len(n.CollatedKeys), capacity)
		copy(collated_keys, n.CollatedKeys)
	}

	if l := len(collated_keys); key_idx >= l {
		collated_keys = collated_keys[:l+1]
		collated_keys[l] = collated_key

	} else {

		// move other refs
		copy(collated_keys[key_idx+1:], collated_keys[key_idx:])

		// set new key
		collated_keys[key_idx] = collated_key

	}

	n.CollatedKeys = collated_keys
}

func (n *node_t) insert_child_ref(ref_idx int, ref *ref_t) {
	children := n.Children

	// make sure we have enough space (B - 1) + 1 (1 extra for insert/split)
	capacity := B + 1
	if n.Type&leaf_node_type > 0 {
		capacity -= 1
	}
	if cap(children) < capacity {
		children = make([]*ref_t, len(n.Children), capacity)
		copy(children, n.Children)
	}

	if l := len(children); ref_idx >= l {
		children = children[:l+1]
		children[l] = ref

	} else {

		// move other refs
		copy(children[ref_idx+1:], children[ref_idx:])

		// set new ref
		children[ref_idx] = ref

	}

	n.Children = children
}
