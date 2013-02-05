package btree

import (
	"errors"
	"github.com/fd/simplex/cas"
)

var EOI = errors.New("end of iterator")

type Iter struct {
	node *node_t
	curr int

	child *Iter
	store cas.GetterSetter
}

func (i *Iter) Next() (key, elt cas.Addr, err error) {
	if i.child != nil {
		k, e, er := i.child.Next()

		if er == EOI {
			i.child = nil
			return i.Next()
		}

		if er != nil {
			return nil, nil, er
		}

		return k, e, nil
	}

	if i.curr >= len(i.node.Children) {
		return nil, nil, EOI
	}

	if i.node.Type&leaf_node_type == 0 {
		it, er := iter_at(i.node, i.curr, i.store)
		if er != nil {
			return nil, nil, err
		}

		i.curr += 1
		i.child = it
		return i.Next()
	}

	ref := i.node.Children[i.curr]
	i.curr += 1
	return ref.Key, ref.Elt, nil
}

func iter_at(n *node_t, idx int, store cas.GetterSetter) (*Iter, error) {
	node, err := n.Children[idx].load_node(store)
	if err != nil {
		return nil, err
	}

	return &Iter{node, 0, nil, store}, nil
}
