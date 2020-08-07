package grpc

import (
	"fmt"

	"github.com/kumarabd/gokit/errors"
)

func ErrPanic(r interface{}) error {
	return errors.New("600", fmt.Sprintf("%v", r))
}

func ErrGrpcListener(err error) error {
	return errors.New("601", fmt.Sprintf("Error during grpc listener initialization : %v", err))
}

func ErrGrpcServer(err error) error {
	return errors.New("602", fmt.Sprintf("Error during grpc server initialization : %v", err))
}
