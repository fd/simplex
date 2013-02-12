package runtime

import (
	"reflect"
	"simplex.sh/cas"
)

type Context struct {
	store cas.Store
}

func (ctx *Context) Load(addr cas.Addr, val interface{}) {
	if err := cas.Decode(ctx.store, addr, val); err != nil {
		panic("cas: " + err.Error())
	}
}
func (ctx *Context) LoadValue(addr cas.Addr, val reflect.Value) {
	if err := cas.DecodeValue(ctx.store, addr, val); err != nil {
		panic("cas: " + err.Error())
	}
}

func (ctx *Context) Save(val interface{}) cas.Addr {
	addr, err := cas.Encode(ctx.store, val, -1)
	if err != nil {
		panic("cas: " + err.Error())
	}
	return addr
}
