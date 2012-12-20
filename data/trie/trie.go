package trie

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
	Set      bool
	Chunk    []byte
	Children []*node_t
	Value    interface{}
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

type operation byte

const (
	lookup operation = iota
	insert
	remove
)

func exec(node *node_t, suffix []byte, op operation, val interface{}) (interface{}, bool) {
	var found bool

	if node == nil {
		return nil, false
	}

	if len(suffix) == 0 {
		switch op {
		case insert:
			found, node.Set = true, true
			val, node.Value = node.Value, val
		case remove:
			found, node.Set = node.Set, false
			val, node.Value = node.Value, nil
		case lookup:
			found = node.Set
			val = node.Value
		}
		// CASE 1.A 1.B 1.C
		//fmt.Println("CASE 1x")
		return val, found
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
		}
	}

	// found a prefix
	if n == nil {
		if op == insert {
			c := push(node, suffix)
			found, c.Set = true, true
			val, c.Value = c.Value, val

			// CASE 2
			//fmt.Println("CASE 2")
			return val, found
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
				return exec(c, suffix[i:], op, val)

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
				return exec(c, suffix[i:], op, val)

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
			return exec(n, suffix[l_c:], op, val)

		case remove:
			val, n.Value = n.Value, nil
			found, n.Set = n.Set, false
			optimize(node, n, n_idx)
			// CASE 9A -> no optimize
			//      9B -> optimize (remove leaf)
			//      9C -> optimize (remove branch)
			//fmt.Println("CASE 9x")
			return val, found

		case lookup:
			// CASE 10
			//fmt.Println("CASE 10")
			return n.Value, n.Set

		}

	} else if l_s > l_c {
		// CASE 11
		//fmt.Println("CASE 11 (cont)")
		return exec(n, suffix[l_c:], op, val)

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
		p.Children[a_idx] = c
	}
}

func split(p, a *node_t, a_idx int, offset int) *node_t {
	b := &node_t{}

	// update b node
	b.Chunk = a.Chunk[:offset]
	b.Children = append(b.Children, a)

	// update a node
	a.Chunk = a.Chunk[offset:]

	// update p node
	p.Children[a_idx] = b

	return b
}

func push(a *node_t, suffix []byte) *node_t {
	b := &node_t{}

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
