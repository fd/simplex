package runtime

import (
	"simplex.sh/cas"
	"simplex.sh/cas/btree"
)

func GetTable(store cas.Store, addr cas.Addr) *btree.Tree {
	tree, err := btree.Open(store, addr)
	if err != nil {
		panic("runtime: " + err.Error())
	}
	return tree
}
