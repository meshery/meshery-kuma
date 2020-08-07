package grpc

import (
	"fmt"

	"github.com/kumarabd/gokit/errors"
)

// ErrPanic is the error object for panic errors
func ErrPanic(r interface{}) error {
	return errors.New("600", fmt.Sprintf("%v", r))
}

// ErrGrpcListener is the error object for grpc listener
func ErrGrpcListener(err error) error {
	return errors.New("601", fmt.Sprintf("Error during grpc listener initialization : %v", err))
}

// ErrGrpcServer is the error object for grpc server
func ErrGrpcServer(err error) error {
	return errors.New("602", fmt.Sprintf("Error during grpc server initialization : %v", err))
}
