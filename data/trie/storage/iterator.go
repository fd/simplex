package storage

import (
	"github.com/fd/simplex/data/storage"
)

type Iterator interface {
	Next() (sha storage.SHA, done bool)
}

type iterator struct {
	root  *node_t
	node  *node_t
	next  int
	child Iterator
}

func (n *node_t) Iter() Iterator {
	return &iterator{
		root: n,
		node: n,
	}
}

func (i *iterator) Next() (sha storage.SHA, done bool) {
	if i.next == 0 {
		i.next += 1
		if i.node.Set {
			return i.node.Value, false
		}
	}

TRY_NEXT:
	if i.child == nil {
		if len(i.node.Children) == (i.next - 1) {
			return storage.ZeroSHA, true
		}

		child_ref := i.node.Children[i.next-1]

		child_node, found := load_node(i.node.storage, child_ref.SHA)
		if !found {
			panic("corrupted trie node")
		}

		i.child = &iterator{
			root: i.root,
			node: child_node,
		}
	}

	if i.child != nil {
		sha, d := i.child.Next()
		if d {
			i.next += 1
			i.child = nil
			goto TRY_NEXT
		}
		return sha, false
	}

	panic("not reached")
}
