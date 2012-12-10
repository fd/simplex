package data

/*
  transaction keeps track of the changes and dependencies
  during a single transformation.
*/
type transaction struct {
	upstream_states []StoreReader

	added   []string
	updated []string
	removed []string
}

type Context struct {
	Id string
}
