package kuma

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/layer5io/gokit/smi"
)

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
