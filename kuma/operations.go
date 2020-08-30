package kuma

import (
	"context"
	"encoding/json"
	"fmt"

	cfg "github.com/layer5io/meshery-kuma/internal/config"
)

// Operation holds the informormation for list of operations
type Operation struct {
	Type       int32             `json:"type,string,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

// Operations hold a map of Operation objects
type Operations map[string]*Operation

// ApplyOperation applies the operation on kuma
func (h *handler) ApplyOperation(ctx context.Context, op string, id string, del bool) error {

	operations := make(Operations, 0)
	h.config.Operations(&operations)

	status := "deploying"
	e := &Event{
		Operationid: id,
		Summary:     "Deploying",
		Details:     "None",
	}

	switch op {
	case cfg.InstallKumav071, cfg.InstallKumav070, cfg.InstallKumav060:
		go func(hh *handler, ee *Event) {
			if status, err := hh.installKuma(del, operations[op].Properties["version"]); err != nil {
				e.Summary = fmt.Sprintf("Error while %s Kuma service mesh", status)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Kuma service mesh %s successfully", status)
			ee.Details = fmt.Sprintf("The Kuma service mesh is now %s.", status)
			hh.StreamInfo(e)
		}(h, e)
	case cfg.InstallSample:
		go func(hh *handler, ee *Event) {
			if status, err := h.installSampleApp(); err != nil {
				e.Summary = fmt.Sprintf("Error while %s Sample application", status)
				e.Details = err.Error()
				h.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Sample application %s successfully", status)
			ee.Details = fmt.Sprintf("The Sample application is now %s.", status)
			hh.StreamInfo(e)
		}(h, e)
	case cfg.RunSmiConformance:
		go func(hh *handler, ee *Event) {
			result, err := h.runSmiTest(context.TODO(), h.kubeClient, ee.Operationid)
			if err != nil {
				e.Summary = fmt.Sprintf("Error while %s running smi-conformance test", result.Status)
				e.Details = err.Error()
				h.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Smi conformance test %s successfully", result.Status)
			jsondata, _ := json.Marshal(result)
			ee.Details = string(jsondata)
			hh.StreamInfo(e)
		}(h, e)
	default:
		h.StreamErr(e, ErrOpInvalid)
	}

	return nil
}

// ListOperations lists the operations available
func (h *handler) ListOperations() (Operations, error) {
	operations := make(Operations, 0)
	h.config.Operations(&operations)
	return operations, nil
}
