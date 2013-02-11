package runtime

import (
	"reflect"
	"simplex.sh/runtime/event"
)

type Deferred interface {
	DeferredId() string
	Resolve(txn *Transaction, events chan<- event.Event)
}

type Table interface {
	TableId() string
	KeyType() reflect.Type
	EltType() reflect.Type
	DeferredId() string
	Resolve(txn *Transaction, events chan<- event.Event)
}

type KeyedView interface {
	KeyType() reflect.Type
	EltType() reflect.Type
	DeferredId() string
	Resolve(txn *Transaction, events chan<- event.Event)
}

type IndexedView interface {
	EltType() reflect.Type
	DeferredId() string
	Resolve(txn *Transaction, events chan<- event.Event)
}
