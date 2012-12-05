package data

type Value interface {

	// -1 => less than other
	//  0 => equal to other
	//  1 => greater than other
	Compare(other Value) int
}
