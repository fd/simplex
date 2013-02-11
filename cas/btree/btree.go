package btree

import (
	"bytes"
	"fmt"
	"simplex.sh/cas"
	"strings"
)

const B = 256

const (
	root_node_type node_type_t = 1 << iota
	inner_node_type
	leaf_node_type
)

const (
	key_is_set ref_flags_t = 1 << iota
	elt_is_set
	ref_is_val
	ref_is_nod
)

type (
	node_type_t byte
	ref_flags_t byte

	node_t struct {
		Type node_type_t

		CollatedKeys [][]byte
		Children     []*ref_t

		ref     *ref_t
		changed bool
	}

	ref_t struct {
		Flags ref_flags_t
		Len   uint64
		Key   cas.Addr
		Elt   cas.Addr

		cache interface{}
	}
)

func (n *node_t) Len() uint64 {
	var l uint64
	for _, ref := range n.Children {
		if ref != nil {
			l += ref.Len
		}
	}
	return l
}

func (n *node_t) String() string {
	var (
		b bytes.Buffer
		j = 0
	)

	fmt.Fprintf(&b, "[NODE %d:%d]\n", len(n.Children), n.Len())

	if n.Type&leaf_node_type == 0 {
		ref := n.Children[j]
		fmt.Fprintf(&b, "- k: [BEFORE]\n  v: %+v\n", strings.Replace(ref.String(), "\n", "\n  ", -1))
		j += 1
	}

	for i, key := range n.CollatedKeys {
		if i != 0 {
			b.WriteByte('\n')
		}
		ref := "Ref(nil)"
		if j < len(n.Children) {
			ref = n.Children[j].String()
		}
		fmt.Fprintf(&b, "- k: %x\n  v: %+v", key, strings.Replace(ref, "\n", "\n  ", -1))
		j += 1
	}

	return b.String()
}

func (ref *ref_t) String() string {
	if ref == nil {
		return "Ref(nil)"
	}

	if ref.Flags&ref_is_nod > 0 {
		c := ref.cache.(*node_t)
		return c.String()
	}

	return fmt.Sprintf("`%s` => `%s`", ref.Key, ref.Elt)
}

func (n *node_t) max_children(order int) int {
	if n.Type&root_node_type > 0 && n.Type&leaf_node_type > 0 {
		return order - 1
	}

	if n.Type&root_node_type > 0 {
		return order
	}

	if n.Type&inner_node_type > 0 {
		return order
	}

	if n.Type&leaf_node_type > 0 {
		return order - 1
	}

	panic("not reached")
}

func (n *node_t) min_children(order int) int {
	if n.Type&root_node_type > 0 && n.Type&leaf_node_type > 0 {
		return 1
	}

	if n.Type&root_node_type > 0 {
		return 2
	}

	if n.Type&inner_node_type > 0 {
		return order / 2
	}

	if n.Type&leaf_node_type > 0 {
		return order / 2
	}

	panic("not reached")
}

func (n *node_t) has_too_many_children(order int) bool {
	return len(n.Children) > n.max_children(order)
}

func (n *node_t) has_too_few_children(order int) bool {
	return len(n.Children) < n.min_children(order)
}

func (ref *ref_t) load_node(store cas.GetterSetter) (*node_t, error) {
	if ref.Flags&elt_is_set == 0 {
		panic("corrupt btree ref")
	}

	if ref.Flags&ref_is_nod == 0 {
		panic("corrupt btree ref")
	}

	if node, ok := ref.cache.(*node_t); ok && node != nil {
		return node, nil
	}

	var (
		node = &node_t{}
	)

	err := cas.NewDecoder(store, ref.Elt).Decode(&node)
	if err != nil {
		return nil, err
	}

	node.ref = ref
	ref.cache = node
	return node, nil
}

func build_ref(store cas.GetterSetter, key, elt cas.Addr) *ref_t {
	var (
		ref = &ref_t{Len: 1, Flags: ref_is_val}
	)

	ref.Flags |= key_is_set
	ref.Key = key

	ref.Flags |= elt_is_set
	ref.Elt = elt

	return ref
}
