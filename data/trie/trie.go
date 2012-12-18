package trie

import (
	"fmt"
	"sort"
	"strings"
	"unsafe"
)

type T struct {
	nodes []node_t
}

func New() *T {
	t := new(T)
	t.nodes = make([]node_t, 1, 64)
	return t
}

type node_t struct {
	chunk    []byte      // 4 + A
	children []uint      // 4 + (4 * B)
	value    interface{} // 4
}

func (t *T) Lookup(key []byte) (val interface{}, found bool) {
	return exec(t, 0, key, lookup, nil)
}

func (t *T) Insert(key []byte, val interface{}) {
	exec(t, 0, key, insert, val)
}

func (t *T) Remove(key []byte) (val interface{}, found bool) {
	return exec(t, 0, key, remove, nil)
}

func (t *T) ConsumedMemory() uintptr {
	m := unsafe.Sizeof(T{})
	m += unsafe.Sizeof(node_t{}) * uintptr(cap(t.nodes))

	for _, n := range t.nodes {
		m += n.ConsumedMemory()
	}

	return m
}

func (n *node_t) ConsumedMemory() uintptr {
	m := uintptr(len(n.chunk))
	m += uintptr(cap(n.children)) * unsafe.Sizeof(uint(0))
	return m
}

func (t *T) String() string {
	return t.nodes[0].String(t)
}

func (n *node_t) String(t *T) string {
	c := []string{}
	for _, id := range n.children {
		c = append(c, strings.Replace(t.nodes[id].String(t), "\n", "\n  ", -1))
	}
	return fmt.Sprintf("{\n  K: %s,\n  V: %+v\n  %s\n}", string(n.chunk), n.value, strings.Join(c, "\n  "))
}

func (t *T) grow() {
	if len(t.nodes) == cap(t.nodes) {
		n := make([]node_t, len(t.nodes), cap(t.nodes)*2)
		copy(n, t.nodes)
		t.nodes = n
	}
	t.nodes = t.nodes[:len(t.nodes)+1]
}

type operation byte

const (
	lookup operation = iota
	insert
	remove
)

func exec(t *T, node_id uint, suffix []byte, op operation, val interface{}) (interface{}, bool) {
	node := &t.nodes[node_id]

	if len(suffix) == 0 {
		switch op {
		case insert:
			node.value = val
		case remove:
			val = node.value
			node.value = nil
		case lookup:
			val = node.value
		}
		return val, true
	}

	var n *node_t
	var n_idx int
	var n_id uint

	// look for the branch to continue on
	for i, child_id := range node.children {
		child := &t.nodes[child_id]

		if len(child.chunk) == 0 {
			panic("zero length prefixes must be merged with self")
		}

		if child.chunk[0] == suffix[0] {
			n = child
			n_id = child_id
			n_idx = i
		}
	}

	// found a prefix
	if n == nil {
		if op == insert {
			p, _ := push(t, node_id, suffix)
			p.value = val

			return p.value, true
		} else {
			return nil, false
		}
	}

	chunk := n.chunk
	l_s := len(suffix)
	l_c := len(chunk)

	// min lenth
	l_m := l_s
	if l_c < l_s {
		l_m = l_c
	}

	for i := 0; i < l_m; i++ {
		if suffix[i] != chunk[i] {
			if op == insert {
				p, p_id := split(t, n_id, i)
				p.value = nil

				if l_s > i {
					return exec(t, p_id, suffix[i:], op, val)
				} else {
					p.value = val
					return p.value, true
				}

			} else {
				return nil, false
			}
		}
	}

	if l_s == l_c {
		if op == insert {
			n.value = val
			return val, true

		} else if op == remove {
			val := n.value
			n.value = nil
			l_t := len(node.children)

			if len(n.children) == 0 {
				if l_t > n_idx {
					copy(node.children[n_idx:], node.children[n_idx+1:])
				}
				node.children = node.children[:l_t-1]
				// TODO shrink t.nodes (remove n_id)
			}

			if len(n.children) == 1 {
				c_id := n.children[0]
				c := &t.nodes[c_id]
				c.chunk = append(n.chunk, c.chunk...)
				node.children[n_idx] = c_id
				// TODO shrink t.nodes (remove n_id)
			}

			return val, true
		}

		return n.value, true
	} else if l_s < l_c {
		if op == insert {
			p, _ := split(t, n_id, len(suffix))
			p.value = val
			return p.value, true

		}

		return nil, true
	}

	return exec(t, n_id, suffix[l_c:], op, val)
}

func split(t *T, b uint, offset int) (*node_t, uint) {
	c := uint(len(t.nodes))

	t.grow()

	t.nodes[c], c, b = t.nodes[b], b, c
	n_c := &t.nodes[c]
	n_b := &t.nodes[b]

	// update c node
	n_c.chunk = n_b.chunk[:offset]
	n_c.value = nil
	n_c.children = []uint{b}

	// update b node
	n_b.chunk = n_b.chunk[offset:]

	return n_c, c
}

func push(t *T, a uint, suffix []byte) (*node_t, uint) {
	b := uint(len(t.nodes))

	t.grow()

	n_a := &t.nodes[a]
	n_b := &t.nodes[b]

	n_b.chunk = make([]byte, len(suffix))
	n_b.value = nil
	n_b.children = nil
	copy(n_b.chunk, suffix)

	//      node.grow()
	n_a.children = append(n_a.children, b)
	sort.Sort(t_sort{t, n_a.children})

	return n_b, b
}

type t_sort struct {
	t *T
	c []uint
}

func (l t_sort) Len() int {
	return len(l.c)
}

func (l t_sort) Less(x, y int) bool {
	x_id := l.c[x]
	y_id := l.c[y]
	return l.t.nodes[x_id].chunk[0] < l.t.nodes[y_id].chunk[0]
}

func (l t_sort) Swap(x, y int) {
	l.c[x], l.c[y] = l.c[y], l.c[x]
}
