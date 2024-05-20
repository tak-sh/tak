package except

import (
	"errors"
	"fmt"
	"github.com/tak-sh/tak/generated/go/except"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"strings"
)

var ReasonString = map[except.Reason]string{
	except.Reason_UNKNOWN:        "unknown",
	except.Reason_NOT_FOUND:      "not found",
	except.Reason_INVALID:        "invalid",
	except.Reason_INTERNAL:       "internal",
	except.Reason_TIMEOUT:        "timeout",
	except.Reason_ALREADY_EXISTS: "already exists",
	except.Reason_ABORTED:        "aborted",
	except.Reason_FAILED:         "failed",
}

func New(reason except.Reason, msg string) error {
	return &Err{
		except.Exception{
			Reason:  reason,
			Message: msg,
		},
	}
}

func Newf(reason except.Reason, msg string, args ...any) error {
	return &Err{
		except.Exception{
			Reason:  reason,
			Message: fmt.Sprintf(msg, args...),
		},
	}
}

func NewNotFound(msg string, args ...any) error {
	return Newf(except.Reason_NOT_FOUND, msg, args...)
}

func NewAborted(msg string, args ...any) error {
	return Newf(except.Reason_ABORTED, msg, args...)
}

func NewInvalid(msg string, args ...any) error {
	return Newf(except.Reason_INVALID, msg, args...)
}

func NewInternal(msg string, args ...any) error {
	return Newf(except.Reason_INTERNAL, msg, args...)
}

func NewTimeout(msg string, args ...any) error {
	return Newf(except.Reason_TIMEOUT, msg, args...)
}

func NewAlreadyExists(msg string, args ...any) error {
	return Newf(except.Reason_ALREADY_EXISTS, msg, args...)
}

func NewFailed(msg string, args ...any) error {
	return Newf(except.Reason_FAILED, msg, args...)
}

func Reason(err error) except.Reason {
	var e *Err
	errors.As(err, &e)
	return e.GetReason()
}

// HasReason returns true if r matches err. err must be an Err.
func HasReason(err error, r ...except.Reason) bool {
	var e *Err
	if errors.As(err, &e) {
		for _, v := range r {
			if v == e.Reason {
				return true
			}
		}
	}
	return false
}

var _ error = &Err{}

type Err struct {
	except.Exception
}

func (e *Err) Error() string {
	out := make([]string, 0, 2)
	if r := ReasonString[e.GetReason()]; r != "" {
		out = append(out, r)
	}

	out = append(out, e.GetMessage())

	return strings.Join(out, ": ")
}

func NewFromGrpcStatus(st *status.Status) *Err {
	if st == nil {
		return nil
	}
	return &Err{
		except.Exception{
			Reason:  GrpcCodeToReason(codes.Code(st.Code)),
			Message: st.Message,
		},
	}
}

func GrpcCodeToReason(c codes.Code) except.Reason {
	switch c {
	case codes.AlreadyExists:
		return except.Reason_ALREADY_EXISTS
	case codes.InvalidArgument:
		return except.Reason_INVALID
	case codes.DeadlineExceeded:
		return except.Reason_TIMEOUT
	case codes.NotFound:
		return except.Reason_NOT_FOUND
	case codes.Aborted:
		return except.Reason_ABORTED
	default:
		return except.Reason_INTERNAL
	}
}

func ReasonToGrpcCode(r except.Reason) codes.Code {
	switch r {
	case except.Reason_UNKNOWN:
		return codes.Unknown
	case except.Reason_NOT_FOUND:
		return codes.NotFound
	case except.Reason_INVALID:
		return codes.InvalidArgument
	case except.Reason_TIMEOUT:
		return codes.DeadlineExceeded
	case except.Reason_ALREADY_EXISTS:
		return codes.AlreadyExists
	case except.Reason_ABORTED:
		return codes.Aborted
	default:
		return codes.Internal
	}
}
