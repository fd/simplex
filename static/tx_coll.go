package static

import (
	"encoding/json"
	"reflect"
)

func (tx *Tx) Coll(name string, typ interface{}) *C {
	tx.mtx.Lock()
	defer tx.mtx.Unlock()

	if tx.collections == nil {
		tx.collections = map[string]*C{}
	}

	if coll := tx.collections[name]; coll != nil {
		return coll
	}

	out := &C{elem_type: reflect.TypeOf(typ), tx: tx}
	tx.collections[name] = out

	out.t.Do(func() error {
		var (
			et          = reflect.TypeOf(typ)
			json_coll   json_coll_t
			elems       []interface{}
			err         error
			has_pointer bool
		)

		if et.Kind() == reflect.Ptr {
			has_pointer = true
			et = et.Elem()
		}

		r, err := tx.src.GetBlob(name + ".json")
		if err != nil {
			return err
		}

		defer r.Close()

		err = json.NewDecoder(r).Decode(&json_coll)
		if err != nil {
			return err
		}

		elems = make([]interface{}, len(json_coll))

		for i, json_elem := range json_coll {
			var elem interface{}
			if has_pointer {
				elem = reflect.New(et).Interface()
			} else {
				elem = reflect.Indirect(reflect.New(et)).Interface()
			}

			err := json.Unmarshal(json_elem, &elem)
			if err != nil {
				return err
			}

			elems[i] = elem
		}

		out.elems = elems
		return nil
	})

	return out
}

type json_coll_t []json.RawMessage
