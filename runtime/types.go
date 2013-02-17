package runtime

import (
	"reflect"
)

type (
	Resolver interface {
		DeferredId() string
		Resolve(*Transaction) IChange
	}

	Terminal interface {
		Resolver
	}

	Table interface {
		Resolver

		TableId() string
		KeyType() reflect.Type
		EltType() reflect.Type
	}

	KeyedView interface {
		Resolver

		KeyType() reflect.Type
		EltType() reflect.Type
	}

	IndexedView interface {
		Resolver

		EltType() reflect.Type
	}
)
