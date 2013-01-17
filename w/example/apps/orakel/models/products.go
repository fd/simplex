package models

import (
	"github.com/fd/simplex/w/data"
)

var Products = data.Where(of_type("product")).Sort(by_property("name"))
