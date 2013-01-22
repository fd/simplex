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
	Resolve(txn *Transaction, events chan<- Event)
}

type KeyedView interface {
	KeyType() reflect.Type
	EltType() reflect.Type
	Resolve(txn *Transaction, events chan<- Event)
}

type IndexedView interface {
	EltType() reflect.Type
	Resolve(txn *Transaction, events chan<- Event)
}

func Dump(v Deferred) {
}
