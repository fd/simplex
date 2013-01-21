package runtime

import (
	"reflect"
)

type Deferred interface {
	Resolve(txn *Transaction, events chan<- Event)
}

type Table interface {
	TableId() string
	KeyType() reflect.Type
	EltType() reflect.Type
	Resolve()
}

type KeyedView interface {
	KeyType() reflect.Type
	EltType() reflect.Type
	Resolve()
}

type IndexedView interface {
	EltType() reflect.Type
	Resolve()
}

func Dump(v Deferred) {
}
