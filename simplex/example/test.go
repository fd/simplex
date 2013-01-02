//line example.smplx:2
package main

import "github.com/fd/w/simplex/runtime"
import "fmt"

//line example.smplx:5
//line example.smplx:7

//line example.smplx:6
func Lower(s runtime.String) runtime.String {
	return runtime.NewUndefined(0, "undefined")
}

//line example.smplx:11

//line example.smplx:10
func Add5(i runtime.Int) runtime.Any {
	return runtime.BINOP_ADD(i, runtime.IntType(5))
}

func main() {
	fmt.Printf("4.add5 => %v\n", Add5(runtime.IntType(4)))
	fmt.Printf("`hello`.Lower() => %v\n", Lower(runtime.StringType("hello")))
}
