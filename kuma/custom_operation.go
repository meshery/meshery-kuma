package kuma

import "github.com/layer5io/meshery-adapter-library/status"

func (kuma *Kuma) applyCustomOperation(namespace string, manifest string, isDel bool) (string, error) {
	st := status.Starting

	err := kuma.applyManifest(isDel, namespace, []byte(manifest))
	if err != nil {
		return st, ErrCustomOperation(err)
	}

	return status.Completed, nil
}
