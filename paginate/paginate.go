package paginate

import (
	"reflect"
	"simplex.sh/static"
)

type Page struct {
	Number    int
	Elements  []interface{}
	elem_type reflect.Type
}

func Paginate(in *static.C, size int) *static.C {
	if size <= 0 {
		size = 15
	}

	return in.Transform(reflect.TypeOf(&Page{}), func(elems []interface{}) ([]interface{}, error) {
		var (
			number = 1
			pages  = make([]interface{}, 0, len(elems)/size+1)
		)

		for l := len(elems); l > 0; l -= size {
			end := size
			if end > l {
				end = l
			}

			page := &Page{
				Number:    number,
				Elements:  elems[:end],
				elem_type: in.ElemType(),
			}
			pages = append(pages, page)

			number += 1
			elems = elems[end:]
		}

		return pages, nil
	})
}
