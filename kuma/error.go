package kuma

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

var (
	// Error code for failed service mesh installation

	// ErrInstallKumaCode represents the errors which are generated
	// during kuma service mesh install process
	ErrInstallKumaCode = "kuma_test_code"

	// ErrMeshConfigCode represents the errors which are generated
	// when an invalid mesh config is found
	ErrMeshConfigCode = "kuma_test_code"

	// ErrFetchManifestCode represents the errors which are generated
	// during the process of fetching manifest
	ErrFetchManifestCode = "kuma_test_code"

	// ErrClientConfigCode represents the errors which are generated
	// during the process of setting client config
	ErrClientConfigCode = "kuma_test_code"

	// ErrClientSetCode represents the errors which are generated
	// during the process of setting clientset
	ErrClientSetCode = "kuma_test_code"

	// ErrStreamEventCode represents the errors which are generated
	// during the process of streaming events
	ErrStreamEventCode = "kuma_test_code"

	// ErrSampleAppCode represents the errors which are generated
	// during the process of installing sample app
	ErrSampleAppCode = "kuma_test_code"

	// ErrGetKumactlCode represents the errors which are generated
	// during the process of using kumactl for installation
	ErrGetKumactlCode = "kuma_test_code"

	// ErrDownloadBinaryCode represents the errors which are generated
	// during the process of downloading binary
	ErrDownloadBinaryCode = "kuma_test_code"

	// ErrInstallBinaryCode represents the errors which are generated
	// during the process of installing binary
	ErrInstallBinaryCode = "kuma_test_code"

	// ErrUntarCode represents the errors which are generated
	// during the process of untaring a compressed file
	ErrUntarCode = "kuma_test_code"

	// ErrUntarDefaultCode represents the errors which are generated
	// during the process of untaring a compressed file
	ErrUntarDefaultCode = "kuma_test_code"

	// ErrMoveBinaryCode represents the errors which are generated
	// during the process of moving binaries
	ErrMoveBinaryCode = "kuma_test_code"

	// ErrCustomOperationCode represents the errors which are generated
	// during the process of handeling a custom process
	ErrCustomOperationCode = "kuma_test_code"

	// ErrOpInvalid represents the errors which are generated
	// when an operation is invalid
	ErrOpInvalid = errors.NewDefault(errors.ErrOpInvalid, "Invalid operation")

	// ErrUntarDefault represents the errors which are generated
	// during the process of untaring a compressed file
	ErrUntarDefault = errors.NewDefault(ErrUntarDefaultCode, "Error untaring operation default")
)

// ErrInstallKuma is the error for install mesh
func ErrInstallKuma(err error) error {
	return errors.NewDefault(ErrInstallKumaCode, fmt.Sprintf("Error with kuma operation: %s", err.Error()))
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.NewDefault(ErrMeshConfigCode, fmt.Sprintf("Error configuration mesh: %s", err.Error()))
}

// ErrFetchManifest is the error occured during the process
// fetching manifest
func ErrFetchManifest(err error, des string) error {
	return errors.NewDefault(ErrFetchManifestCode, fmt.Sprintf("Error fetching mesh manifest: %s", des))
}

// ErrClientConfig is the error for setting client config
func ErrClientConfig(err error) error {
	return errors.NewDefault(ErrClientConfigCode, fmt.Sprintf("Error setting client config: %s", err.Error()))
}

// ErrClientSet is the error for setting clientset
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

// ErrDownloadBinary is the error for downloading binary
func ErrDownloadBinary(err error) error {
	return errors.NewDefault(ErrDownloadBinaryCode, fmt.Sprintf("Error downloadinf kumactl binary: %s", err.Error()))
}

// ErrUntar is the error for streaming event
func ErrUntar(err error) error {
	return errors.NewDefault(ErrUntarCode, fmt.Sprintf("Error untaring package: %s", err.Error()))
}

// ErrInstallBinary is the error for streaming event
func ErrInstallBinary(err error) error {
	return errors.NewDefault(ErrInstallBinaryCode, fmt.Sprintf("Error installing kumactl: %s", err.Error()))
}

// ErrMoveBinary is the error for streaming event
func ErrMoveBinary(err error) error {
	return errors.NewDefault(ErrMoveBinaryCode, fmt.Sprintf("Error with moving binary kumactl: %s", err.Error()))
}

// ErrCustomOperation is the error occured during the process of
// applying custom operation
func ErrCustomOperation(err error) error {
	return errors.NewDefault(ErrCustomOperationCode, fmt.Sprintf("Error applying custom operation: %s", err.Error()))
}
