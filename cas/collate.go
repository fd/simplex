package cas

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"sort"
	"unicode"
)

func Collate(v interface{}) []byte {
	var buf bytes.Buffer
	write_consistent_rep(&buf, v)
	return buf.Bytes()
}

func write_consistent_rep(w *bytes.Buffer, value interface{}) {
	switch v := value.(type) {

	case nil:
		w.WriteByte(1)

	case bool:
		if v == false {
			w.WriteByte(2)
		} else {
			w.WriteByte(3)
		}

	case int8:
		w.WriteByte(4)
		binary.Write(w, binary.BigEndian, float64(v))

	case int16:
		w.WriteByte(4)
		binary.Write(w, binary.BigEndian, float64(v))

	case int32:
		w.WriteByte(4)
		binary.Write(w, binary.BigEndian, float64(v))

	case int64:
		w.WriteByte(4)
		binary.Write(w, binary.BigEndian, float64(v))

	case int:
		w.WriteByte(4)
		binary.Write(w, binary.BigEndian, float64(v))

	case uint8:
		w.WriteByte(4)
		binary.Write(w, binary.BigEndian, float64(v))

	case uint16:
		w.WriteByte(4)
		binary.Write(w, binary.BigEndian, float64(v))

	case uint32:
		w.WriteByte(4)
		binary.Write(w, binary.BigEndian, float64(v))

	case uint64:
		w.WriteByte(4)
		binary.Write(w, binary.BigEndian, float64(v))

	case uint:
		w.WriteByte(4)
		binary.Write(w, binary.BigEndian, float64(v))

	case uintptr:
		w.WriteByte(4)
		binary.Write(w, binary.BigEndian, float64(v))

	case float32:
		w.WriteByte(4)
		binary.Write(w, binary.BigEndian, float64(v))

	case float64:
		w.WriteByte(4)
		binary.Write(w, binary.BigEndian, float64(v))

	case complex64:
		w.WriteByte(5)
		binary.Write(w, binary.BigEndian, complex128(v))

	case complex128:
		w.WriteByte(5)
		binary.Write(w, binary.BigEndian, complex128(v))

	// TODO(fd) use "exp/locale/collate".(*Collator).KeyFromString(buf, v)
	case string:
		w.WriteByte(6)
		w.WriteString(v)
		w.WriteByte(0) // terminator

	// TODO(fd) use "exp/locale/collate".(*Collator).Key(buf, v)
	case []byte:
		w.WriteByte(6)
		w.Write(v)
		w.WriteByte(0) // terminator

	// TODO(fd) use "exp/locale/collate".(*Collator).Key(buf, v)
	case []rune:
		w.WriteByte(6)
		for _, r := range v {
			w.WriteRune(r)
		}
		w.WriteByte(0) // terminator

	case reflect.Value:
		write_consistent_rep_reflect(w, v)

	default:
		write_consistent_rep_reflect(w, reflect.ValueOf(v))

	}
}

func write_consistent_rep_reflect(w *bytes.Buffer, v reflect.Value) {
	typ := v.Type()
	switch typ.Kind() {

	case reflect.Bool:
		write_consistent_rep(w, v.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		write_consistent_rep(w, v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		write_consistent_rep(w, v.Uint())

	case reflect.Float32, reflect.Float64:
		write_consistent_rep(w, v.Float())

	case reflect.Array, reflect.Slice:
		w.WriteByte(7)
		l := v.Len()
		for i := 0; i < l; i++ {
			write_consistent_rep_reflect(w, v.Index(i))
		}
		w.WriteByte(0) // terminator

		//case reflect.Chan:
		//case reflect.Func:

	case reflect.Interface, reflect.Ptr:
		write_consistent_rep_reflect(w, v.Elem())

	case reflect.Map:
		keys := v.MapKeys()
		pairs := make([]map_kv_pair, len(keys))
		for i, k := range keys {
			pairs[i].cmp = Collate(k)
			pairs[i].val = v.MapIndex(k)
		}
		sort.Sort(map_kv_pairs(pairs))

		w.WriteByte(8)
		for _, pair := range pairs {
			w.Write(pair.cmp)
		}
		w.WriteByte(0) // terminator
		for _, pair := range pairs {
			write_consistent_rep_reflect(w, pair.val)
		}
		w.WriteByte(0) // terminator

	case reflect.String:
		write_consistent_rep(w, v.Bytes())

	case reflect.Struct:
		l := v.NumField()
		typ := v.Type()

		pairs := make([]map_kv_pair, 0, l)

		for i := 0; i < l; i++ {
			field := v.Field(i)
			field_typ := typ.Field(i)
			if len(field_typ.Name) > 0 && unicode.IsUpper(rune(field_typ.Name[0])) {
				pairs = append(pairs, map_kv_pair{
					cmp: Collate(field_typ.Name),
					val: field,
				})
			}
		}

		sort.Sort(map_kv_pairs(pairs))

		w.WriteByte(8)
		for _, pair := range pairs {
			w.Write(pair.cmp)
		}
		w.WriteByte(0) // terminator
		for _, pair := range pairs {
			write_consistent_rep_reflect(w, pair.val)
		}
		w.WriteByte(0) // terminator

		//case reflect.UnsafePointer:
	default:
		panic("inavlid compair type")

	}
}

type map_kv_pair struct {
	val reflect.Value
	cmp []byte
}

type map_kv_pairs []map_kv_pair

func (s map_kv_pairs) Len() int { return len(s) }
func (s map_kv_pairs) Less(i, j int) bool {
	return bytes.Compare(s[i].cmp, s[j].cmp) == -1
}
func (s map_kv_pairs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
