package digest

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"fmt"
	"hash"
	"reflect"
	"simplex.sh/store/collate"
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

func Digest(v interface{}) ([]byte, error) {
	return New().Digest(v)
}

type Digestor struct {
	hash hash.Hash
	enc  *gob.Encoder
}

func New() *Digestor {
	hash := sha1.New()
	return &Digestor{hash: hash, enc: gob.NewEncoder(hash)}
}

func (digestor *Digestor) Digest(v interface{}) ([]byte, error) {
	digestor.hash.Reset()

	err := digestor.digest(v)
	if err != nil {
		return nil, err
	}

	return digestor.hash.Sum(nil), nil
}

func (digestor *Digestor) digest(v interface{}) error {

	if v == nil {
		return digestor.enc.Encode(v)
	}

	switch v.(type) {

	case bool:
		return digestor.enc.Encode(v)

	case int, int8, int16, int32, int64:
		return digestor.enc.Encode(v)

	case uint, uint8, uint16, uint32, uint64:
		return digestor.enc.Encode(v)

	case float32, float64:
		return digestor.enc.Encode(v)

	case string:
		return digestor.enc.Encode(v)

	case []byte:
		return digestor.enc.Encode(v)

	default:
		return digestor.digest_value(reflect.ValueOf(v))

	}
}

func (digestor *Digestor) digest_value(v reflect.Value) error {

	switch v.Kind() {

	case reflect.Ptr:
		if v.IsNil() {
			return digestor.enc.Encode(nil)
		} else {
			return digestor.digest_value(v.Elem())
		}

	case reflect.Bool:
		return digestor.enc.EncodeValue(v)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return digestor.enc.EncodeValue(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return digestor.enc.EncodeValue(v)

	case reflect.Float32, reflect.Float64:
		return digestor.enc.EncodeValue(v)

	case reflect.String:
		return digestor.enc.EncodeValue(v)

	case reflect.Slice:
		var (
			err error
		)

		if v.Elem().Kind() == reflect.Uint8 {
			err = digestor.enc.EncodeValue(v)
			if err != nil {
				return err
			}
		}

		err = digestor.enc.Encode(v.Type().String())
		if err != nil {
			return err
		}
		err = digestor.enc.Encode(v.Len())
		if err != nil {
			return err
		}
		for i, l := 0, v.Len(); i < l; i++ {
			err = digestor.digest_value(v.Index(i))
			if err != nil {
				return err
			}
		}

		return nil

	case reflect.Map:
		var (
			err      error
			buf      collate.Buffer
			collator = collate.New("en")
			map_keys = make(map_keys_t, v.Len())
		)

		for i, k := range v.MapKeys() {
			buf.Reset()

			bytea, err := collator.KeyForValue(&buf, k)
			if err != nil {
				return err
			}

			map_keys[i].k = k
			map_keys[i].v = v.MapIndex(k)
			map_keys[i].o = make([]byte, len(bytea))
			copy(map_keys[i].o, bytea)
		}

		sort.Sort(map_keys)

		err = digestor.enc.Encode(v.Type().String())
		if err != nil {
			return err
		}

		for _, map_key := range map_keys {
			err = digestor.digest_value(map_key.k)
			if err != nil {
				return err
			}
			err = digestor.digest_value(map_key.v)
			if err != nil {
				return err
			}
		}

		return nil

	case reflect.Struct:
		var (
			t   = v.Type()
			f   reflect.StructField
			err error
		)

		err = digestor.enc.Encode(v.Type().String())
		if err != nil {
			return err
		}

		for i, l := 0, v.NumField(); i < l; i++ {
			f = t.Field(i)

			if f.PkgPath != "" {
				continue
			}

			err = digestor.enc.Encode(f.Name)
			if err != nil {
				return err
			}
			err = digestor.digest_value(v.Field(i))
			if err != nil {
				return err
			}
		}
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
	k, v reflect.Value
	o    []byte
}

type map_keys_t []map_key_t

func (s map_keys_t) Len() int           { return len(s) }
func (s map_keys_t) Less(i, j int) bool { return bytes.Compare(s[i].o, s[j].o) == -1 }
func (s map_keys_t) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
