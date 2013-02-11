package types

// A View represents a view type view[Key]Elt.
type View struct {
	Key, Elt Type
}

// A View represents a table type table[Key]Elt.
type Table struct {
	Key, Elt Type
}

type Viewish interface {
	KeyType() Type
	EltType() Type
}

func (v *View) KeyType() Type  { return v.Key }
func (v *Table) KeyType() Type { return v.Key }

func (v *View) EltType() Type  { return v.Elt }
func (v *Table) EltType() Type { return v.Elt }

func (*View) aType()  {}
func (*Table) aType() {}
