package types

// A View represents a view type view[Key]Elt.
type View struct {
	implementsType
	Key, Elt Type
}

// A View represents a table type table[Key]Elt.
type Table struct {
	implementsType
	Key, Elt Type
}
