package runtime

import (
	"github.com/fd/simplex/runtime/event"
	"reflect"
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
