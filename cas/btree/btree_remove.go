package btree

func (n *node_t) remove_ref(collated_key []byte) (prev *ref_t, err error) {

	key_idx, ref_idx, middle_ref := n.search_ref(collated_key)

	if n.Type&leaf_node_type > 0 {

		// nothing is found
		if middle_ref == nil {
			return nil, nil
		}

		prev = middle_ref

		{ // delete key
			var (
				length        = len(n.CollatedKeys)
				collated_keys = n.CollatedKeys
				right         = collated_keys[key_idx+1:]
			)
			if len(right) > 0 {
				copy(collated_keys[key_idx:], right)
			}
			collated_keys = collated_keys[:length-1]
			n.CollatedKeys = collated_keys
		}

		{ // delete ref
			var (
				length   = len(n.Children)
				children = n.Children
				right    = children[ref_idx+1:]
			)
			if len(right) > 0 {
				copy(children[ref_idx:], right)
			}
			children = children[:length-1]
			n.Children = children
		}

		// update Len
		if n.ref != nil {
			n.ref.Len = n.Len()
		}

		// mark as changed
		n.changed = true

		return
	}

	// fatal: no bucket found
	if middle_ref == nil {
		panic("ref should never be nil for inner nodes")
	}

	// find the middle node
	middle_node, err := middle_ref.load_node(n.store, n)
	if err != nil {
		return nil, err
	}

	// propagate delete
	prev, err = middle_node.remove_ref(collated_key)
	if err != nil {
		return nil, err
	}
	if prev == nil {
		return nil, nil
	}

	// update ref key
	if middle_node.Type&leaf_node_type > 0 {
		n.CollatedKeys[key_idx] = middle_node.CollatedKeys[0]
	}

	// merge nodes
	if middle_node.has_too_few_children() {
		// find siblings
		var (
			left_ref   *ref_t
			right_ref  *ref_t
			left_node  *node_t
			right_node *node_t
		)

		if ref_idx > 0 {
			left_ref = n.Children[ref_idx-1]
			left_node, err = left_ref.load_node(n.store, n)
			if err != nil {
				return nil, err
			}

			// attempt to steal children
			if key, ok := borrow(middle_node, left_node, true, n.CollatedKeys[key_idx]); ok {
				n.CollatedKeys[key_idx] = key
				n.changed = true
				return prev, nil
			}
		}

		if (ref_idx + 1) < len(n.Children) {
			right_ref = n.Children[ref_idx+1]
			right_node, err = right_ref.load_node(n.store, n)
			if err != nil {
				return nil, err
			}

			// attempt to steal children
			if key, ok := borrow(middle_node, right_node, false, n.CollatedKeys[key_idx+1]); ok {
				n.CollatedKeys[key_idx+1] = key
				n.changed = true
				return prev, nil
			}
		}

		// attempt to merge siblings

		if merge(left_node, middle_node, n.CollatedKeys[key_idx]) {
			// remove ref at ref_idx
		}

		if merge(middle_node, right_node, n.CollatedKeys[key_idx+1]) {
			// remove ref at ref_idx + 1
		}
	}

	return prev, nil
}
