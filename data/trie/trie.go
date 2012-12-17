package trie

import (
	"sort"
)

type T struct {
	chunk    []byte
	children []*T
	value    interface{}
}

func (t *T) Lookup(key []byte) (val interface{}, found bool) {
	return t.exec(key, lookup, nil)
}

func (t *T) Insert(key []byte, val interface{}) {
	t.exec(key, insert, val)
}

func (t *T) Remove(key []byte) (val interface{}, found bool) {
	return t.exec(key, remove, nil)
}

type operation byte

const (
	lookup operation = iota
	insert
	remove
)

func (t *T) exec(suffix []byte, op operation, val interface{}) (interface{}, bool) {
	if len(suffix) == 0 {
		switch op {
		case insert:
			t.value = val
		case remove:
			val = t.value
			t.value = nil
		case lookup:
			val = t.value
		}
		return val, true
	}

	var n *T
	var n_idx int

	// look for the branch to continue on
	for i, child := range t.children {
		if len(child.chunk) == 0 {
			panic("zero length prefixes must be merged with self")
		}

		if child.chunk[0] == suffix[0] {
			n = child
			n_idx = i
		}
	}

	// found a prefix
	if n == nil {
		if op == insert {
			p := &T{chunk: suffix, value: val}
			t.children = append(t.children, p)
			sort.Sort(t_sort(t.children))
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
				p := &T{chunk: chunk[:i]}
				n.chunk = chunk[i:]
				p.children = append(p.children, n)
				t.children[n_idx] = p

				if l_s > i {
					return p.exec(suffix[i:], op, val)
				} else {
					p.value = val
				}

				return val, true
			} else {
				return nil, false
			}
		}
	}

	if l_s == l_m {
		if op == insert {
			n.value = val
			return val, true

		} else if op == remove {
			val := n.value
			n.value = nil
			l_t := len(t.children)

			if len(n.children) == 0 {
				if l_t > n_idx {
					copy(t.children[n_idx:], t.children[n_idx+1:])
				}
				t.children = t.children[:l_t-1]
			}

			if len(n.children) == 1 {
				c := n.children[0]
				c.chunk = append(n.chunk, c.chunk...)
				t.children[n_idx] = c
			}

			return val, true
		}

		return n.value, true
	}

	return n.exec(suffix[l_c:], op, val)
}

type t_sort []*T

func (l t_sort) Len() int {
	return len(l)
}

func (l t_sort) Less(x, y int) bool {
	return l[x].chunk[0] < l[y].chunk[0]
}

func (l t_sort) Swap(x, y int) {
	l[x], l[y] = l[y], l[x]
}
