// Package kuma - Common operations for the adapter
package kuma

import (
	"context"
	"fmt"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-adapter-library/status"
	internalconfig "github.com/layer5io/meshery-kuma/internal/config"
	"github.com/layer5io/meshery-kuma/kuma/oam"
	meshkitCfg "github.com/layer5io/meshkit/config"
	"github.com/layer5io/meshkit/errors"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
)

const (
	// SMIManifest is the manifest.yaml file for smi conformance tool
	SMIManifest = "https://raw.githubusercontent.com/layer5io/learn-layer5/master/smi-conformance/manifest.yml"
)

// Kuma represents the kuma adapter and embeds adapter.Adapter
type Kuma struct {
	adapter.Adapter // Type Embedded
}

// New initializes kuma handler.
func New(c meshkitCfg.Handler, l logger.Handler, kc meshkitCfg.Handler) adapter.Handler {
	return &Kuma{
		Adapter: adapter.Adapter{
			Config:            c,
			Log:               l,
			KubeconfigHandler: kc,
		},
	}
}

// ApplyOperation applies the operation on kuma
func (kuma *Kuma) ApplyOperation(ctx context.Context, opReq adapter.OperationRequest, hchan *chan interface{}) error {
	err := kuma.CreateKubeconfigs(opReq.K8sConfigs)
	if err != nil {
		return err
	}
	kuma.SetChannel(hchan)
	kubeconfigs := opReq.K8sConfigs
	operations := adapter.Operations{}
	err = kuma.Config.GetObject(adapter.OperationsKey, &operations)
	if err != nil {
		return err
	}

	e := &adapter.Event{
		Operationid: opReq.OperationID,
		Summary:     status.Deploying,
		Details:     "Operation is not supported",
		Component:   internalconfig.ServerConfig["type"],
		ComponentName: internalconfig.ServerConfig["name"],
	}

	switch opReq.OperationName {
	case internalconfig.KumaOperation:
		go func(hh *Kuma, ee *adapter.Event) {
			version := string(operations[opReq.OperationName].Versions[0])
			stat, err := hh.installKuma(opReq.IsDeleteOperation, false, opReq.Namespace, version, kubeconfigs)
			if err != nil {
				summary := fmt.Sprintf("Error while %s Kuma service mesh", stat)
				hh.streamErr(summary, e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Kuma service mesh %s successfully", stat)
			ee.Details = fmt.Sprintf("The Kuma service mesh is now %s.", stat)
			hh.StreamInfo(e)
		}(kuma, e)
	case common.BookInfoOperation, common.HTTPBinOperation, common.ImageHubOperation, common.EmojiVotoOperation:
		go func(hh *Kuma, ee *adapter.Event) {
			appName := operations[opReq.OperationName].AdditionalProperties[common.ServiceName]
			stat, err := hh.installSampleApp(opReq.IsDeleteOperation, opReq.Namespace, operations[opReq.OperationName].Templates, kubeconfigs)
			if err != nil {
				summary := fmt.Sprintf("Error while %s %s application", stat, appName)
				hh.streamErr(summary, e, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s application %s successfully", appName, stat)
			ee.Details = fmt.Sprintf("The %s application is now %s.", appName, stat)
			hh.StreamInfo(e)
		}(kuma, e)
	case common.SmiConformanceOperation:
		go func(hh *Kuma, ee *adapter.Event) {
			name := operations[opReq.OperationName].Description
			_, err := hh.RunSMITest(adapter.SMITestOptions{
				Ctx:         context.TODO(),
				OperationID: ee.Operationid,
				Kubeconfigs: kubeconfigs,
				Manifest:    string(operations[opReq.OperationName].Templates[0]),
				Namespace:   "meshery",
				Labels: map[string]string{
					"kuma.io/gateway": "enabled",
				},
				Annotations: make(map[string]string),
			})
			if err != nil {
				summary := fmt.Sprintf("Error while %s %s test", status.Running, name)
				hh.streamErr(summary, e, err)
				return
			}
		}(kuma, e)
	case common.CustomOperation:
		go func(hh *Kuma, ee *adapter.Event) {
			stat, err := hh.applyCustomOperation(opReq.Namespace, opReq.CustomBody, opReq.IsDeleteOperation, kubeconfigs)
			if err != nil {
				summary := fmt.Sprintf("Error while %s custom operation", stat)
				hh.streamErr(summary, e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Manifest %s successfully", status.Deployed)
			ee.Details = ""
			hh.StreamInfo(e)
		}(kuma, e)
	default:
		kuma.streamErr("Invalid operation", e, ErrOpInvalid)
	}

	return nil
}

// ProcessOAM will handles the grpc invocation for handling OAM objects
func (kuma *Kuma) ProcessOAM(ctx context.Context, oamReq adapter.OAMRequest, hchan *chan interface{}) (string, error) {
	err := kuma.CreateKubeconfigs(oamReq.K8sConfigs)
	if err != nil {
		return "", err
	}
	kuma.SetChannel(hchan)
	kubeconfigs := oamReq.K8sConfigs
	var comps []v1alpha1.Component
	for _, acomp := range oamReq.OamComps {
		comp, err := oam.ParseApplicationComponent(acomp)
		if err != nil {
			kuma.Log.Error(ErrParseOAMComponent)
			continue
		}

		comps = append(comps, comp)
	}

	config, err := oam.ParseApplicationConfiguration(oamReq.OamConfig)
	if err != nil {
		kuma.Log.Error(ErrParseOAMConfig)
	}

	// If operation is delete then first HandleConfiguration and then handle the deployment
	if oamReq.DeleteOp {
		// Process configuration
		msg2, err := kuma.HandleApplicationConfiguration(config, oamReq.DeleteOp, kubeconfigs)
		if err != nil {
			return msg2, ErrProcessOAM(err)
		}

		// Process components
		msg1, err := kuma.HandleComponents(comps, oamReq.DeleteOp, kubeconfigs)
		if err != nil {
			return msg1 + "\n" + msg2, ErrProcessOAM(err)
		}

		return msg1 + "\n" + msg2, nil
	}

	// Process components
	msg1, err := kuma.HandleComponents(comps, oamReq.DeleteOp, kubeconfigs)
	if err != nil {
		return msg1, ErrProcessOAM(err)
	}

	// Process configuration
	msg2, err := kuma.HandleApplicationConfiguration(config, oamReq.DeleteOp, kubeconfigs)
	if err != nil {
		return msg1 + "\n" + msg2, ErrProcessOAM(err)
	}

	return msg1 + "\n" + msg2, nil
}

func(kuma *Kuma) streamErr(summary string, e *adapter.Event, err error) {
	e.Summary = summary
	e.Details = err.Error()
	e.ErrorCode = errors.GetCode(err)
	e.ProbableCause = errors.GetCause(err)
	e.SuggestedRemediation = errors.GetRemedy(err)
	kuma.StreamErr(e, err)
}
