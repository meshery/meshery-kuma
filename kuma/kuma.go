package kuma

import (
	"context"
	"fmt"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/common"
	adapterconfig "github.com/layer5io/meshery-adapter-library/config"
	"github.com/layer5io/meshery-adapter-library/status"
	internalconfig "github.com/layer5io/meshery-kuma/internal/config"
	"github.com/layer5io/meshkit/logger"
)

type Kuma struct {
	adapter.Adapter // Type Embedded
}

// New initializes kuma handler.
func New(c adapterconfig.Handler, l logger.Handler, kc adapterconfig.Handler) adapter.Handler {
	return &Kuma{
		Adapter: adapter.Adapter{
			Config:            c,
			Log:               l,
			KubeconfigHandler: kc,
		},
	}
}

// ApplyOperation applies the operation on kuma
func (kuma *Kuma) ApplyOperation(ctx context.Context, opReq adapter.OperationRequest) error {

	operations := make(adapter.Operations, 0)
	err := kuma.Config.GetObject(adapter.OperationsKey, &operations)
	if err != nil {
		return err
	}

	st := status.Deploying

	e := &adapter.Event{
		Operationid: opReq.OperationID,
		Summary:     status.Deploying,
		Details:     status.None,
	}

	switch opReq.OperationName {
	case internalconfig.KumaOperation:
		go func(hh *Kuma, ee *adapter.Event) {
			version := string(operations[opReq.OperationName].Versions[0])
			if stat, err := hh.installKuma(opReq.IsDeleteOperation, version); err != nil {
				e.Summary = fmt.Sprintf("Error while %s Kuma service mesh", stat)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Kuma service mesh %s successfully", stat)
			ee.Details = fmt.Sprintf("The Kuma service mesh is now %s.", stat)
			hh.StreamInfo(e)
		}(kuma, e)
	case common.SmiConformanceOperation:
		go func(hh *Kuma, ee *adapter.Event) {
			err := hh.ValidateSMIConformance(&adapter.SmiTestOptions{
				Ctx:  context.TODO(),
				OpID: ee.Operationid,
			})
			if err != nil {
				return
			}
		}(kuma, e)
	default:
		kuma.StreamErr(e, ErrOpInvalid)
	}

	return nil
}
