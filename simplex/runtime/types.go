package runtime

import (
	"fmt"
	"runtime"
)

type Undefined struct {
	Position string
	Reason   string
}

type (
	NilType     byte // always 0
	BooleanType bool
	IntType     int
	FloatType   float64
	StringType  string
	ArrayType   []Any
	ObjectType  map[String]Any
)

const (
	Nil NilType = 0
)

const (
	False BooleanType = BooleanType(false)
	True              = BooleanType(true)
)

func NewUndefined(skip int, format string, args ...interface{}) Undefined {
	undef := Undefined{
		Position: "???",
		Reason:   fmt.Sprintf(format, args...),
	}

	_, file, line, ok := runtime.Caller(skip + 1)
	if ok {
		undef.Position = fmt.Sprintf("%s:%d", file, line)
	}

	return undef
}

type Any interface {
	any_type() Any
}

func (v Undefined) any_type() Any   { return v }
func (v NilType) any_type() Any     { return v }
func (v BooleanType) any_type() Any { return v }
func (v IntType) any_type() Any     { return v }
func (v FloatType) any_type() Any   { return v }
func (v StringType) any_type() Any  { return v }
func (v ArrayType) any_type() Any   { return v }
func (v ObjectType) any_type() Any  { return v }

type Boolean interface {
	Any
	boolean_type() Boolean
}

func (v Undefined) boolean_type() Boolean   { return v }
func (v NilType) boolean_type() Boolean     { return False }
func (v BooleanType) boolean_type() Boolean { return v }
func (v IntType) boolean_type() Boolean     { return True }
func (v FloatType) boolean_type() Boolean   { return True }
func (v StringType) boolean_type() Boolean  { return True }
func (v ArrayType) boolean_type() Boolean   { return True }
func (v ObjectType) boolean_type() Boolean  { return True }

type Int interface {
	Any
	int_type() Int
}

func (v Undefined) int_type() Int { return v }
func (v IntType) int_type() Int   { return v }
func (v FloatType) int_type() Int { return IntType(float64(v)) }

type Float interface {
	Any
	float_type() Float
}

func (v Undefined) float_type() Float { return v }
func (v IntType) float_type() Float   { return FloatType(int(v)) }
func (v FloatType) float_type() Float { return v }

type Array interface {
	Any
	array_type() Array
}

func (v Undefined) array_type() Array { return v }
func (v ArrayType) array_type() Array { return v }

type Object interface {
	Any
	object_type() Object
}

func (v Undefined) object_type() Object  { return v }
func (v ObjectType) object_type() Object { return v }
