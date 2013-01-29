package runtime

import (
	"github.com/fd/simplex/data/storage"
	"reflect"
)

type Context struct {
	txn *Transaction
}

func (ctx *Context) Load(sha SHA, val interface{}) {
	if !ctx.txn.env.store.Get(storage.SHA(sha), val) {
		panic("corrupted data store")
	}
}
func (ctx *Context) LoadValue(sha SHA, val reflect.Value) {
	if !ctx.txn.env.store.GetValue(storage.SHA(sha), val) {
		panic("corrupted data store")
	}
}

func (ctx *Context) Save(val interface{}) SHA {
	return SHA(ctx.txn.env.store.Set(val))
}
