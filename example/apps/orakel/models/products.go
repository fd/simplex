package models

import (
	"github.com/fd/w/data"
)

var Products = data.Select(of_type("product")).Sort(by_property("name"))
