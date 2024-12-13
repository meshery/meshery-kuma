package kuma

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/layer5io/meshery-adapter-library/meshes"
	"github.com/layer5io/meshery-kuma/internal/config"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"gopkg.in/yaml.v2"
)

// CompHandler is the type for functions which can handle OAM components
type CompHandler func(*Kuma, v1alpha1.Component, bool, []string) (string, error)

// HandleComponents handles the processing of OAM components
func (kuma *Kuma) HandleComponents(comps []v1alpha1.Component, isDel bool, kubeconfigs []string) (string, error) {
	var errs []error
	var msgs []string

	compFuncMap := map[string]CompHandler{
		"KumaMesh": handleComponentKumaMesh,
	}

	stat1 := "deploying"
	stat2 := "deployed"
	if isDel {
		stat1 = "removing"
		stat2 = "removed"
	}
	for _, comp := range comps {
		ee := &meshes.EventsResponse{
			OperationId:   uuid.New().String(),
			Component:     config.ServerConfig["type"],
			ComponentName: config.ServerConfig["name"],
		}
		fnc, ok := compFuncMap[comp.Spec.Type]
		if !ok {
			msg, err := handleKumaCoreComponent(kuma, comp, isDel, "", "", kubeconfigs)
			if err != nil {
				kuma.streamErr(fmt.Sprintf("Error while %s %s", stat1, strings.TrimSuffix(comp.Spec.Type, ".KUMA")), ee, err)
				errs = append(errs, err)
				continue
			}
			ee.Summary = fmt.Sprintf("%s %s.", strings.TrimSuffix(comp.Spec.Type, ".KUMA"), stat2)
			ee.Details = fmt.Sprintf("The %s is now %s.", strings.TrimSuffix(comp.Spec.Type, ".KUMA"), stat2)
			kuma.StreamInfo(ee)
			msgs = append(msgs, msg)
			continue
		}

		msg, err := fnc(kuma, comp, isDel, kubeconfigs)
		if err != nil {
			kuma.streamErr(fmt.Sprintf("Error while %s %s", stat1, comp.Spec.Type), ee, err)
			errs = append(errs, err)
			continue
		}
		ee.Summary = fmt.Sprintf("%s %s %s.", comp.Name, comp.Spec.Type, stat2)
		ee.Details = fmt.Sprintf("The %s %s is now %s.", comp.Name, comp.Spec.Type, stat2)
		kuma.StreamInfo(ee)
		msgs = append(msgs, msg)
	}

	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), err
	}

	return mergeMsgs(msgs), nil
}

// HandleApplicationConfiguration handles the processing of OAM application configuration
func (kuma *Kuma) HandleApplicationConfiguration(config v1alpha1.Configuration, isDel bool, kubeconfigs []string) (string, error) {
	var errs []error
	var msgs []string
	for _, comp := range config.Spec.Components {
		for _, trait := range comp.Traits {
			if trait.Name == "automaticSidecarInjection.Kuma" {
				namespaces := castSliceInterfaceToSliceString(trait.Properties["namespaces"].([]interface{}))
				if err := handleNamespaceLabel(kuma, namespaces, isDel, kubeconfigs); err != nil {
					errs = append(errs, err)
				}
			}

			msgs = append(msgs, fmt.Sprintf("applied trait \"%s\" on service \"%s\"", trait.Name, comp.ComponentName))
		}
	}

	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), err
	}

	return mergeMsgs(msgs), nil
}

func handleNamespaceLabel(kuma *Kuma, namespaces []string, isDel bool, kubeconfigs []string) error {
	var errs []error
	for _, ns := range namespaces {
		if err := kuma.sidecarInjection(ns, isDel, kubeconfigs); err != nil {
			errs = append(errs, err)
		}
	}

	return mergeErrors(errs)
}

func handleComponentKumaMesh(kuma *Kuma, comp v1alpha1.Component, isDel bool, kubeconfigs []string) (string, error) {
	// Get the kuma version from the settings
	// we are sure that the version of kuma would be present
	// because the configuration is already validated against the schema
	version := comp.Spec.Version
	if version == "" {
		return "", fmt.Errorf("pass valid version inside service for Istio installation")
	}
	//TODO: When no version is passed in service, use the latest Kuma version

	msg, err := kuma.installKuma(isDel, false, comp.Namespace, version, kubeconfigs)
	if err != nil {
		return fmt.Sprintf("%s: %s", comp.Name, msg), err
	}

	return fmt.Sprintf("%s: %s", comp.Name, msg), nil
}

func handleKumaCoreComponent(
	kuma *Kuma,
	comp v1alpha1.Component,
	isDel bool,
	apiVersion,
	kind string,
	kubeconfigs []string) (string, error) {
	if apiVersion == "" {
		apiVersion = getAPIVersionFromComponent(comp)
		if apiVersion == "" {
			return "", ErrKumaCoreComponentFail(fmt.Errorf("failed to get API Version for: %s", comp.Name))
		}
	}

	if kind == "" {
		kind = getKindFromComponent(comp)
		if kind == "" {
			return "", ErrKumaCoreComponentFail(fmt.Errorf("failed to get kind for: %s", comp.Name))
		}
	}

	component := map[string]interface{}{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata": map[string]interface{}{
			"name":        comp.Name,
			"annotations": comp.Annotations,
			"labels":      comp.Labels,
		},
		"spec": comp.Spec.Settings,
	}

	// Convert to yaml
	yamlByt, err := yaml.Marshal(component)
	if err != nil {
		err = ErrParseKumaCoreComponent(err)
		kuma.Log.Error(err)
		return "", err
	}

	msg := fmt.Sprintf("created %s \"%s\" in namespace \"%s\"", kind, comp.Name, comp.Namespace)
	if isDel {
		msg = fmt.Sprintf("deleted %s config \"%s\" in namespace \"%s\"", kind, comp.Name, comp.Namespace)
	}

	return msg, kuma.applyManifest(isDel, comp.Namespace, yamlByt, kubeconfigs)
}

func getAPIVersionFromComponent(comp v1alpha1.Component) string {
	return comp.Annotations["pattern.meshery.io.mesh.workload.k8sAPIVersion"]
}

func getKindFromComponent(comp v1alpha1.Component) string {
	return comp.Annotations["pattern.meshery.io.mesh.workload.k8sKind"]
}

func castSliceInterfaceToSliceString(in []interface{}) []string {
	var out []string

	for _, v := range in {
		cast, ok := v.(string)
		if ok {
			out = append(out, cast)
		}
	}

	return out
}

func mergeErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var errMsgs []string

	for _, err := range errs {
		errMsgs = append(errMsgs, err.Error())
	}

	return fmt.Errorf("%s", strings.Join(errMsgs, "\n"))
}

func mergeMsgs(strs []string) string {
	return strings.Join(strs, "\n")
}
