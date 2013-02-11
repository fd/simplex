package models

import (
	"simplex.sh/w/data"
)

type Location struct {
	data.Type

	Name    string
	Address string
	Website url.URL

	GPS struct {
		Coordinates struct {
			Lng float64
			Lat float64
		}

		MapUrl url.URL
	}
}

var Locations = data.Where(M._type == "location").Sort(M.name)
var ByEvent = Locations.GroupN(M.event_ids)
