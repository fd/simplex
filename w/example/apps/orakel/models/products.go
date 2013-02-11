package models

import (
	"simplex.sh/w/data"
)

var Products = data.Where(of_type("product")).Sort(by_property("name"))
