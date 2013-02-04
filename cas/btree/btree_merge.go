package btree

func merge(dst, src *node_t, placeholder_key []byte, order int) bool {
	if dst == nil || src == nil {
		return false
	}

	if dst.min_children(order) < len(dst.Children) {
		return false
	}

	if src.min_children(order) < len(src.Children) {
		return false
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
	)

	if dst.Type&leaf_node_type > 0 {
		dst_keys = dst_keys[:dst_keys_n+src_keys_n]
		copy(dst_keys[dst_keys_n:], src_keys)
	} else {
		dst_keys = dst_keys[:dst_keys_n+src_keys_n+1]
		dst_keys[dst_keys_n] = placeholder_key
		copy(dst_keys[dst_keys_n+1:], src_keys)
	}

	dst_children = dst_children[:dst_children_n+src_children_n]
	copy(dst_children[dst_children_n:], src_children)

	dst.changed = true
	dst.CollatedKeys = dst_keys
	dst.Children = dst_children
	if dst.ref != nil {
		dst.ref.Len = dst.Len()
	}

	return true
}
