package data

import (
	"bytes"
	"encoding/binary"
)

func Compair(a, b Value) int {
	cmp_a := CompairString(a)
	cmp_b := CompairString(b)
	if cmp_a < cmp_b {
		return 1
	}
	if cmp_a > cmp_b {
		return -1
	}
	return 0
}

func CompairString(v Value) string {
	l := compair_string_len(v)
	buf := bytes.NewBuffer(make([]byte, 0, l))
	write_compair_string(v, buf)
	return string(buf.Bytes())
}

func write_compair_string(v Value, buf *bytes.Buffer) {
	if v == nil {
		buf.WriteByte(0)
		buf.WriteByte(0)
		return
	}

	switch a := v.(type) {
	case bool:
		if a == false {
			buf.WriteByte(0)
			buf.WriteByte(1)
		} else {
			buf.WriteByte(0)
			buf.WriteByte(2)
		}

	case int:
		buf.WriteByte(0)
		buf.WriteByte(3)
		binary.Write(buf, binary.BigEndian, float64(a))

	case float64:
		buf.WriteByte(0)
		buf.WriteByte(3)
		binary.Write(buf, binary.BigEndian, a)

	case string:
		buf.WriteByte(0)
		buf.WriteByte(4)
		buf.WriteString(a)

	case []byte:
		buf.WriteByte(0)
		buf.WriteByte(4)
		buf.Write(a)

	case []interface{}:
		buf.WriteByte(0)
		buf.WriteByte(5)

		for _, v := range a {
			write_compair_string(v, buf)
		}
	case []Value:
		buf.WriteByte(0)
		buf.WriteByte(5)

		for _, v := range a {
			write_compair_string(v, buf)
		}

	case map[string]interface{}:
		buf.WriteByte(0)
		buf.WriteByte(6)

		ks := make([]string, 0, len(a))
		for k := range a {
			ks = append(ks, k)
		}
		for _, k := range ks {
			write_compair_string(k, buf)
		}
		for _, k := range ks {
			write_compair_string(a[k], buf)
		}

	case map[string]Value:
		buf.WriteByte(0)
		buf.WriteByte(6)

		ks := make([]string, 0, len(a))
		for k := range a {
			ks = append(ks, k)
		}
		for _, k := range ks {
			write_compair_string(k, buf)
		}
		for _, k := range ks {
			write_compair_string(a[k], buf)
		}

	default:
		panic("Uncompairable type")

	}
}

func compair_string_len(v Value) int {
	if v == nil {
		return 2
	}

	switch a := v.(type) {
	case bool:
		return 2
	case int:
		return 10
	case float64:
		return 10
	case string:
		return 2 + len(a)
	case []byte:
		return 2 + len(a)
	case []Value:
		c := 2
		for _, v := range a {
			c += compair_string_len(v)
		}
		return c
	case []interface{}:
		c := 2
		for _, v := range a {
			c += compair_string_len(v)
		}
		return c
	case map[string]interface{}:
		c := 2
		for k, v := range a {
			c += compair_string_len(k)
			c += compair_string_len(v)
		}
		return c
	case map[string]Value:
		c := 2
		for k, v := range a {
			c += compair_string_len(k)
			c += compair_string_len(v)
		}
		return c
	}

	panic("Uncompairable type")
}
