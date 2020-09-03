package kuma

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/layer5io/gokit/smi"
	"k8s.io/client-go/kubernetes"
)

type SmiTest struct {
	ctx        context.Context
	kubeClient *kubernetes.Clientset
	smiAddress string
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

func (h *handler) smiTest(id string) error {

	e := &Event{
		Operationid: id,
		Summary:     "Deploying",
		Details:     "None",
	}

	annotations := map[string]string{
		"kuma.io/gateway": "enabled",
	}

	test, err := smi.New(context.TODO(), "kuma", h.kubeClient)
	if err != nil {
		e.Summary = "Error while creating smi-conformance tool"
		e.Details = err.Error()
		h.StreamErr(e, err)
		return err
	}

	result, err := test.Run(nil, annotations)
	if err != nil {
		e.Summary = fmt.Sprintf("Error while %s running smi-conformance test", result.Status)
		e.Details = err.Error()
		h.StreamErr(e, err)
		return err
	}

	e.Summary = fmt.Sprintf("Smi conformance test %s successfully", result.Status)
	jsondata, _ := json.Marshal(result)
	e.Details = string(jsondata)
	h.StreamInfo(e)

	return nil
}
