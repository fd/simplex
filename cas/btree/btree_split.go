package btree

func split(n *node_t, order int) (right_key []byte, right_ref *ref_t) {
	if n == nil {
		return nil, nil
	}

	if !n.has_too_many_children(order) {
		return nil, nil
	}

	// find the split idx :: floor(len(keys) / 2)
	split_idx := len(n.CollatedKeys) / 2

	right_node := &node_t{
		Type:         n.Type,
		CollatedKeys: make([][]byte, 0, order),
		Children:     make([]*ref_t, 0, order+1),
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
	}

	right_ref = &ref_t{
		Flags: elt_is_set | ref_is_nod,
		Len:   right_node.Len(),
		cache: right_node,
	}
	right_node.ref = right_ref

	right_key = right_node.CollatedKeys[0]

	if n.Type&leaf_node_type == 0 {
		right_node.CollatedKeys = right_node.CollatedKeys[1:]
	}

	return right_key, right_ref
}
