package kuma

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/layer5io/learn-layer5/smi-conformance/conformance"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/layer5io/gokit/utils"
)

type SmiTest struct {
	ctx            context.Context
	kubeClient     *kubernetes.Clientset
	kubeConfigPath string
	smiAddress     string
	helmPath       string
	namespace      string
}

// ConformanceResponse holds the response object of the test
type ConformanceResponse struct {
	Tests    string                       `json:"tests,omitempty"`
	Failures string                       `json:"failures,omitempty"`
	Results  []*SingleConformanceResponse `json:"results,omitempty"`
	Status   string                       `json:"status,omitempty"`
}

// Failure is the failure response object
type Failure struct {
	Text    string `json:"text,omitempty"`
	Message string `json:"message,omitempty"`
}

// SingleConformanceResponse holds the result of one particular test case
type SingleConformanceResponse struct {
	Name       string   `json:"name,omitempty"`
	Time       string   `json:"time,omitempty"`
	Assertions string   `json:"assertions,omitempty"`
	Failure    *Failure `json:"failure,omitempty"`
}

func (h *handler) runSmiTest(ctx context.Context, client *kubernetes.Clientset, opid string) (ConformanceResponse, error) {

	e := &Event{
		Operationid: opid,
	}

	test := SmiTest{
		ctx:            ctx,
		kubeClient:     client,
		helmPath:       h.config.GetKey("smi-helm-path"),
		namespace:      h.config.GetKey("smi-namespace"),
		kubeConfigPath: h.kubeConfigPath,
	}

	response := ConformanceResponse{
		Tests:    "None",
		Failures: "None",
		Results:  make([]*SingleConformanceResponse, 0),
		Status:   "deploying",
	}

	err := test.installConformanceTool()
	if err != nil {
		response.Status = "installing"
		return response, ErrInstallSmi(err)
	}
	e.Summary = "Tool deployed"
	h.StreamInfo(e)

	err = test.connectConformanceTool()
	if err != nil {
		response.Status = "connecting"
		return response, ErrConnectSmi(err)
	}
	e.Summary = "Tool connected"
	h.StreamInfo(e)

	err = test.runConformanceTest("kuma", &response)
	if err != nil {
		response.Status = "running"
		return response, ErrRunSmi(err)
	}
	e.Summary = "Executing the test"
	h.StreamInfo(e)

	err = test.deleteConformanceTool()
	if err != nil {
		response.Status = "deleting"
		return response, ErrDeleteSmi(err)
	}

	response.Status = "completed"

	return response, nil
}

// installConformanceTool installs the smi conformance tool
func (test *SmiTest) installConformanceTool() error {

	_, err := test.kubeClient.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: test.namespace,
			Annotations: map[string]string{
				"meta.helm.sh/release-name":      "smi-conformance",
				"meta.helm.sh/release-namespace": test.namespace,
			},
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "Helm",
			},
		}}, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	path := "/tmp/smi-conformance.tar.gz"
	err = utils.DownloadFile(path, test.helmPath)
	if err != nil {
		return err
	}

	chart, err := loader.Load(path)
	if err != nil {
		return err
	}

	actionConfig := &action.Configuration{}
	if err := actionConfig.Init(kube.GetConfig(test.kubeConfigPath, "", test.namespace), test.namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}); err != nil {
		return err
	}

	iCli := action.NewInstall(actionConfig)
	iCli.Namespace = test.namespace
	iCli.ReleaseName = "smi-conformance"
	_, err = iCli.Run(chart, nil)
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second) // Required for all the resources to be created

	return nil
}

// deleteConformanceTool deletes the smi conformance tool
func (test *SmiTest) deleteConformanceTool() error {
	err := test.kubeClient.CoreV1().Namespaces().Delete(context.TODO(), test.namespace, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// connectConformanceTool initiates the connection
func (test *SmiTest) connectConformanceTool() error {
	var host string
	var port int32

	svc, err := test.kubeClient.CoreV1().Services("meshery").Get(test.ctx, "smi-conformance", metav1.GetOptions{})
	if err != nil {
		return err
	}

	nodes, err := test.kubeClient.CoreV1().Nodes().List(test.ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	addresses := make(map[string]string, 0)
	for _, addr := range nodes.Items[0].Status.Addresses {
		addresses[string(addr.Type)] = addr.Address
	}
	host = addresses["ExternalIP"]
	port = svc.Spec.Ports[0].NodePort
	if tcpCheck(addresses["InternalIP"], port) {
		host = addresses["InternalIP"]
	}

	test.smiAddress = fmt.Sprintf("%s:%d", host, port)
	return nil
}

func tcpCheck(ip string, port int32) bool {
	timeout := 5 * time.Second
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		return true
	}
	return false
}

// runConformanceTest runs the conformance test
func (test *SmiTest) runConformanceTest(adaptorname string, response *ConformanceResponse) error {
	labels := make(map[string]string, 0)
	annotations := map[string]string{
		"kuma.io/gateway": "enabled",
	}

	cClient, err := conformance.CreateClient(context.TODO(), test.smiAddress)
	if err != nil {
		return err
	}
	defer cClient.Close()

	result, err := cClient.CClient.RunTest(context.TODO(), &conformance.Request{
		Annotations: annotations,
		Labels:      labels,
		Meshname:    adaptorname,
	})
	if err != nil {
		return err
	}

	if result == nil {
		return err
	}

	response.Tests = result.Tests
	response.Failures = result.Failures

	for _, res := range result.SingleTestResult {
		response.Results = append(response.Results, &SingleConformanceResponse{
			Name:       res.Name,
			Time:       res.Time,
			Assertions: res.Assertions,
			Failure: &Failure{
				Text:    res.Failure.Test,
				Message: res.Failure.Message,
			},
		})
	}

	return nil
}
