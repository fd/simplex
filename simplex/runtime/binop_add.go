package runtime

/*

source(T)[]T
[]T.where(func(T)bool)[]T
[]T.map(func(T)T2)[]T2
[]T.sort(func(T)T2)[]T
[]T.group(func(T)T2)[]T
[]T.union([]T)[]T
[]T.render(func(T)Fragment)[]Fragment
[]T.route(func(T)Route)[]Route
[]T.schedule(func(T)Cron)[]Cron
[]T.count()int
[]T.sum(func(T)int)int
[]T.sum(func(T)float64)float64
[]T.avg(func(T)float64)float64

*/

/*
  Undefined + Any                          => Undefined
  Any       + Undefined                    => Undefined
  Int       + Int|Float                    => Int
  Float     + Int|Float                    => Float
  String    + Nil|Boolean|Int|Float|String => String
  Array     + Array                        => Array
  Any       + Any                          => Undefined
*/
func BINOP_ADD(left, right Any) Any {
	if _, ok := left.(Undefined); ok {
		return left
	}

	if _, ok := right.(Undefined); ok {
		return right
	}

	switch l := left.(type) {
	case IntType:
		if r_cat, ok := right.(Int); ok {
			if r, ok := r_cat.int_type().(IntType); ok {
				return IntType(int(l) + int(r))
			}
		}

	case FloatType:
		if r_cat, ok := right.(Float); ok {
			if r, ok := r_cat.float_type().(FloatType); ok {
				return FloatType(float64(l) + float64(r))
			}
		}

	case StringType:
		if r_cat, ok := right.(String); ok {
			if r, ok := r_cat.string_type().(StringType); ok {
				return StringType(string(l) + string(r))
			}
		}

	case ArrayType:
		if r_cat, ok := right.(Array); ok {
			if r, ok := r_cat.array_type().(ArrayType); ok {
				res := make(ArrayType, len(l)+len(r))
				copy(res[0:], l)
				copy(res[len(l):], r)
				return res
			}
		}

	}

	return NewUndefined(1, "Unable to add (+) %T to %T", right, left)
}
