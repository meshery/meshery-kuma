package kuma

import (
	"context"
	"fmt"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kuma *Kuma) installSampleApp(del bool, namespace string, templates []adapter.Template) (string, error) {
	st := status.Installing

	if del {
		st = status.Removing
	}

	for _, template := range templates {
		err := kuma.applyManifest(del, namespace, []byte(template.String()))
		if err != nil {
			return st, ErrSampleApp(err, st)
		}
	}

	return status.Installed, nil
}

func (kuma *Kuma) sidecarInjection(namespace string, del bool) error {
	kclient := kuma.KubeClient
	if kclient == nil {
		return ErrNilClient
	}

	// updating the label on the namespace
	ns, err := kclient.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if ns.ObjectMeta.Labels == nil {
		ns.ObjectMeta.Labels = map[string]string{}
	}
	ns.ObjectMeta.Labels["kuma.io/monitored-by"] = "kuma"

	if del {
		delete(ns.ObjectMeta.Labels, "kuma.io/monitored-by")
	}

	// updating the annotations on the namespace
	if ns.ObjectMeta.Annotations == nil {
		ns.ObjectMeta.Annotations = map[string]string{}
	}
	ns.ObjectMeta.Annotations["kuma.io/sidecar-injection"] = "enabled"

	if del {
		delete(ns.ObjectMeta.Annotations, "kuma.io/sidecar-injection")
	}

	fmt.Println(ns.ObjectMeta)

	_, err = kclient.CoreV1().Namespaces().Update(context.TODO(), ns, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
