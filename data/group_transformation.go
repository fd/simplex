package data

import (
	"sort"
)

type GroupFunc func(Context, Value) Value
type GroupNFunc func(Context, Value) []Value

type GroupView struct {
	View
	transformation *group_transformation
}

/*
  Returns a View referencing the members of a group.

  The view looks like this:

    [
      // member objects
      // ...
    ]
*/
func (gv GroupView) Members() View {
	return gv.View
}

/*
  Returns a View referencing the groups as collection members.

  The view looks like this:

    [
      {
        "key":        ...,
        "member_ids": [ ... ]
      }
    ]
*/
func (gv GroupView) Groups() View {
	v := gv.View
	v.current = &group_sidestream{gv.transformation}
	return v
}

/*
  This is a shorthand for data.ScopedView().Group(f)
*/
func Group(f GroupFunc) GroupView {
	return current_engine.ScopedView().Group(f)
}

/*
  This is a shorthand for data.ScopedView().GroupN(f)
*/
func GroupN(f GroupNFunc) GroupView {
	return current_engine.ScopedView().GroupN(f)
}

/*
  Group the members of a collection by a Value.
  Each member can only be in one group at a time.

    var Upstream = ...
    var GroupedByCategory = Upstream.Group(by_category)

    func by_category (ctx data.Context, val data.Value) data.Value {
      return val.Get("category")
    }
*/
func (v View) Group(f GroupFunc) GroupView {
	return v.GroupN(func(ctx Context, val Value) []Value {
		return []Value{f(ctx, val)}
	})
}

/*
  Group the members of a collection by a Value.
  Each member can be in multiple groups at a time.

    var Upstream = ...
    var GroupedByTag = Upstream.GroupN(by_tag)

    func by_tag (ctx data.Context, val data.Value) []data.Value {
      return val.Get("tags")
    }
*/
func (v View) GroupN(f GroupNFunc) GroupView {
	t := &group_transformation{
		id:       v.new_id() + ":Group",
		upstream: v.current,
		f:        f,
	}
	return GroupView{v.push(t), t}
}

type group_transformation struct {
	id         string
	upstream   transformation
	downstream []transformation // each individual group
	sidestream []transformation // the collection of groups
	f          GroupNFunc
}

func (t *group_transformation) Id() string {
	return t.id
}

func (t *group_transformation) Chain() []transformation {
	if t.upstream == nil {
		return []transformation{t}
	}
	return append(t.upstream.Chain(), t)
}

func (t *group_transformation) Dependencies() []transformation {
	if t.upstream == nil {
		return []transformation{}
	}
	return append(t.upstream.Dependencies(), t.upstream)
}

func (t *group_transformation) PushDownstream(d transformation) {
	t.downstream = append(t.downstream, d)
}

func (t *group_transformation) Transform(upstream upstream_state, txn *transaction) {
	var (
		group_state = upstream.NewState(t.id, "collection")
		group_info  = &group_transformation_collection_state{}
	)

	group_info.upstream = upstream
	txn.Restore(group_state, &group_info)
	group_state.Info = group_info

	if group_info.Groups == nil {
		group_info.Groups = make(map[string]*group_state_t)
	}

	member_states := map[string]*state{}

	{
		groups := make(map[string]*group_state_t)
		for key_str, group := range group_info.Groups {
			groups[key_str] = group
		}

		// build the reverse index
		reverse_index := map[string]map[string]byte{}
		for key_str, group := range groups {
			for _, id := range group.Members {
				m, p := reverse_index[id]
				if !p {
					reverse_index[id] = map[string]byte{}
				}

				m[key_str] = 0 // unchanged
			}
		}

		for _, id := range upstream.Removed() {
			m := reverse_index[id]
			for key_str := range m {
				m[key_str] = 3 // removed
			}
		}

		for _, id := range upstream.Added() {
			var (
				val      = upstream.Get(id)
				key_vals = t.f(Context{Id: id}, val)
				m        = map[string]byte{}
			)

			for _, key_val := range key_vals {
				key_str := CompairString(key_val)

				// add group if missing
				if _, p := groups[key_str]; !p {
					groups[key_str] = &group_state_t{KeyValue: key_val}
				}

				// add to group
				m[key_str] = 1 // added
			}

			reverse_index[id] = m
		}

		for _, id := range upstream.Changed() {
			var (
				val      = upstream.Get(id)
				key_vals = t.f(Context{Id: id}, val)
				m        = reverse_index[id]
			)

			// remove from groups
			for key_str := range m {
				m[key_str] = 3 // removed
			}

			for _, key_val := range key_vals {
				key_str := CompairString(key_val)

				// add group if missing
				if _, p := groups[key_str]; !p {
					groups[key_str] = &group_state_t{KeyValue: key_val}
				}

				// add to group
				if _, p := m[key_str]; p {
					m[key_str] = 2 // changed
				} else {
					m[key_str] = 1 // added
				}
			}

			reverse_index[id] = m
		}

		for _, group := range groups {
			group.Members = []string{}
		}

		for _, id := range upstream.Ids() {
			m := reverse_index[id]
			for key_str, s := range m {
				switch s {
				case 0: // unchanged
					group := groups[key_str]
					group.Members = append(group.Members, id)

				case 1: // added
					group := groups[key_str]
					group.Members = append(group.Members, id)

					member_state, p := member_states[key_str]
					if !p {
						member_state = upstream.NewState(t.id, "members", key_str)
						member_state.Info = &group_transformation_member_state{upstream: upstream}
						member_states[key_str] = member_state
					}
					member_state.added = append(member_state.added, id)

				case 2: // changed
					group := groups[key_str]
					group.Members = append(group.Members, id)

					member_state, p := member_states[key_str]
					if !p {
						member_state = upstream.NewState(t.id, "members", key_str)
						member_state.Info = &group_transformation_member_state{upstream: upstream}
						member_states[key_str] = member_state
					}
					member_state.changed = append(member_state.changed, id)

				case 3: // removed
					member_state, p := member_states[key_str]
					if !p {
						member_state = upstream.NewState(t.id, "members", key_str)
						member_state.Info = &group_transformation_member_state{upstream: upstream}
						member_states[key_str] = member_state
					}
					member_state.removed = append(member_state.removed, id)

				}
			}
		}

		for key_str, group := range groups {
			// remove empty groups
			if len(group.Members) == 0 {
				delete(groups, key_str)
				group_state.removed = append(group_state.removed, key_str)
				continue
			}

			// added group
			if _, p := group_info.Groups[key_str]; !p {
				group_state.added = append(group_state.added, key_str)
				continue
			}

			member_state := member_states[key_str]
			if len(member_state.added) > 0 || len(member_state.changed) > 0 || len(member_state.removed) > 0 {
				group_state.changed = append(group_state.changed, key_str)
			}
		}

		sorted := make([]string, 0, len(groups))
		for key_str := range groups {
			sorted = append(sorted, key_str)
		}
		sort.Strings(sorted)

		group_info.Sorted = sorted
		group_info.Groups = groups
	}

	txn.Save(group_state)
	txn.Propagate(t.sidestream, group_state)

	for _, member_state := range member_states {
		txn.Propagate(t.downstream, member_state)
	}
}

type group_state_t struct {
	KeyValue Value
	Members  []string
}

type group_transformation_collection_state struct {
	upstream upstream_state

	Sorted []string
	Groups map[string]*group_state_t
}

func (s *group_transformation_collection_state) Ids() []string {
	return s.Sorted
}

func (s *group_transformation_collection_state) Get(id string) Value {
	if group, p := s.Groups[id]; p {
		return Object{
			"key":        group.KeyValue,
			"member_ids": group.Members,
		}
	}
	return nil
}

type group_transformation_member_state struct {
	upstream upstream_state
	Memebers []string
}

func (s *group_transformation_member_state) Ids() []string {
	return s.Memebers
}

func (s *group_transformation_member_state) Get(id string) Value {
	return s.upstream.Get(id)
}

type group_sidestream struct {
	transformation *group_transformation
}

func (t *group_sidestream) Id() string {
	return t.transformation.Id()
}

func (t *group_sidestream) Chain() []transformation {
	return t.transformation.Chain()
}

func (t *group_sidestream) Dependencies() []transformation {
	return t.transformation.Dependencies()
}

func (t *group_sidestream) PushDownstream(downstream transformation) {
	t.transformation.sidestream = append(t.transformation.sidestream, downstream)
}

func (t *group_sidestream) Transform(upstream upstream_state, txn *transaction) {
	panic("group_sidestream cannot transform")
}
