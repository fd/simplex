package runtime

func BINOP_SUB(left, right Any) Any {
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
				return IntType(int(l) - int(r))
			}
		}

	case FloatType:
		if r_cat, ok := right.(Float); ok {
			if r, ok := r_cat.float_type().(FloatType); ok {
				return FloatType(float64(l) - float64(r))
			}
		}

	}

	return NewUndefined(1, "Unable to add (+) %T to %T", right, left)
}
