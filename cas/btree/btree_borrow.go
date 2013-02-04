package btree

func borrow(dst, src *node_t, to_begin bool, placeholder_key []byte, order int) (new_placeholder_key []byte, ok bool) {
	if dst == nil || src == nil {
		return nil, false
	}

	if dst.min_children(order) < len(dst.Children) {
		return nil, false
	}

	if src.min_children(order) >= len(src.Children) {
		return nil, false
	}

	var (
		dst_keys       = dst.CollatedKeys
		dst_children   = dst.Children
		dst_keys_n     = len(dst_keys)
		dst_children_n = len(dst_children)

		src_keys       = src.CollatedKeys
		src_children   = src.Children
		src_keys_n     = len(src_keys)
		src_children_n = len(src_children)

		borrow_n = (src_children_n - (src.min_children(order) - 1)) / 2
	)

	if borrow_n == 0 {
		return nil, false
	}

	dst_keys = dst_keys[:dst_keys_n+borrow_n]
	dst_children = dst_children[:dst_children_n+borrow_n]

	if to_begin {
		// [1 2 3] 4 [5 6 7]
		// [1] 2 [3 4 5 6 7]
		if src.Type&leaf_node_type > 0 {
			copy(dst_keys[borrow_n:], dst_keys)
			copy(dst_keys, src_keys[src_keys_n-borrow_n:])
			new_placeholder_key = dst_keys[0]
		} else {
			copy(dst_keys[borrow_n:], dst_keys)
			dst_keys[borrow_n-1] = placeholder_key
			new_placeholder_key = src_keys[src_keys_n-borrow_n]
			copy(dst_keys, src_keys[src_keys_n-borrow_n+1:])
		}

		copy(dst_children[borrow_n:], dst_children)
		copy(dst_children, src_children[src_children_n-borrow_n:])

	} else {
		// [1 2 3] 4 [5 6 7]
		// [1 2 3 4 5] 6 [7]
		if src.Type&leaf_node_type > 0 {
			copy(dst_keys[dst_keys_n-borrow_n+1:], src_keys)
			copy(src_keys, src_keys[borrow_n:])
			new_placeholder_key = src_keys[0]
		} else {
			copy(dst_keys[dst_keys_n-borrow_n+2:], src_keys)
			dst_keys[dst_keys_n-borrow_n+1] = placeholder_key
			new_placeholder_key = src_keys[borrow_n-1]
			copy(src_keys, src_keys[borrow_n:])
		}

		copy(dst_children[dst_children_n-borrow_n+1:], src_children)
		copy(src_children, src_children[borrow_n:])
	}

	src_keys = src_keys[:src_keys_n-borrow_n]
	src_children = src_children[:src_children_n-borrow_n]

	src.changed = true
	src.CollatedKeys = src_keys
	src.Children = src_children
	if src.ref != nil {
		src.ref.Len = src.Len()
	}

	dst.changed = true
	dst.CollatedKeys = dst_keys
	dst.Children = dst_children
	if dst.ref != nil {
		dst.ref.Len = dst.Len()
	}

	return new_placeholder_key, true
}
