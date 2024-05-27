package testutils

import (
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func DiffProto(expected, actual proto.Message) string {
	return cmp.Diff(expected, actual, protocmp.Transform())
}

func EqualProtos[T proto.Message, S ~[]T](expected, actual S) []string {
	out := make([]string, 0, len(actual))
	for i, v := range actual {
		var exp T
		if i < len(expected) {
			exp = expected[i]
		}
		out = append(out, DiffProto(exp, v))
	}
	return out
}

func AllEmpty[T comparable](s *suite.Suite, t []T) bool {
	zeroSl := make([]T, 0, len(t))
	for range t {
		var zero T
		zeroSl = append(zeroSl, zero)
	}

	return s.Equal(zeroSl, t)
}
