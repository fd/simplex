package transaction

import (
	"github.com/fd/simplex/w/data/ident"
)

type Action struct {
	// the source collection
	SourceSHA ident.SHA

	// the previously calculated collection
	PreviousSHA ident.SHA

	// SHA(SHA(key) + SHA(value)) of added members (in the source)
	Added []ident.SHA

	// keys of removed members (in the source)
	Removed []ident.SHA
}
