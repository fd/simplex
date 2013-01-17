package models

import (
	"github.com/fd/w/data"
)

var Products = data.Where(of_type("product")).Sort(by_property("name"))
