package transaction

import (
	"github.com/fd/simplex/w/data/ident"
)

type Transaction struct {
	Parent ident.SHA
	Source ident.SHA
	Bags   map[string]ident.SHA
}

func Apply(parent ident.SHA) ident.SHA {
	// get previous source
	// apply changes
	// store new source
	// run transformation meta function
	// store transcation
	// return sha
	return ident.ZeroSHA
}
