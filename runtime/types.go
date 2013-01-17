package runtime

import (
	"reflect"
)

type GenericTable interface {
	GenericKeyedView
	InnerTable()
}

type GenericKeyedView interface {
	GenericIndexedView
	KeyType() reflect.Type
}

type GenericIndexedView interface {
	GenericView
	EltType() reflect.Type
}

type GenericView interface {
	InnerView()
}

func Dump(v GenericView) {
}
