package kuma

import (
	"github.com/layer5io/meshkit/errors"
)

var (
	// Error code for failed service mesh installation

	// ErrOpInvalidCode represents the errors which are generated
	// when an invalid operation is invoked
	ErrOpInvalidCode = "1002"

	// ErrInstallKumaCode represents the errors which are generated
	// during kuma service mesh install process
	ErrInstallKumaCode = "1003"

	// ErrMeshConfigCode represents the errors which are generated
	// when an invalid mesh config is found
	ErrMeshConfigCode = "1004"

	// ErrFetchManifestCode represents the errors which are generated
	// during the process of fetching manifest
	ErrFetchManifestCode = "1005"

	// ErrClientConfigCode represents the errors which are generated
	// during the process of setting client config
	ErrClientConfigCode = "1006"

	// ErrClientSetCode represents the errors which are generated
	// during the process of setting clientset
	ErrClientSetCode = "1007"

	// ErrStreamEventCode represents the errors which are generated
	// during the process of streaming events
	ErrStreamEventCode = "1008"

	// ErrSampleAppCode represents the errors which are generated
	// during the process of installing sample app
	ErrSampleAppCode = "1009"

	// ErrGetKumactlCode represents the errors which are generated
	// during the process of using kumactl for installation
	ErrGetKumactlCode = "1010"

	// ErrDownloadBinaryCode represents the errors which are generated
	// during the process of downloading binary
	ErrDownloadBinaryCode = "1011"

	// ErrInstallBinaryCode represents the errors which are generated
	// during the process of installing binary
	ErrInstallBinaryCode = "1012"

	// ErrUntarCode represents the errors which are generated
	// during the process of untaring a compressed file
	ErrUntarCode = "1013"

	// ErrUntarDefaultCode represents the errors which are generated
	// during the process of untaring a compressed file
	ErrUntarDefaultCode = "1014"

	// ErrMoveBinaryCode represents the errors which are generated
	// during the process of moving binaries
	ErrMoveBinaryCode = "1015"

	// ErrCustomOperationCode represents the errors which are generated
	// during the process of handeling a custom process
	ErrCustomOperationCode = "1016"

	// ErrOpInvalid represents the errors which are generated
	// when an operation is invalid
	ErrOpInvalid = errors.New(ErrOpInvalidCode, errors.Alert, []string{"Invalid operation"}, []string{"Invalid operation detected by the Adaptor"}, []string{}, []string{})

	// ErrUntarDefault represents the errors which are generated
	// during the process of untaring a compressed file
	ErrUntarDefault = errors.New(ErrUntarDefaultCode, errors.Alert, []string{"Error untaring opeartion default"}, []string{"Error occured in the process of untaring a compressed file"}, []string{}, []string{})
)

// ErrInstallKuma is the error for install mesh
func ErrInstallKuma(err error) error {
	return errors.New(ErrInstallKumaCode, errors.Alert, []string{"Error occured while installing Kuma"}, []string{err.Error()}, []string{}, []string{})
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.New(ErrMeshConfigCode, errors.Alert, []string{"Error occured while configuring the mesh"}, []string{err.Error()}, []string{}, []string{})
}

// ErrFetchManifest is the error occured during the process
// fetching manifest
func ErrFetchManifest(err error, des string) error {
	return errors.New(ErrFetchManifestCode, errors.Alert, []string{"Error occured while fetching the mainfest", des}, []string{err.Error()}, []string{}, []string{})
}

// ErrClientConfig is the error for setting client config
func ErrClientConfig(err error) error {
	return errors.New(ErrClientConfigCode, errors.Alert, []string{"Error occured while setting client config"}, []string{err.Error()}, []string{}, []string{})
}

// ErrClientSet is the error for setting clientset
func ErrClientSet(err error) error {
	return errors.New(ErrClientSetCode, errors.Alert, []string{"Error occured while setting clientset"}, []string{err.Error()}, []string{}, []string{})
}

// ErrStreamEvent is the error for streaming event
func ErrStreamEvent(err error) error {
	return errors.New(ErrStreamEventCode, errors.Alert, []string{"Error occured while streaming events"}, []string{err.Error()}, []string{}, []string{})
}

// ErrSampleApp is the error for applying/deleting Sample App
func ErrSampleApp(err error, status string) error {
	return errors.New(ErrSampleAppCode, errors.Alert, []string{"Error occured while applying sample app mainifest"}, []string{"Error occured while " + status + " Sample App", err.Error()}, []string{}, []string{})
}

// ErrGetKumactl is the error for getting `kumactl`
func ErrGetKumactl(err error) error {
	return errors.New(ErrGetKumactlCode, errors.Alert, []string{"Error getting kumactl comamndline"}, []string{err.Error()}, []string{}, []string{})
}

// ErrDownloadBinary is the error for downloading binary
func ErrDownloadBinary(err error) error {
	return errors.New(ErrDownloadBinaryCode, errors.Alert, []string{"Error occured while downloading kumactl binary"}, []string{err.Error()}, []string{}, []string{})
}

// ErrUntar is the error for streaming event
func ErrUntar(err error) error {
	return errors.New(ErrUntarCode, errors.Alert, []string{"Error occured while untaring a pakage"}, []string{err.Error()}, []string{}, []string{})
}

// ErrInstallBinary is the error for installing binary
func ErrInstallBinary(err error) error {
	return errors.New(ErrInstallBinaryCode, errors.Alert, []string{"Error installing kumactl"}, []string{err.Error()}, []string{}, []string{})
}

// ErrMoveBinary is the error for moving binary
func ErrMoveBinary(err error) error {
	return errors.New(ErrMoveBinaryCode, errors.Alert, []string{"Error occured while moving the kumactl binary"}, []string{err.Error()}, []string{}, []string{})
}

// ErrCustomOperation is the error occured during the process of
// applying custom operation
func ErrCustomOperation(err error) error {
	return errors.New(ErrCustomOperationCode, errors.Alert, []string{"Error occured while applying custom opearaion"}, []string{err.Error()}, []string{}, []string{})
}
