package runtime

import (
	"reflect"
)

type Table struct {
}

type GenericTable interface {
	InnerTable() *Table
}

type GenericKeyedView interface {
	KeyType() reflect.Type
}

type GenericIndexedView interface {
	EltType() reflect.Type
}

type GenericView interface {
	InnerView() *Table
}

func Dump(v GenericView) {
}

func Select(v GenericView, f interface{}) *Table {
	return nil
}

func Collect(v GenericView, f interface{}) *Table {
	return nil
}

func Sort(v GenericView, f interface{}) *Table {
	return nil
}
