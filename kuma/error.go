package kuma

import (
	"fmt"

	"github.com/kumarabd/gokit/errors"
)

// ErrInstallMesh is the error for install mesh
func ErrInstallMesh(err error) error {
	return errors.New("1001", fmt.Sprintf("Error installing mesh: %s", err.Error()))
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.New("1002", fmt.Sprintf("Error configuration mesh: %s", err.Error()))
}

// ErrPortForward is the error for mesh port forward
func ErrPortForward(err error) error {
	return errors.New("1003", fmt.Sprintf("Error portforwarding mesh gui: %s", err.Error()))
}
