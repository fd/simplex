package storage

import (
	"github.com/fd/simplex/data/storage"
	"sort"
)

type ReduceFunc func(reductions []storage.SHA, value storage.SHA) storage.SHA

type T interface {
	Store() *storage.S
	Commit() storage.SHA

	Lookup(key []byte) (storage.SHA, bool)
	Insert(key []byte, sha storage.SHA) (old_sha storage.SHA, inserted bool)
	Remove(key []byte) (old_sha storage.SHA, removed bool)
	Reduce(cache T, f ReduceFunc) storage.SHA
	Iter() Iterator

	Empty() bool
}

type node_t struct {
	sha     storage.SHA
	storage *storage.S
	parent  *node_t
	changed bool

	Set      bool
	Value    storage.SHA
	Children []*child_ref_t
}

type child_ref_t struct {
	node  *node_t
	Chunk []byte
	SHA   storage.SHA
}

func New(s *storage.S) T {
	return &node_t{
		storage: s,
		changed: true,
	}
}

func Get(s *storage.S, sha storage.SHA) (T, bool) {
	return load_node(s, sha)
}

func load_node(s *storage.S, sha storage.SHA) (*node_t, bool) {
	node := &node_t{
		storage: s,
		sha:     sha,
	}
	found := s.Get(sha, node)
	return node, found
}

func (s *node_t) Reduce(cache T, f ReduceFunc) storage.SHA {
	return s.reduce(cache, f, []byte{})
}

func (s *node_t) reduce(cache T, f ReduceFunc, prefix []byte) storage.SHA {
	reductions := make([]storage.SHA, len(s.Children))

	for i, c := range s.Children {
		if c.node == nil || !c.node.changed {
			sha, found := cache.Lookup(append(prefix, c.Chunk...))
			if found {
				reductions[i] = sha
				continue
			}
		}

		if c.node == nil {
			n, f := load_node(s.storage, c.SHA)
			if !f {
				panic("corrupted trie node")
			}
			n.parent = s
			c.node = n
		}

		reductions[i] = c.node.reduce(cache, f, append(prefix, c.Chunk...))
	}

	if len(reductions) == 1 && !s.Set {
		cache.Insert(prefix, reductions[0])
		return reductions[0]
	}

	sha := f(reductions, s.Value)
	cache.Insert(prefix, sha)

	return sha
}

func (s *node_t) Empty() bool {
	return s.Set == false && len(s.Children) == 0
}

func (s *node_t) Store() *storage.S {
	return s.storage
}

func (s *node_t) Commit() storage.SHA {
	if s.changed != true {
		return s.sha
	}

	for _, c := range s.Children {
		if c.node != nil && c.node.changed == true {
			sha := c.node.Commit()
			c.SHA = sha
		}
	}

	sha := s.storage.Set(s)
	s.sha = sha
	s.changed = false
	return sha
}

func (s *node_t) Lookup(key []byte) (storage.SHA, bool) {
	node, rem := s.lookup_node_chain(key)

	if len(rem) == 0 {
		if node.Value != storage.ZeroSHA {
			return node.Value, true
		}
	}

	return storage.ZeroSHA, false
}

func (s *node_t) Remove(key []byte) (old_sha storage.SHA, removed bool) {
	node, rem := s.lookup_node_chain(key)

	if len(rem) != 0 {
		return storage.ZeroSHA, false
	}

	if node.Set == false {
		return storage.ZeroSHA, true
	}

	old_sha = node.Value

	node.Set = false
	node.Value = storage.ZeroSHA
	node.cleanup()

	for node != nil {
		node.changed = true
		node = node.parent
	}

	return old_sha, true
}

func (s *node_t) Insert(key []byte, sha storage.SHA) (old_sha storage.SHA, inserted bool) {
	node, rem := s.lookup_node_chain(key)

	var (
		lower_split_child_ref *child_ref_t
		upper_split_child_ref *child_ref_t
		new_child_ref         *child_ref_t
		new_child             *node_t
		i                     int
	)

	if len(rem) == 0 {
		goto set_value
	}

	for _, c := range node.Children {
		if c.Chunk[0] == rem[0] {
			lower_split_child_ref = c
			goto split_node
		}
	}

	// push
	if lower_split_child_ref == nil {
		goto push_node
	}

split_node:
	i = 0
	for ; i < len(rem); i++ {
		if rem[i] != lower_split_child_ref.Chunk[i] {
			break
		}
	}

	new_child = &node_t{
		parent:   node,
		storage:  node.storage,
		Children: []*child_ref_t{lower_split_child_ref},
	}

	upper_split_child_ref = &child_ref_t{
		node:  new_child,
		Chunk: rem[:i],
		SHA:   storage.ZeroSHA,
	}

	lower_split_child_ref.Chunk = lower_split_child_ref.Chunk[i:]
	rem = rem[i:]

	if lower_split_child_ref.node != nil {
		lower_split_child_ref.node.parent = new_child
	}

	for i, c := range node.Children {
		if c == lower_split_child_ref {
			node.Children[i] = upper_split_child_ref
		}
	}

	node = new_child

	if len(rem) == 0 {
		goto set_value
	} else {
		goto push_node
	}

push_node:
	new_child = &node_t{
		parent:  node,
		storage: node.storage,
	}

	new_child_ref = &child_ref_t{
		node:  new_child,
		Chunk: rem,
		SHA:   storage.ZeroSHA,
	}

	node.Children = append(node.Children, new_child_ref)
	sort.Sort(t_sort(node.Children))

	rem = []byte{}
	node = new_child
	goto set_value

set_value:
	node.Set = true
	node.changed = true
	old_sha, node.Value = node.Value, sha

	for node != nil {
		node.changed = true
		node = node.parent
	}

	return old_sha, true
}

func (s *node_t) cleanup() {
	if len(s.Children) == 0 {
		s.remove_from_parent()
	} else if len(s.Children) == 1 {
		s.merge_with_only_child()
	}
}

func (s *node_t) remove_from_parent() {
	p := s.parent

	if p == nil {
		return
	}

	for i, c := range p.Children {
		if c.node == s {
			l := len(p.Children)

			if i < l-1 {
				copy(p.Children[i:], p.Children[i+1:])
			}

			p.Children = p.Children[:l-1]
			p.changed = true
			p.cleanup()

			break
		}
	}
}

func (s *node_t) merge_with_only_child() {
	p := s.parent

	if p == nil {
		return
	}

	var (
		c_ref = s.Children[0]
		s_ref *child_ref_t
	)

	for _, c := range p.Children {
		if c.node == s {
			s_ref = c
			break
		}
	}

	s_ref.node = c_ref.node
	s_ref.SHA = c_ref.SHA
	s_ref.Chunk = append(s_ref.Chunk, c_ref.Chunk...)
	p.changed = true
	p.cleanup()
}

func (s *node_t) lookup_node_chain(key []byte) (*node_t, []byte) {
	if len(key) == 0 {
		return s, key
	}

	var (
		child_ref  *child_ref_t
		i          int
		found      bool
		child_node *node_t
	)

	for _, c := range s.Children {
		if c.Chunk[0] == key[0] {
			child_ref = c
		}
	}

	// no child is found (push)
	if child_ref == nil {
		return s, key
	}

	for i = 0; i < len(key); i++ {

		// a prefix of key is equal to child.Chunk but key is longer (continue)
		if i >= len(child_ref.Chunk) {

			if child_ref.node == nil {
				child_node, found = load_node(s.storage, child_ref.SHA)
				if !found {
					panic("corrupted trie node")
				}
				child_node.parent = s
				child_ref.node = child_node
			} else {
				child_node = child_ref.node
			}

			return child_node.lookup_node_chain(key[i:])
		}

		// key shares a common prefix with child.Chunk (split)
		if child_ref.Chunk[i] != key[i] {
			return s, key
		}

	}

	// key and child.Chunk are equal (continue)
	if len(key) == len(child_ref.Chunk) {

		if child_ref.node == nil {
			child_node, found = load_node(s.storage, child_ref.SHA)
			if !found {
				panic("corrupted trie node")
			}
			child_node.parent = s
			child_ref.node = child_node
		} else {
			child_node = child_ref.node
		}

		return child_node.lookup_node_chain(key[i:])
	}

	// key is shorter than child_ref.Chunk (split)
	if len(key) < len(child_ref.Chunk) {
		return s, key
	}

	panic("not reached")
}

type t_sort []*child_ref_t

func (l t_sort) Len() int {
	return len(l)
}

func (l t_sort) Less(x, y int) bool {
	return l[x].Chunk[0] < l[y].Chunk[0]
}

func (l t_sort) Swap(x, y int) {
	l[x], l[y] = l[y], l[x]
}
