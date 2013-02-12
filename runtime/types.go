package runtime

import (
	"reflect"
	"simplex.sh/runtime/promise"
)

type (
	Terminal interface {
		promise.Deferred
	}

	Table interface {
		promise.Deferred

		TableId() string
		KeyType() reflect.Type
		EltType() reflect.Type
	}

	KeyedView interface {
		promise.Deferred

		KeyType() reflect.Type
		EltType() reflect.Type
	}

	IndexedView interface {
		promise.Deferred

		EltType() reflect.Type
	}
)
