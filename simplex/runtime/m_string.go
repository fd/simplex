package runtime

import (
	"strconv"
	"strings"
)

type String interface {
	Any
	string_type() String
	S_contains(sub String) Boolean
}

func (v Undefined) string_type() String  { return v }
func (v NilType) string_type() String    { return StringType("nil") }
func (v IntType) string_type() String    { return StringType(strconv.Itoa(int(v))) }
func (v StringType) string_type() String { return v }

func (v BooleanType) string_type() String {
	if v == True {
		return StringType("true")
	}
	return StringType("false")
}

func (v FloatType) string_type() String {
	return StringType(strconv.FormatFloat(float64(v), 'f', 6, 64))
}

func (str StringType) S_contains(sub String) Boolean {
	if u, ok := sub.(Undefined); ok {
		return u
	}

	go_sub, ok := sub.string_type().(StringType)
	if !ok {
		return NewUndefined(1, "Unable to cast as string %v", sub)
	}

	return BooleanType(strings.Contains(string(str), string(go_sub)))
}
func (u Undefined) S_contains(sub String) Boolean {
	return u
}

func M_String_strip(str String) String {
	if u, ok := str.(Undefined); ok {
		return u
	}

	go_str, ok := str.string_type().(StringType)
	if !ok {
		return NewUndefined(1, "Unable to cast as string %v", str)
	}

	return StringType(strings.TrimSpace(string(go_str)))
}

func M_String_lower(str String) String {
	if u, ok := str.(Undefined); ok {
		return u
	}

	go_str, ok := str.string_type().(StringType)
	if !ok {
		return NewUndefined(1, "Unable to cast as string %v", str)
	}

	return StringType(strings.ToLower(string(go_str)))
}

func M_String_upper(str String) String {
	if u, ok := str.(Undefined); ok {
		return u
	}

	go_str, ok := str.string_type().(StringType)
	if !ok {
		return NewUndefined(1, "Unable to cast as string %v", str)
	}

	return StringType(strings.ToUpper(string(go_str)))
}

func M_String_title(str String) String {
	if u, ok := str.(Undefined); ok {
		return u
	}

	go_str, ok := str.string_type().(StringType)
	if !ok {
		return NewUndefined(1, "Unable to cast as string %v", str)
	}

	return StringType(strings.ToTitle(string(go_str)))
}

func M_String_split(str, del String) Array {
	return M_String_splitN(str, del, IntType(-1))
}

func M_String_splitN(str, del String, cnt Int) Array {
	if u, ok := str.(Undefined); ok {
		return u
	}

	if u, ok := del.(Undefined); ok {
		return u
	}

	if u, ok := cnt.(Undefined); ok {
		return u
	}

	var (
		go_str StringType
		go_del StringType
		go_cnt IntType
		ok     bool
	)

	go_str, ok = str.string_type().(StringType)
	if !ok {
		return NewUndefined(1, "Unable to cast as string %v", str)
	}

	go_del, ok = del.string_type().(StringType)
	if !ok {
		return NewUndefined(1, "Unable to cast as string %v", del)
	}

	go_cnt, ok = cnt.int_type().(IntType)
	if !ok {
		return NewUndefined(1, "Unable to cast as int %v", cnt)
	}

	parts := strings.SplitN(string(go_str), string(go_del), int(go_cnt))
	arr := make(ArrayType, len(parts))
	for i, part := range parts {
		arr[i] = StringType(part)
	}

	return arr
}
