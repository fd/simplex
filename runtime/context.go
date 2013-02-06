package runtime

import (
	"github.com/fd/simplex/cas"
	"reflect"
)

type Context struct {
	txn *Transaction
}

func (ctx *Context) Load(addr cas.Addr, val interface{}) {
	if err := cas.Decode(ctx.txn.env.Store, addr, val); err != nil {
		panic("cas: " + err.Error())
	}
}
func (ctx *Context) LoadValue(addr cas.Addr, val reflect.Value) {
	if err := cas.DecodeValue(ctx.txn.env.Store, addr, val); err != nil {
		panic("cas: " + err.Error())
	}
}

func (ctx *Context) Save(val interface{}) cas.Addr {
	addr, err := cas.Encode(ctx.txn.env.Store, val, -1)
	if err != nil {
		panic("cas: " + err.Error())
	}
	return addr
}
