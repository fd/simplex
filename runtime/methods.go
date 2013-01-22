package runtime

func Select(v Deferred, f func(interface{}) bool) Deferred {
	return nil
}

func Reject(v Deferred, f func(interface{}) bool) Deferred {
	return nil
}

func Detect(v Deferred, f func(interface{}) bool) interface{} {
	return nil
}

func Collect(v Deferred, f func(interface{}) interface{}) Deferred {
	return nil
}

func Inject(v Deferred, f func(interface{}, []interface{}) interface{}) Deferred {
	return nil
}

func Group(v Deferred, f func(interface{}) interface{}) Deferred {
	return nil
}

func Index(v Deferred, f func(interface{}) interface{}) Deferred {
	return nil
}

func Sort(v Deferred, f func(interface{}) interface{}) Deferred {
	return nil
}

func Union(v ...Deferred) Deferred {
	return nil
}
