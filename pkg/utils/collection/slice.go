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

func ContainsList[E comparable, S ~[]E](s S, sub S) (contains bool) {
	n := len(sub)
	if n > len(s) {
		return false
	}

	contains = true
	for i := 0; i < n && contains; i++ {
		contains = contains && s[i] == sub[i]
	}

	return
}

func Rel[E comparable, S ~[]E](base S, target S) (out S, matched bool) {
	t := len(target)
	n := len(base)
	out = make(S, 0, t)
	if matched = ContainsList(target, base); matched && t > n {
		out = append(out, target[n:]...)
	}
	return
}
