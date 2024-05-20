package ptr

func Ptr[T any](v T) (out *T) {
	return &v
}

func Deref[T any](v *T) (out T) {
	if v != nil {
		out = *v
	}
	return
}

func PtrOrNil[T comparable](v T) *T {
	var out T
	if v == out {
		return nil
	}
	return &v
}
