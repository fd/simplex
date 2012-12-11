package data

type transaction struct {
	upstream_states []StoreReader

	added   []string
	updated []string
	removed []string
}

type Context struct {
	Id string
}
