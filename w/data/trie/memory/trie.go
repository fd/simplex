package memory

import (
	"fmt"
	"sort"
	"strings"
	"unsafe"
)

type T struct {
	root *node_t
}

func New() *T {
	t := new(T)
	t.root = &node_t{}
	return t
}

type node_t struct {
	parent   *node_t
	Set      bool
	Chunk    []byte
	Children []*node_t
	Value    interface{}
	Len      int
}

func (t *T) Lookup(key []byte) (val interface{}, found bool) {
	return exec(t.root, key, lookup, nil)
}

func (t *T) Insert(key []byte, val interface{}) (old_val interface{}, found bool) {
	return exec(t.root, key, insert, val)
}

func (t *T) Remove(key []byte) (old_val interface{}, found bool) {
	return exec(t.root, key, remove, nil)
}

func (t *T) Len() int {
	return t.root.Len
}

func (t *T) At(index int) (key []byte, val interface{}, found bool) {
	node, found := t.root.At(index)

	if found {
		c := 0

		n := node
		for n != nil {
			c += len(n.Chunk)
			n = n.parent
		}

		key = make([]byte, c)

		n = node
		for n != nil {
			c -= len(n.Chunk)
			copy(key[c:], n.Chunk)
			n = n.parent
		}
	}

	return key, node.Value, found
}

func (n *node_t) At(index int) (*node_t, bool) {
	min := 0
	max := 0

	if index >= n.Len {
		return nil, false
	}

	if n.Set && index == 0 {
		return n, true
	}

	if n.Set {
		min = 1
		max = 1
	}

	for _, c := range n.Children {
		max = min + c.Len
		if index < max {
			return c.At(index - min)
		}
		min = max
	}

	return nil, false
}

func (t *T) ConsumedMemory() uintptr {
	m := unsafe.Sizeof(T{})
	m += t.root.ConsumedMemory()

	return m
}

func (n *node_t) ConsumedMemory() uintptr {
	m := unsafe.Sizeof(node_t{})
	m += uintptr(len(n.Chunk))
	m += uintptr(cap(n.Children)) * unsafe.Sizeof(&node_t{})
	for _, c := range n.Children {
		m += c.ConsumedMemory()
	}
	return m
}

func (t *T) String() string {
	return t.root.String()
}

func (n *node_t) String() string {
	if len(n.Children) > 0 {
		cs := []string{}
		for _, c := range n.Children {
			cs = append(cs, strings.Replace(c.String(), "\n", "\n  ", -1))
		}
		return fmt.Sprintf("{\n  K: %s,\n  V: %+v,\n  %s\n}", string(n.Chunk), n.Value, strings.Join(cs, "\n  "))
	}
	return fmt.Sprintf("{ K: %s, V: %+v }", string(n.Chunk), n.Value)
}

func (t *T) Values() []interface{} {
	l := make([]interface{}, t.Len())
	t.root.Values(l)
	return l
}

func (n *node_t) Values(tip []interface{}) []interface{} {
	if n.Set {
		tip[0] = n.Value
		tip = tip[1:]
	}

	for _, c := range n.Children {
		tip = c.Values(tip)
	}

	return tip
}

type operation byte

const (
	lookup operation = iota
	insert
	remove
)

func exec(node *node_t, suffix []byte, op operation, val interface{}) (interface{}, bool) {
	node, found := exec_node(node, suffix, op)

	if !found {
		return nil, false
	}

	switch op {
	case insert:
		// update counts
		if !node.Set {
			n := node
			for n != nil {
				n.Len += 1
				n = n.parent
			}
		}

		found, node.Set = true, true
		val, node.Value = node.Value, val

	case remove:
		// update counts
		if node.Set {
			n := node
			for n != nil {
				n.Len -= 1
				n = n.parent
			}
		}

		found, node.Set = node.Set, false
		val, node.Value = node.Value, nil

	case lookup:
		found = node.Set
		val = node.Value

	}

	return val, found
}

func exec_node(node *node_t, suffix []byte, op operation) (*node_t, bool) {

	if node == nil {
		return nil, false
	}

	if len(suffix) == 0 {
		return node, true
	}

	var n *node_t
	var n_idx int

	// look for the branch to continue on
	for i, c := range node.Children {
		if len(c.Chunk) == 0 {
			panic("zero length prefixes must be merged with self")
		}

		if c.Chunk[0] == suffix[0] {
			n = c
			n_idx = i
			break
		}
	}

	// found a prefix
	if n == nil {
		if op == insert {
			c := push(node, suffix)
			// CASE 2
			//fmt.Println("CASE 2")
			return c, true
		} else {
			// CASE 3A 3B
			//fmt.Println("CASE 3x")
			return nil, false
		}
	}

	chunk := n.Chunk
	l_s := len(suffix)
	l_c := len(chunk)

	for i := 0; i < l_c; i++ {

		if i >= l_s {
			if op == insert {
				c := split(node, n, n_idx, i)
				// CASE 4
				//fmt.Println("CASE 4")
				return c, true

			} else {
				// CASE 5A 5B
				//fmt.Println("CASE 5x")
				return nil, false
			}
		}

		if suffix[i] != chunk[i] {
			if op == insert {
				c := split(node, n, n_idx, i)
				// CASE 6
				//fmt.Println("CASE 6")
				return exec_node(c, suffix[i:], op)

			} else {
				// CASE 7A 7B
				//fmt.Println("CASE 7x")
				return nil, false
			}
		}

	}

	if l_s == l_c {
		switch op {
		case insert:
			// CASE 8
			//fmt.Println("CASE 8")
			return n, true

		case remove:
			// CASE 9A -> no optimize
			//      9B -> optimize (remove leaf)
			//      9C -> optimize (remove branch)
			//fmt.Println("CASE 9x")
			optimize(node, n, n_idx)
			return n, true

		case lookup:
			// CASE 10
			//fmt.Println("CASE 10")
			return n, true

		}

	} else if l_s > l_c {
		// CASE 11
		//fmt.Println("CASE 11 (cont)")
		return exec_node(n, suffix[l_c:], op)

	}

	panic("not reached")
}

func optimize(p, a *node_t, a_idx int) {
	if len(a.Children) == 0 {
		l := len(p.Children)
		copy(p.Children[a_idx:], p.Children[a_idx+1:])
		p.Children = p.Children[:l-1]
	}

	if len(a.Children) == 1 {
		c := a.Children[0]
		c.Chunk = append(a.Chunk, c.Chunk...)
		c.parent = p
		p.Children[a_idx] = c
	}
}

func split(p, a *node_t, a_idx int, offset int) *node_t {
	b := &node_t{}

	// update b node
	b.parent = p
	b.Chunk = a.Chunk[:offset]
	b.Children = append(b.Children, a)
	b.Len = a.Len

	// update a node
	a.Chunk = a.Chunk[offset:]
	a.parent = b

	// update p node
	p.Children[a_idx] = b

	return b
}

func push(a *node_t, suffix []byte) *node_t {
	b := &node_t{}

	b.parent = a
	b.Chunk = make([]byte, len(suffix))
	copy(b.Chunk, suffix)

	a.Children = append(a.Children, b)
	sort.Sort(t_sort(a.Children))

	return b
}

type t_sort []*node_t

func (l t_sort) Len() int {
	return len(l)
}

func (l t_sort) Less(x, y int) bool {
	return l[x].Chunk[0] < l[y].Chunk[0]
}

func (l t_sort) Swap(x, y int) {
	l[x], l[y] = l[y], l[x]
}
