package btree

import (
	"github.com/fd/simplex/cas"
)

func commit(n *node_t, store cas.GetterSetter) (cas.Addr, error) {
	for _, ref := range n.Children {
		if child, ok := ref.cache.(*node_t); ok &&
			child != nil {

			if child.changed {
				addr, err := commit(child, store)
				if err != nil {
					return nil, err
				}

				ref.Key = nil
				ref.Elt = addr
				child.changed = false
			}

			child.ref = nil
			ref.cache = nil
		}
	}

	return cas.Encode(store, n)
}
