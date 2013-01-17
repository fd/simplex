package value

import (
	"fmt"
)

func Get(base interface{}, path string) interface{} {
	return fmt.Sprintf("value.Get(%+v, %s)", base, path)
}
