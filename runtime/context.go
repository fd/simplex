package runtime

import (
	"github.com/fd/simplex/data/storage"
)

type Context struct {
	txn *Transaction
}

func (ctx *Context) Load(sha SHA, val interface{}) {
	found := ctx.txn.env.store.Get(storage.SHA(sha), val)
	if !found {
		panic("corrupted data store")
	}
}
