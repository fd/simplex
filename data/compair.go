package data

func Compair(a, b Value) int {
	val_a, type_a := type_index_for_compair(a)
	val_b, type_b := type_index_for_compair(b)

	if type_a != type_b {
		return simple_compair_result(type_a, type_b)
	}

	switch type_a {
	case 1: // nil is always nil
		return 0

	case 2:
		var (
			a_a bool
			a_b bool
		)
		a_a = val_a.(bool)
		a_b = val_b.(bool)

		if a_a == a_b {
			return 0
		}
		if a_a == false {
			return 1
		}
		return -1

	case 3:
		var (
			a_a float64
			a_b float64
		)
		a_a = val_a.(float64)
		a_b = val_b.(float64)

		if a_a < a_b {
			return 1
		}
		if a_a > a_b {
			return -1
		}
		return 0

	case 4:
		var (
			a_a string
			a_b string
		)
		a_a = val_a.(string)
		a_b = val_b.(string)

		if a_a < a_b {
			return 1
		}
		if a_a > a_b {
			return -1
		}
		return 0

	case 5:
		var (
			a_a []interface{}
			a_b []interface{}
			l_a int
			l_b int
			l   int
		)
		a_a = val_a.([]interface{})
		a_b = val_b.([]interface{})
		l_a = len(a_a)
		l_b = len(a_b)

		if l_a > l_b {
			l = l_b
		} else {
			l = l_a
		}

		for i := 0; i < l; i++ {
			c := Compair(a_a[i], a_b[i])
			if c != 0 {
				return c
			}
		}

		if l_a < l_b {
			return 1
		}
		if l_a > l_b {
			return -1
		}
		return 0

	case 6:
		// TODO

	}

	panic("Uncompairable type")
}

func type_index_for_compair(v Value) (Value, int) {
	if v == nil {
		return nil, 1
	}

	switch a := v.(type) {
	case bool:
		return v, 2
	case int:
		return float64(a), 3
	case float64:
		return v, 3
	case string:
		return v, 4
	case []byte:
		return string(a), 4
	case []rune:
		return string(a), 4
	case []interface{}:
		return v, 5
	case map[string]interface{}:
		return v, 6
	}

	panic("Uncompairable type")
}

func simple_compair_result(a, b int) int {
	if a > b {
		return -1
	}
	if a < b {
		return 1
	}
	return 0
}
