package data

/*
  transaction keeps track of the changes and dependencies
  during a single transformation.
*/
type transaction struct {
	added   []int
	updated []int
	removed []int
}
