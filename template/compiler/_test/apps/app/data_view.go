package app

import (
	"github.com/fd/w/data"
)

func AllPosts() data.View {
	return data.Window(10, 0)
}
