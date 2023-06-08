package eerrors

import (
	"errors"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Error defines an grpc error that can be transformed between micro-service caller and micro-service callee.
type Error interface {
	error
	WithMetadata(map[string]string) Error
	WithMd(map[string]string) Error
	WithMessage(string) Error
	WithMsg(string) Error
}

const (
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

var _ Error = &GoError{}

type errKey string

var errs = map[errKey]*GoError{}

// Register registers error instance.
func Register(goError *GoError) {
	errs[errKey(goError.Reason)] = goError
}

// Error returns error detail message.
func (x *GoError) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v", x.Code, x.Reason, x.Message, x.Metadata)
}

// Is will be called in errors.Is method to check error type.
func (x *GoError) Is(err error) bool {
	egoErr, flag := err.(*GoError)
	if !flag {
		return false
	}
	if x == nil {
		return x == egoErr
	}
	if egoErr == nil {
		return x.Reason == ""
	}
	return x.Reason == egoErr.Reason
}

// GRPCStatus returns the Status represented by se.
func (x *GoError) GRPCStatus() *status.Status {
	s, _ := status.New(codes.Code(x.Code), x.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   x.Reason,
			Metadata: x.Metadata,
		})
	return s
}

// WithMetadata with an MD formed by the mapping of key, value.
// Deprecated: Will be removed in future versions, use WithMd instead.
func (x *GoError) WithMetadata(md map[string]string) Error {
	err := proto.Clone(x).(*GoError)
	err.Metadata = md
	return err
}

// WithMd with an MD formed by the mapping of key, value.
func (x *GoError) WithMd(md map[string]string) Error {
	err := proto.Clone(x).(*GoError)
	err.Metadata = md
	return err
}

// WithMessage set message to current GoError
// Deprecated: Will be removed in future versions, use WithMsg instead.
func (x *GoError) WithMessage(msg string) Error {
	err := proto.Clone(x).(*GoError)
	err.Message = msg
	return err
}

// WithMsg set message to current GoError
func (x *GoError) WithMsg(msg string) Error {
	err := proto.Clone(x).(*GoError)
	err.Message = msg
	return err
}

// New returns an error object for the code, message.
func New(code int, reason, message string) *GoError {
	return &GoError{
		Code:    int32(code),
		Message: message,
		Reason:  reason,
	}
}

// ToHTTPStatusCode Get equivalent HTTP status code from x.Code
func (x *GoError) ToHTTPStatusCode() int {
	return GrpcToHTTPStatusCode(codes.Code(x.Code))
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *GoError {
	if err == nil {
		return nil
	}
	if se := new(GoError); errors.As(err, &se) {
		return se
	}

	gs, ok := status.FromError(err)
	if !ok {
		return New(int(codes.Unknown), UnknownReason, err.Error())
	}

	ret := New(int(gs.Code()), UnknownReason, gs.Message())
	for _, detail := range gs.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			e, ok := errs[errKey(d.Reason)]
			if ok {
				return e.WithMsg(gs.Message()).WithMetadata(d.Metadata).(*GoError)
			}

			ret.Reason = d.Reason
			return ret.WithMd(d.Metadata).(*GoError)
		}
	}

	return ret
}
