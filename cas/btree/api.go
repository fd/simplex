package btree

import (
	"github.com/fd/simplex/cas"
)

type Tree struct {
	root  *node_t
	store cas.GetterSetter
}

func New(s cas.GetterSetter) *Tree {
	return &Tree{store: s, root: &node_t{
		Type:         root_node_type | leaf_node_type,
		CollatedKeys: make([][]byte, 0, B),
		Children:     make([]*ref_t, 0, B+1),
		store:        s,
		parent:       nil,
		changed:      true,
	}}
}

func (t *Tree) String() string {
	return t.root.String()
}

func (t *Tree) Len() uint64 {
	return t.root.Len()
}

func (t *Tree) Get(collated_key []byte) (key, elt cas.Addr, err error) {
	ref, err := t.root.get(collated_key)
	if err != nil {
		return nil, nil, err
	}

	if ref == nil {
		return nil, nil, nil
	}

	return ref.Key, ref.Elt, nil
}

func (t *Tree) GetAt(idx uint64) (key, elt cas.Addr, err error) {
	return
}

func (t *Tree) Del(collated_key []byte) (key, elt cas.Addr, err error) {
	ref, err := t.root.remove_ref(collated_key, B)
	if err != nil {
		return nil, nil, err
	}

	if ref == nil {
		return nil, nil, nil
	}

	if len(t.root.Children) == 1 {
		t.root = t.root.Children[0].cache.(*node_t)
	}

	return ref.Key, ref.Elt, nil
}

func (t *Tree) DelAt(idx uint64) (key, elt cas.Addr, err error) {
	return
}

func (t *Tree) Set(collated_key []byte, key, elt cas.Addr) (prev cas.Addr, err error) {
	ref := build_ref(t.store, key, elt)

	prev_ref, root_node, err := t.root.insert_ref(collated_key, ref, B)
	if err != nil {
		return nil, err
	}

	if root_node != nil {
		t.root = root_node
	}

	if prev_ref != nil && prev_ref.Flags&elt_is_set > 0 {
		prev = prev_ref.Elt
	}

	return prev, nil
}
