package runtime

import (
	"bytes"
	"simplex.sh/cas"
	"sort"
	"sync"
)

type (
	ChangeType int8

	IChange struct {
		A             cas.Addr
		B             cas.Addr
		MemberChanges []MemberChange

		Err   error
		Stack []byte

		mutex sync.Mutex
	}

	MemberChange struct {
		CollatedKey []byte
		Key         cas.Addr
		IChange
	}
)

const (
	ChangeNone ChangeType = iota
	ChangeInsert
	ChangeUpdate
	ChangeRemove
)

func (c IChange) Type() ChangeType {
	if c.A == nil && c.B == nil {
		return ChangeNone
	}

	if c.A == nil {
		return ChangeInsert
	}

	if c.B == nil {
		return ChangeRemove
	}

	if bytes.Compare(c.A, c.B) == 0 {
		return ChangeNone
	}

	return ChangeUpdate
}

func (c *IChange) MemberChanged(collated_key []byte, key_addr cas.Addr, ichange IChange) {
	change := MemberChange{
		CollatedKey: collated_key,
		Key:         key_addr,
		IChange:     ichange,
	}

	if change.Type() == ChangeNone {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	i := sort.Search(len(c.MemberChanges), func(i int) bool {
		return bytes.Compare(c.MemberChanges[i].CollatedKey, change.CollatedKey) != -1
	})

	if cap(c.MemberChanges) < len(c.MemberChanges)+1 {
		src := c.MemberChanges
		dst := make([]MemberChange, len(c.MemberChanges), cap(c.MemberChanges)+100)
		copy(dst, src)
		c.MemberChanges = dst
	}

	// is append
	if i >= len(c.MemberChanges) {
		c.MemberChanges = append(c.MemberChanges, change)
		return
	}

	if bytes.Compare(c.MemberChanges[i].CollatedKey, change.CollatedKey) == 0 {
		panic("already record change for collated_key")
	}

	// make room
	dst := c.MemberChanges[:len(c.MemberChanges)+1]
	if i > 0 {
		copy(dst, c.MemberChanges[:i])
	}
	copy(dst[i+1:], c.MemberChanges[i:])
	c.MemberChanges = dst

	// set change
	c.MemberChanges[i] = change
}
