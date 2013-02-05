package btree

import (
	"bytes"
	"sort"

	"github.com/fd/simplex/cas"
)

func (n *node_t) get(key []byte, store cas.GetterSetter) (ref *ref_t, err error) {
	_, _, ref = n.search_ref(key)

	if n.Type&leaf_node_type > 0 {
		return ref, nil
	}

	if ref == nil {
		return nil, nil
	}

	n, err = ref.load_node(store)
	if err != nil {
		return nil, err
	}

	return n.get(key, store)
}

func (n *node_t) get_at(idx uint64, store cas.GetterSetter) (ref *ref_t, err error) {
	if idx < 0 {
		return nil, nil
	}

	if idx >= n.Len() {
		return nil, nil
	}

	for _, ref := range n.Children {
		l := ref.Len

		if l > idx {
			n, err = ref.load_node(store)
			if err != nil {
				return nil, err
			}

			return n.get_at(idx, store)
		}

		idx -= l
	}

	return nil, nil
}

func (n *node_t) search_ref(key []byte) (key_idx, ref_idx int, ref *ref_t) {

	if len(n.CollatedKeys) == 0 {
		return 0, 0, nil
	}

	// find ref idx
	ref_idx = sort.Search(len(n.CollatedKeys), func(i int) bool {
		return bytes.Compare(n.CollatedKeys[i], key) > 0
	})

	key_idx = ref_idx - 1

	if n.Type&leaf_node_type > 0 {
		ref_idx -= 1

		if key_idx < 0 {
			key_idx += 1
			ref_idx += 1
			ref = nil
			return
		}

		if bytes.Compare(n.CollatedKeys[key_idx], key) < 0 {
			key_idx += 1
			ref_idx += 1
			ref = nil
			return
		}
	}

	if ref_idx >= 0 && ref_idx < len(n.Children) {
		ref = n.Children[ref_idx]
	}

	return
}
