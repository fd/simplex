package collate

import (
	"bytes"
	"code.google.com/p/go.exp/locale/collate"
	"code.google.com/p/go.exp/locale/collate/colltab"
	"encoding/binary"
	"fmt"
	"reflect"
	"sort"
)

const (
	typ_nil byte = 1 + iota
	typ_false
	typ_true
	typ_int_neg
	typ_int_pos
	typ_float
	typ_string
	typ_bytea
	typ_slice
	typ_map
	typ_struct
)

type Buffer struct {
	buf  [4096]byte
	key  []byte
	coll collate.Buffer
}

func (b *Buffer) init() {
	if b.key == nil {
		b.key = b.buf[:0]
	}
}

func (buf *Buffer) write_byte(b byte) {
	buf.key = append(buf.key, b)
}

func (buf *Buffer) ewrite_bytes(x []byte) {
	for _, b := range x {
		if b == 0x00 || b == 0x01 {
			buf.write_byte(0x01)
			buf.write_byte(b)
			continue
		}
		if b == 0xFF || b == 0xFE {
			buf.write_byte(0xFE)
			buf.write_byte(b)
			continue
		}
		buf.write_byte(b)
	}
}

func (buf *Buffer) ewrite_string(x string) {
	buf.ewrite_bytes([]byte(x))
}

func (buf *Buffer) write_bytes(x []byte) {
	buf.key = append(buf.key, x...)
}

func (buf *Buffer) write_string(x string) {
	buf.write_bytes([]byte(x))
}

func (buf *Buffer) Write(x []byte) (int, error) {
	buf.write_bytes(x)
	return 0, nil
}

func (b *Buffer) Reset() {
	b.key = b.key[:0]
	b.coll.Reset()
}

type Collator struct {
	coll *collate.Collator
}

func New(locale string) *Collator {
	c := &Collator{coll: collate.New(locale)}
	c.coll.Strength = colltab.Secondary
	c.coll.CaseLevel = true
	// c.coll.Numeric = true
	return c
}

func (collator *Collator) Key(buf *Buffer, v interface{}) ([]byte, error) {
	buf.init()
	err := collator.key(buf, v)
	if err != nil {
		return nil, err
	}
	return buf.key, nil
}

func (collator *Collator) KeyForValue(buf *Buffer, v reflect.Value) ([]byte, error) {
	buf.init()
	err := collator.key_for_value(buf, v)
	if err != nil {
		return nil, err
	}
	return buf.key, nil
}

func (collator *Collator) key(buf *Buffer, v interface{}) error {

	if v == nil {
		buf.write_byte(typ_nil)
		return nil
	}

	switch x := v.(type) {

	case bool:
		if x {
			buf.write_byte(typ_true)
		} else {
			buf.write_byte(typ_false)
		}
		return nil

	case int64:
		if x >= 0 {
			buf.write_byte(typ_int_pos)
		} else {
			buf.write_byte(typ_int_neg)
		}
		return binary.Write(buf, binary.BigEndian, uint64(x))

	case uint64:
		buf.write_byte(typ_int_pos)
		return binary.Write(buf, binary.BigEndian, x)

	case float64:
		buf.write_byte(typ_float)
		return binary.Write(buf, binary.BigEndian, x)

	case string:
		buf.write_byte(typ_bytea)
		buf.ewrite_string(x)
		buf.write_byte(0x00)
		return nil

	case []byte:
		buf.write_byte(typ_bytea)
		buf.ewrite_bytes(x)
		buf.write_byte(0x00)
		return nil

	case String:
		buf.write_byte(typ_string)
		buf.coll.Reset()
		buf.write_bytes(collator.coll.KeyFromString(&buf.coll, string(x)))
		buf.write_byte(0x00)
		return nil

	case Runea:
		buf.write_byte(typ_string)
		buf.coll.Reset()
		buf.write_bytes(collator.coll.KeyFromString(&buf.coll, string(x)))
		buf.write_byte(0x00)
		return nil

	case Bytea:
		buf.write_byte(typ_string)
		buf.coll.Reset()
		buf.write_bytes(collator.coll.Key(&buf.coll, []byte(x)))
		buf.write_byte(0x00)
		return nil

	default:
		return collator.key_for_value(buf, reflect.ValueOf(v))

	}
}

func (collator *Collator) key_for_value(buf *Buffer, v reflect.Value) error {

	switch v.Kind() {

	case reflect.Ptr:
		if v.IsNil() {
			buf.write_byte(typ_nil)
			return nil
		} else {
			return collator.key_for_value(buf, v.Elem())
		}

	case reflect.Bool:
		return collator.key(buf, v.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return collator.key(buf, v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return collator.key(buf, v.Uint())

	case reflect.Float32, reflect.Float64:
		return collator.key(buf, v.Float())

	case reflect.String:
		return collator.key(buf, v.String())

	case reflect.Slice:
		if v.Elem().Kind() == reflect.Uint8 {
			return collator.key(buf, v.Bytes())
		}

		buf.write_byte(typ_slice)
		for i, l := 0, v.Len(); i < l; i++ {
			err := collator.key_for_value(buf, v.Index(i))
			if err != nil {
				return err
			}
		}
		buf.write_byte(0x00)
		return nil

	case reflect.Map:
		var (
			err      error
			map_keys = make(map_keys_t, v.Len())
			buf_ctx  = buf.key
		)

		buf.key = make([]byte, 0, 4096)

		for i, k := range v.MapKeys() {
			err = collator.key_for_value(buf, k)
			if err != nil {
				return err
			}
			k_buf := buf.key
			buf.key = buf.key[len(k_buf):]

			err = collator.key_for_value(buf, k)
			if err != nil {
				return err
			}
			v_buf := buf.key
			buf.key = buf.key[len(v_buf):]

			map_keys[i].k = k_buf
			map_keys[i].v = v_buf
		}

		sort.Sort(map_keys)
		buf.key = buf_ctx

		buf.write_byte(typ_map)
		for _, map_key := range map_keys {
			buf.key = append(buf.key, map_key.k...)
		}
		buf.write_byte(0x00)
		for _, map_key := range map_keys {
			buf.key = append(buf.key, map_key.v...)
		}
		buf.write_byte(0x00)

		return nil

	case reflect.Struct:
		var (
			t   = v.Type()
			f   reflect.StructField
			err error
		)
		buf.write_byte(typ_struct)
		for i, l := 0, v.NumField(); i < l; i++ {
			f = t.Field(i)

			if f.PkgPath != "" {
				continue
			}

			err = collator.key(buf, f.Name)
			if err != nil {
				return err
			}
		}
		buf.write_byte(0x00)
		for i, l := 0, v.NumField(); i < l; i++ {
			f = t.Field(i)

			if f.PkgPath != "" {
				continue
			}

			err = collator.key_for_value(buf, v.Field(i))
			if err != nil {
				return err
			}
		}
		buf.write_byte(0x00)
		return nil

	default:
		panic(&UnsupportedType{v})

	}
}

type UnsupportedType struct {
	value reflect.Value
}

func (u *UnsupportedType) Error() string {
	return fmt.Sprintf("collate: unsupported type: %s", u.value.Type())
}

type map_key_t struct {
	k, v []byte
}

type map_keys_t []map_key_t

func (s map_keys_t) Len() int           { return len(s) }
func (s map_keys_t) Less(i, j int) bool { return bytes.Compare(s[i].k, s[j].k) == -1 }
func (s map_keys_t) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
