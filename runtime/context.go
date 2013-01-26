package runtime

import (
	"github.com/fd/simplex/data/storage"
)

type Context struct {
	txn *Transaction
}

func (ctx *Context) Load(sha SHA, val interface{}) {
	if !ctx.txn.env.store.Get(storage.SHA(sha), val) {
		panic("corrupted data store")
	}
}
