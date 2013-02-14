package runtime

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"simplex.sh/cas"
	"simplex.sh/cas/btree"
	"simplex.sh/runtime/event"
	"simplex.sh/runtime/promise"
)

func (op *group_op) Resolve(state promise.State, events chan<- event.Event) {
	var (
		src_events  = state.Resolve(op.src)
		table       = state.GetTable(op.name)
		table_cache = map[string]*sub_table_t{}
	)

	defer func() {
		for _, sub_table := range table_cache {
			if sub_table.events != nil {
				close(sub_table.events)
			}
		}
	}()

	for e := range src_events.C {
		// propagate error events
		if err, ok := e.(event.Error); ok {
			events <- err
			continue
		}

		i_change, ok := e.(*ChangedMember)
		if !ok {
			continue
		}

		var (
			coll_key_a []byte
			coll_key_b []byte
			key_b      interface{}
		)

		// calculate collated group key for a and b
		if i_change.a != nil {
			group_key := op.fun(&Context{state.Store()}, i_change.a)
			coll_key_a = cas.Collate(group_key)
		}
		if i_change.b != nil {
			key_b = op.fun(&Context{state.Store()}, i_change.b)
			coll_key_b = cas.Collate(key_b)
		}

		// propagate event
		// - to sub table at coll_key_b
		// - to groups table
		if bytes.Compare(coll_key_a, coll_key_b) == 0 {
			var (
				sub_table = get_sub_table(state.Store(), table, coll_key_b, table_cache)
			)

			if sub_table.key == nil {
				sub_table.key = key_b
			}

			prev_elt_addr, err := sub_table.tree.Set(coll_key_b, i_change.key, i_change.b)
			if err != nil {
				panic(err)
			}

			if cas.CompareAddr(prev_elt_addr, i_change.b) != 0 {
				if sub_table.events == nil {
					sub_table.events = get_sub_events(state, op.name, coll_key_b)
				}

				sub_table.changed = true
				sub_table.events <- &ChangedMember{op.name, coll_key_b, i_change.key, prev_elt_addr, i_change.b}
			}

			continue
		}

		// remove old entry from sub table
		if i_change.a != nil {
			var (
				sub_table = get_sub_table(state.Store(), table, coll_key_a, table_cache)
			)

			prev_key_addr, prev_elt_addr, err := sub_table.tree.Del(coll_key_a)
			if err != nil {
				panic(err)
			}

			if prev_key_addr != nil && prev_elt_addr != nil {
				if sub_table.events == nil {
					sub_table.events = get_sub_events(state, op.name, coll_key_a)
				}

				sub_table.changed = true
				sub_table.events <- &ChangedMember{op.name, coll_key_a, prev_key_addr, prev_elt_addr, nil}
			}
		}

		// add new entry to sub table (while potentially adding new subtables)
		if i_change.b != nil {
			var (
				sub_table = get_sub_table(state.Store(), table, coll_key_b, table_cache)
			)

			if sub_table.key == nil {
				sub_table.key = key_b
			}

			prev_elt_addr, err := sub_table.tree.Set(coll_key_b, i_change.key, i_change.b)
			if err != nil {
				panic(err)
			}

			if cas.CompareAddr(prev_elt_addr, i_change.b) != 0 {
				if sub_table.events == nil {
					sub_table.events = get_sub_events(state, op.name, coll_key_b)
				}

				sub_table.changed = true
				sub_table.events <- &ChangedMember{op.name, coll_key_b, i_change.key, prev_elt_addr, i_change.b}
			}
		}
	}

	// remove empty sub tables
	for _, sub_table := range table_cache {
		// delete table
		if sub_table.tree.Len == 0 {
			prev_key_addr, prev_elt_addr, err := table.Del(sub_table.coll_key)
			if err != nil {
				panic(err)
			}

			if prev_key_addr != nil && prev_elt_addr != nil {
				events <- &ChangedMember{op.name, sub_table.coll_key, prev_key_addr, prev_elt_addr, nil}
			}

			// don't commit this table
			sub_table.changed = false
		}

		// commit table
		if sub_table.changed {
			elt_addr, err := sub_table.tree.Commit()
			if err != nil {
				panic(err)
			}

			if sub_table.key_addr == nil {
				key_addr, err := cas.Encode(state.Store(), sub_table.key, -1)
				if err != nil {
					panic(err)
				}

				sub_table.key_addr = key_addr
			}

			prev_elt_addr, err := table.Set(sub_table.coll_key, sub_table.key_addr, elt_addr)
			if err != nil {
				panic(err)
			}

			if cas.CompareAddr(prev_elt_addr, elt_addr) != 0 {
				sub_table.events <- &ConsistentTable{op.name, prev_elt_addr, elt_addr}
				events <- &ChangedMember{op.name, sub_table.coll_key, sub_table.key_addr, prev_elt_addr, elt_addr}
			}
		}
	}

	tab_addr_a, tab_addr_b := state.CommitTable(op.name, table)
	events <- &ConsistentTable{op.name, tab_addr_a, tab_addr_b}
}

type sub_table_t struct {
	coll_key []byte
	key_addr cas.Addr
	key      interface{}
	tree     *btree.Tree
	events   chan<- event.Event
	changed  bool
}

func get_sub_table(store cas.Store, super_table *btree.Tree, coll_key []byte, cache map[string]*sub_table_t) *sub_table_t {
	var (
		tree         *btree.Tree
		changed      bool
		coll_key_str = string(coll_key)
	)

	if sub_table, ok := cache[coll_key_str]; ok {
		return sub_table
	}

	key_addr, sub_table_addr, err := super_table.Get(coll_key)
	if err != nil {
		panic(err)
	}

	if sub_table_addr == nil {
		tree = GetTable(store, sub_table_addr)
	} else {
		tree = btree.New(store)
		changed = true
	}

	if tree == nil {
		return nil
	}

	sub_table := &sub_table_t{
		coll_key,
		key_addr,
		nil,
		tree,
		nil,
		changed,
	}

	cache[coll_key_str] = sub_table
	return sub_table
}

func get_sub_events(state promise.State, super_name string, sub_name []byte) chan<- event.Event {
	sha := sha1.New()
	sha.Write(sub_name)
	name := super_name + "#" + hex.EncodeToString(sha.Sum(nil))

	return state.RegisterPublisher(name)
}
