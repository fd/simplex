package btree

import (
	"errors"
	"github.com/fd/simplex/cas"
)

var EOI = errors.New("end of iterator")

type Iter struct {
	node *node_t
	curr int
	err  error

	child *Iter
	store cas.GetterSetter
}

func (t *Tree) Iter() *Iter {
	return &Iter{t.root, 0, nil, nil, t.store}
}

func (t *Tree) IterFrom(collated_key []byte) *Iter {
	return t.root.get_iter_from(collated_key, t.store)
}

func (t *Tree) IterAt(idx uint64) *Iter {
	return t.root.get_iter_at(idx, t.store)
}

func (n *node_t) get_iter_from(key []byte, store cas.GetterSetter) *Iter {
	_, idx, ref := n.search_ref(key)

	if n.Type&leaf_node_type > 0 {
		return &Iter{n, idx, nil, nil, store}
	}

	if ref == nil {
		return &Iter{n, 0, EOI, nil, store}
	}

	c, err := ref.load_node(store)
	if err != nil {
		return &Iter{n, 0, err, nil, store}
	}

	return &Iter{n, idx, nil, c.get_iter_from(key, store), store}
}

func (n *node_t) get_iter_at(idx uint64, store cas.GetterSetter) *Iter {
	if idx < 0 {
		return &Iter{n, 0, EOI, nil, store}
	}

	if idx >= n.Len() {
		return &Iter{n, 0, EOI, nil, store}
	}

	for i, ref := range n.Children {
		l := ref.Len

		if l > idx {
			c, err := ref.load_node(store)
			if err != nil {
				return &Iter{c, 0, err, nil, store}
			}

			return &Iter{n, i, nil, c.get_iter_at(idx, store), store}
		}

		idx -= l
	}

	return &Iter{n, 0, EOI, nil, store}
}

func (i *Iter) Next() (key, elt cas.Addr, err error) {
	if i.err != nil {
		return nil, nil, i.err
	}

	if i.child != nil {
		k, e, er := i.child.Next()

		if er == EOI {
			i.child = nil
			return i.Next()
		}

		if er != nil {
			i.err = er
			return i.Next()
		}

		return k, e, nil
	}

	if i.curr >= len(i.node.Children) {
		i.err = EOI
		return i.Next()
	}

	if i.node.Type&leaf_node_type == 0 {
		i.child = iter_at(i.node, i.curr, i.store)
		i.curr += 1
		return i.Next()
	}

	ref := i.node.Children[i.curr]
	i.curr += 1
	return ref.Key, ref.Elt, nil
}

func iter_at(n *node_t, idx int, store cas.GetterSetter) *Iter {
	node, err := n.Children[idx].load_node(store)
	return &Iter{node, 0, err, nil, store}
}
