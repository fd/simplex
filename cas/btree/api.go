package btree

import (
	"github.com/fd/simplex/cas"
)

type Tree struct {
	RootAddr cas.Addr
	Len      uint64

	root  *node_t
	store cas.GetterSetter
}

/*
  Make a new B+Tree
*/
func New(s cas.GetterSetter) *Tree {
	return &Tree{
		store: s,
		root: &node_t{
			Type:         root_node_type | leaf_node_type,
			CollatedKeys: make([][]byte, 0, B),
			Children:     make([]*ref_t, 0, B+1),
			changed:      true,
		},
	}
}

func Open(s cas.GetterSetter, addr cas.Addr) (*Tree, error) {
	t := &Tree{
		store: s,
		root: &node_t{
			Type: root_node_type | leaf_node_type,
		},
	}

	err := cas.Decode(s, addr, t)
	if err != nil {
		return nil, err
	}

	err = cas.Decode(s, t.RootAddr, t.root)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Tree) String() string {
	return t.root.String()
}

/*
  Get the key and element addr (cas.Addr) for a collated key (cas.Collate(...)).
  This function return a nil key and elt when the collated key is not found.
*/
func (t *Tree) Get(collated_key []byte) (key, elt cas.Addr, err error) {
	ref, err := t.root.get(collated_key, t.store)
	if err != nil {
		return nil, nil, err
	}

	if ref == nil {
		return nil, nil, nil
	}

	return ref.Key, ref.Elt, nil
}

/*
  Get the key and element addr (cas.Addr) at the specified index.
  This function return a nil key and elt when the index is out of bounds.
*/
func (t *Tree) GetAt(idx uint64) (key, elt cas.Addr, err error) {
	ref, err := t.root.get_at(idx, t.store)
	if err != nil {
		return nil, nil, err
	}

	if ref == nil {
		return nil, nil, nil
	}

	return ref.Key, ref.Elt, nil
}

/*
  Delete the key and element for a collated key (cas.Collate(...)).
  The old key and element addresses are retuned if they were previously set.
  This function return a nil key and elt when the there were no previous values.
*/
func (t *Tree) Del(collated_key []byte) (key, elt cas.Addr, err error) {
	ref, err := t.root.remove_ref(collated_key, B, t.store)
	if err != nil {
		return nil, nil, err
	}

	if ref == nil {
		return nil, nil, nil
	}

	if len(t.root.Children) == 1 && t.root.Type&leaf_node_type == 0 {
		t.root = t.root.Children[0].cache.(*node_t)
	}

	t.Len = t.root.Len()

	return ref.Key, ref.Elt, nil
}

/*
  Delete the key and element at the specified index.
  The old key and element addresses are retuned if they were previously set.
  This function return a nil key and elt when the index is out of bounds.
*/
func (t *Tree) DelAt(idx uint64) (key, elt cas.Addr, err error) {
	return
}

/*
  Set the key and element address (cas.Addr) for a collated key (cas.Collate(...)).
  The old key and element addresses are retuned if they were previously set.
  This function return a nil key and elt when the there were no previous values.
*/
func (t *Tree) Set(collated_key []byte, key, elt cas.Addr) (prev cas.Addr, err error) {
	ref := build_ref(t.store, key, elt)

	prev_ref, root_node, err := t.root.insert_ref(collated_key, ref, B, t.store)
	if err != nil {
		return nil, err
	}

	if root_node != nil {
		t.root = root_node
	}

	if prev_ref != nil && prev_ref.Flags&elt_is_set > 0 {
		prev = prev_ref.Elt
	}

	t.Len = t.root.Len()

	return prev, nil
}

func (t *Tree) Commit() (cas.Addr, error) {
	if t.root.changed {

		addr, err := commit(t.root, t.store)
		if err != nil {
			return nil, err
		}
		t.RootAddr = addr
	}

	return cas.Encode(t.store, &t, -1)
}
