package transformations

import (
	"fmt"
	"github.com/fd/simplex/w/data/value"
)

type (
	ScalarMapFunc    func(e Emiter, key, value value.Any)
	ScalarReduceFunc func(e Emiter, key value.Any, values []value.Any)
	ScalarMergeFunc  func(e Emiter, key value.Any, left, right []value.Any)
)

func IdentifyReduceFunc(e Emiter, key value.Any, values []value.Any) {
	if len(values) == 1 {
		e.Emit(key, values[0])
	} else if len(values) > 1 {
		panic(fmt.Sprintf("IdentifyReduceFunc: cannot be applied to multi-value keys (%v)", key))
	}
}
