package grpcutils

import "google.golang.org/protobuf/proto"

type ProtoWrapper[T proto.Message] interface {
	ToProto() T
}
