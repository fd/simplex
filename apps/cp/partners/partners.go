package partners

import (
	"github.com/fd/w/data"
)

var All = data.Select(of_type("partner")).Sort(by_property("name"))

func of_type(type_name string) data.SelectFunc {
	return func(ctx data.Context, val data.Value) bool {
		object, ok := val.(data.Object)
		if !ok {
			return false
		}

		return object["_type"] == type_name
	}
}

func by_property(name string) data.SortFunc {
	return func(ctx data.Context, val data.Value) data.Value {
		object, ok := val.(data.Object)
		if !ok {
			return nil
		}

		return object[name]
	}
}
