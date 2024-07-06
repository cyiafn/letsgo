package gslice

func Map[T, U any](s []T, f func(T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

func Filter[T any](s []T, f func(T) bool) []T {
	r := make([]T, 0)
	for _, v := range s {
		v := v
		if f(v) {
			r = append(r, v)
		}
	}

	return r
}

