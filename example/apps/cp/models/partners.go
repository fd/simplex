package models

import (
	"github.com/fd/w/data"
)

var Partners = data.Select(of_type("partner")).Sort(by_property("name"))
