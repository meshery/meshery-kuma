package kuma

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

var (
	ErrInstallKumaCode   = "kuma_test_code"
	ErrMeshConfigCode    = "kuma_test_code"
	ErrFetchManifestCode = "kuma_test_code"
	ErrClientConfigCode  = "kuma_test_code"
	ErrClientSetCode     = "kuma_test_code"
	ErrStreamEventCode   = "kuma_test_code"
	ErrSampleAppCode     = "kuma_test_code"
	ErrGetKumactlCode    = "kuma_test_code"
	ErrUntarCode         = "kuma_test_code"
	ErrUntarDefaultCode  = "kuma_test_code"
	ErrMoveBinaryCode    = "kuma_test_code"

	ErrOpInvalid    = errors.NewDefault(errors.ErrOpInvalid, "Invalid operation")
	ErrUntarDefault = errors.NewDefault(ErrUntarDefaultCode, "Error untaring operation default")
)

// ErrInstallMesh is the error for install mesh
func ErrInstallKuma(err error) error {
	return errors.NewDefault(ErrInstallKumaCode, fmt.Sprintf("Error with kuma operation: %s", err.Error()))
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.NewDefault(ErrMeshConfigCode, fmt.Sprintf("Error configuration mesh: %s", err.Error()))
}

// ErrPortForward is the error for mesh port forward
func ErrFetchManifest(err error, des string) error {
	return errors.NewDefault(ErrFetchManifestCode, fmt.Sprintf("Error fetching mesh manifest: %s", des))
}

// ErrClientConfig is the error for setting client config
func ErrClientConfig(err error) error {
	return errors.NewDefault(ErrClientConfigCode, fmt.Sprintf("Error setting client config: %s", err.Error()))
}

// ErrPortForward is the error for setting clientset
func ErrClientSet(err error) error {
	return errors.NewDefault(ErrClientSetCode, fmt.Sprintf("Error setting clientset: %s", err.Error()))
}

// ErrStreamEvent is the error for streaming event
func ErrStreamEvent(err error) error {
	return errors.NewDefault(ErrStreamEventCode, fmt.Sprintf("Error streaming event: %s", err.Error()))
}

// ErrSampleApp is the error for streaming event
func ErrSampleApp(err error) error {
	return errors.NewDefault(ErrSampleAppCode, fmt.Sprintf("Error with sample app operation: %s", err.Error()))
}

// ErrGetKumactl is the error for streaming event
func ErrGetKumactl(err error) error {
	return errors.NewDefault(ErrGetKumactlCode, fmt.Sprintf("Error getting kumactl commandline: %s", err.Error()))
}

// ErrUntar is the error for streaming event
func ErrUntar(err error) error {
	return errors.NewDefault(ErrUntarCode, fmt.Sprintf("Error untaring package: %s", err.Error()))
}

// ErrMoveBinary is the error for streaming event
func ErrMoveBinary(err error) error {
	return errors.NewDefault(ErrMoveBinaryCode, fmt.Sprintf("Error with moving binary kumactl: %s", err.Error()))
}
