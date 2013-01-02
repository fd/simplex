package runtime

/*
  Undefined == Any       => Undefined
  Any       == Undefined => Undefined
  Any       == Any       => Boolean
*/
func BINOP_EQL(left, right Any) Boolean {
	if l, ok := left.(Undefined); ok {
		return l
	}

	if r, ok := right.(Undefined); ok {
		return r
	}

	switch l := left.(type) {
	case BooleanType:
		if r_cat, ok := right.(Boolean); ok {
			if r, ok := r_cat.boolean_type().(BooleanType); ok {
				return BooleanType(bool(l) == bool(r))
			}
		}

	case IntType:
		if r_cat, ok := right.(Int); ok {
			if r, ok := r_cat.int_type().(IntType); ok {
				return BooleanType(int(l) == int(r))
			}
		}

	case FloatType:
		if r_cat, ok := right.(Float); ok {
			if r, ok := r_cat.float_type().(FloatType); ok {
				return BooleanType(float64(l) == float64(r))
			}
		}

	case StringType:
		if r_cat, ok := right.(String); ok {
			if r, ok := r_cat.string_type().(StringType); ok {
				return BooleanType(string(l) == string(r))
			}
		}

	case ArrayType:
		if r_cat, ok := right.(Array); ok {
			if r, ok := r_cat.array_type().(ArrayType); ok {
				if len(l) != len(r) {
					return False
				}

				for i, e := range l {
					if BINOP_EQL(e, r[i]) == False {
						return False
					}
				}

				return True
			}
		}

	case ObjectType:
		if r_cat, ok := right.(Object); ok {
			if r, ok := r_cat.object_type().(ObjectType); ok {
				if len(l) != len(r) {
					return False
				}

				for k, ml := range l {
					mr, p := r[k]
					if !p || BINOP_EQL(ml, mr) == False {
						return False
					}
				}

				return True
			}
		}

	}

	return False
}
