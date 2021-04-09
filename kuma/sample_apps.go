package kuma

import (
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
)

func (kuma *Kuma) installSampleApp(del bool, namespace string, templates []adapter.Template) (string, error) {
	st := status.Installing

	if del {
		st = status.Removing
	}

	for _, template := range templates {
		err := kuma.applyManifest(del, namespace, []byte(template.String()))
		if err != nil {
			return st, ErrSampleApp(err)
		}
	}

	return status.Installed, nil
}
