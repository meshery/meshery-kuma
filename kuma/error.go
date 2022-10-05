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

	// ErrNilClientCode represents the error code which is
	// generated when Kubernetes client is nil
	ErrNilClientCode = "1017"

	// ErrApplyHelmChartCode represents the error which are generated
	// during the process of applying helm chart
	ErrApplyHelmChartCode = "1018"

	// ErrInvalidOAMComponentTypeCode represents the error code which is
	// generated when an invalid oam component is requested
	ErrInvalidOAMComponentTypeCode = "1021"

	// ErrKumaCoreComponentFailCode represents the error code which is
	// generated when an kuma core operations fails
	ErrKumaCoreComponentFailCode = "1022"

	// ErrProcessOAMCode represents the error code which is
	// generated when an OAM operations fails
	ErrProcessOAMCode = "1023"

	// ErrParseKumaCoreComponentCode represents the error code which is
	// generated when kuma core component manifest parsing fails
	ErrParseKumaCoreComponentCode = "1024"

	// ErrParseOAMComponentCode represents the error code which is
	// generated during the OAM component parsing
	ErrParseOAMComponentCode = "1025"

	// ErrParseOAMConfigCode represents the error code which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfigCode = "1026"

	//ErrLoadNamespaceCode occur during the process of applying namespace
	ErrLoadNamespaceCode = "1028"

	// ErrOpInvalid represents the errors which are generated
	// when an operation is invalid
	ErrOpInvalid = errors.New(ErrOpInvalidCode, errors.Alert, []string{"Invalid operation"}, []string{"Kuma adapter received an invalid operation from the meshey server"}, []string{"The operation is not supported by the adapter", "Invalid operation name"}, []string{"Check if the operation name is valid and supported by the adapter"})

	// ErrUntarDefault represents the errors which are generated
	// during the process of untaring a compressed file
	ErrUntarDefault = errors.New(ErrUntarDefaultCode, errors.Alert, []string{"Error untaring opeartion default"}, []string{"Error occurred in the process of untaring a compressed file"}, []string{"The compressed file might be corrupted"}, []string{"Clear the cache and retry the operation"})

	// ErrParseOAMComponent represents the error which is
	// generated during the OAM component parsing
	ErrParseOAMComponent = errors.New(ErrParseOAMComponentCode, errors.Alert, []string{"error parsing the component"}, []string{"Error occurred while parsing application component in the OAM request made by Meshery server"}, []string{"Could not unmarshall configuration component received via ProcessOAM gRPC call into a valid Component struct"}, []string{"Check if Meshery Server is creating valid component for ProcessOAM gRPC call. This error should never happen and can be reported as a bug in Meshery Server. Also check if Meshery Server and adapters are referring to same component struct provided in MeshKit."})

	// ErrParseOAMConfig represents the error which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfig = errors.New(ErrParseOAMConfigCode, errors.Alert, []string{"error parsing the configuration"}, []string{"Error occurred while parsing configuration in the request made by Meshery Server"}, []string{"Could not unmarshall OAM config received via ProcessOAM gRPC call into a valid Config struct"}, []string{"Check if Meshery Server is creating valid config for ProcessOAM gRPC call. This error should never happen and can be reported as a bug in Meshery Server. Also, confirm that Meshery Server and Adapters are referring to same config struct provided in MeshKit"})

	// ErrNilClient represents the error which is
	// generated when Kubernetes client is nil
	ErrNilClient = errors.New(ErrNilClientCode, errors.Alert, []string{"Kubernetes client not initialized"}, []string{"Kubernetes client is nil"}, []string{"Kubernetes client not initialized"}, []string{"Reconnect the adapter to Meshery Server"})

	// ErrGetLatestReleaseCode represents the error which is
	// generated when the latest stable version could not
	// be fetched during runtime component registeration
	ErrGetLatestReleaseCode = "1019"
)

// ErrInstallKuma is the error for install mesh
func ErrInstallKuma(err error) error {
	return errors.New(ErrInstallKumaCode, errors.Alert, []string{"Error with kuma operation"}, []string{"Error occurred while installing kuma mesh through kumactl", err.Error()}, []string{}, []string{})
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.New(ErrMeshConfigCode, errors.Alert, []string{"Error configuration mesh"}, []string{err.Error(), "Error getting MeshSpecKey config from in-memory configuration"}, []string{}, []string{"Reconnect the adapter to the Meshery Server"})
}

// ErrFetchManifest is the error occurred during the process
// fetching manifest
func ErrFetchManifest(err error, des string) error {
	return errors.New(ErrFetchManifestCode, errors.Alert, []string{"Error occurred while fetching the manifest", des}, []string{err.Error()}, []string{}, []string{})
}

// ErrClientConfig is the error for setting client config
func ErrClientConfig(err error) error {
	return errors.New(ErrClientConfigCode, errors.Alert, []string{"Error occurred while setting client config"}, []string{err.Error()}, []string{}, []string{})
}

// ErrClientSet is the error for setting clientset
func ErrClientSet(err error) error {
	return errors.New(ErrClientSetCode, errors.Alert, []string{"Error occurred while setting clientset"}, []string{err.Error()}, []string{}, []string{})
}

// ErrStreamEvent is the error for streaming event
func ErrStreamEvent(err error) error {
	return errors.New(ErrStreamEventCode, errors.Alert, []string{"Error occurred while streaming events"}, []string{err.Error()}, []string{}, []string{})
}

// ErrSampleApp is the error for applying/deleting Sample App
func ErrSampleApp(err error, status string) error {
	return errors.New(ErrSampleAppCode, errors.Alert, []string{"Error with sample app operation"}, []string{err.Error(), "Error occurred while trying to install a sample application using manifests"}, []string{"Invalid kubeclient config", "Invalid manifest"}, []string{"Reconnect your adapter to Meshery Server to refresh the kubeclient"})
}

// ErrGetKumactl is the error for getting `kumactl`
func ErrGetKumactl(err error) error {
	return errors.New(ErrGetKumactlCode, errors.Alert, []string{"Error getting kumactl comamndline"}, []string{"Error occurred while downloading`kumactl` and moving it to .meshery/bin]", err.Error()}, []string{"https://download.konghq.com/mesh-alpine/kuma-<release>-<platform>-<arch>.tar.gz might be deprecated"}, []string{})
}

// ErrDownloadBinary is the error for downloading binary
func ErrDownloadBinary(err error) error {
	return errors.New(ErrDownloadBinaryCode, errors.Alert, []string{"Error downloading kuma binary"}, []string{err.Error(), "Error occurred while download kuma binary from its release url"}, []string{"Checkout https://download.konghq.com/mesh-alpine/kuma-<release>-<platform>-<arch>.tar.gz for more details"}, []string{})
}

// ErrUntar is the error for streaming event
func ErrUntar(err error) error {
	return errors.New(ErrUntarCode, errors.Alert, []string{"Error while extracting file"}, []string{err.Error()}, []string{"The gzip might be corrupt"}, []string{"Retry the operation"})
}

// ErrInstallBinary is the error for installing binary
func ErrInstallBinary(err error) error {
	return errors.New(ErrInstallBinaryCode, errors.Alert, []string{"Error installing kumactl"}, []string{"Error occurred while installing kuma mesh through kumactl", err.Error()}, []string{}, []string{})
}

// ErrMoveBinary is the error for moving binary
func ErrMoveBinary(err error) error {
	return errors.New(ErrMoveBinaryCode, errors.Alert, []string{"Error occurred while moving the kumactl binary"}, []string{err.Error()}, []string{"Meshery adapter might not have write access"}, []string{})
}

// ErrCustomOperation is the error occurred during the process of
// applying custom operation
func ErrCustomOperation(err error) error {
	return errors.New(ErrCustomOperationCode, errors.Alert, []string{"Error with custom operation"}, []string{"Error occurred while applying custom manifest to the cluster", err.Error()}, []string{"Invalid kubeclient config", "Invalid manifest"}, []string{"Reupload the kubconfig in the Meshery Server and reconnect the adapter"})
}

// ErrParseKumaCoreComponent is the error when kuma core component manifest parsing fails
func ErrParseKumaCoreComponent(err error) error {
	return errors.New(ErrParseKumaCoreComponentCode, errors.Alert, []string{"kuma core component manifest parsing failing"}, []string{err.Error()}, []string{}, []string{})
}

// ErrInvalidOAMComponentType is the error when the OAM component name is not valid
func ErrInvalidOAMComponentType(compName string) error {
	return errors.New(ErrInvalidOAMComponentTypeCode, errors.Alert, []string{"invalid OAM component name: ", compName}, []string{}, []string{}, []string{})
}

// ErrKumaCoreComponentFail is the error when core kuma component processing fails
func ErrKumaCoreComponentFail(err error) error {
	return errors.New(ErrKumaCoreComponentFailCode, errors.Alert, []string{"error in kuma core component"}, []string{err.Error()}, []string{}, []string{})
}

// ErrProcessOAM is a generic error which is thrown when an OAM operations fails
func ErrProcessOAM(err error) error {
	return errors.New(ErrProcessOAMCode, errors.Alert, []string{"error performing OAM operations"}, []string{err.Error()}, []string{}, []string{})
}

// ErrApplyHelmChart is the error for applying helm chart
func ErrApplyHelmChart(err error) error {
	return errors.New(ErrApplyHelmChartCode, errors.Alert, []string{"Error with helm chart operation"}, []string{"Error occurred while applying Helm Chart"}, []string{err.Error()}, []string{"Invalid helm chart configuration"})
}

// ErrGetLatestRelease is the error for get latest versions
func ErrGetLatestRelease(err error) error {
	return errors.New(ErrGetLatestReleaseCode, errors.Alert, []string{"Could not get latest version"}, []string{err.Error()}, []string{"Latest version could not be found at the specified url"}, []string{})
}

// ErrLoadNamespace is the occurend while applying namespace
func ErrLoadNamespace(err error, s string) error {
	return errors.New(ErrLoadNamespaceCode, errors.Alert, []string{"Error occurred while applying namespace "}, []string{err.Error()}, []string{"Trying to access a namespace which is not available"}, []string{"Verify presence of namespace. Confirm Meshery ServiceAccount permissions"})
}
