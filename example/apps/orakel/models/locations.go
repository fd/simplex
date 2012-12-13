package models

import (
	"github.com/fd/w/data"
)

var Locations = data.Select(of_type("location")).Sort(by_property("name"))
var ByEvent = Locations.GroupN(by_event)

func by_event(ctx data.Context, val data.Value) []data.Value {
	object, ok := val.(data.Object)
	if !ok {
		return []data.Value{}
	}

	ids_val, ok := object["event_ids"]
	if !ok {
		return []data.Value{}
	}

	ids, ok := ids_val.([]data.Value)
	if !ok {
		return []data.Value{}
	}

	return ids
}

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
