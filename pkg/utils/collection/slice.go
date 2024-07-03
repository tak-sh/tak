package collection

func All[T any](t []T, f func(t T) bool) bool {
	all := true
	for _, v := range t {
		all = all && f(v)
	}
	return all
}

func Any[T any](t []T, f func(t T) bool) bool {
	all := false
	for _, v := range t {
		all = all || f(v)
	}
	return all
}

func AllZero[T comparable](t []T) bool {
	var zero T
	return All(t, func(e T) bool {
		return e == zero
	})
}

func AnyZero[T comparable](t []T) bool {
	var zero T
	return Any(t, func(e T) bool {
		return e == zero
	})
}

type CastFunc[F, T any] func(f F) T

func Cast[F, T any](s []F, c CastFunc[F, T]) []T {
	out := make([]T, 0, len(s))
	for _, v := range s {
		out = append(out, c(v))
	}
	return out
}
