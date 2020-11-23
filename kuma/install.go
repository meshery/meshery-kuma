package kuma

import (
	"bytes"
	"os/exec"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

func (kuma *Kuma) installKuma(del bool, namespace string, version string) (string, error) {
	st := status.Installing

	if del {
		st = status.Removing
	}

	err := kuma.Config.GetObject(adapter.MeshSpecKey, kuma)
	if err != nil {
		return st, ErrMeshConfig(err)
	}

	manifest, err := kuma.fetchManifest(version)
	if err != nil {
		kuma.Log.Error(ErrInstallKuma(err))
		return st, ErrInstallKuma(err)
	}

	err = kuma.applyManifest(del, namespace, []byte(manifest))
	if err != nil {
		kuma.Log.Error(ErrInstallKuma(err))
		return st, ErrInstallKuma(err)
	}

	if del {
		return status.Removed, nil
	}
	return status.Installed, nil
}

func (kuma *Kuma) fetchManifest(version string) (string, error) {

	var (
		out bytes.Buffer
		er  bytes.Buffer
	)

	Executable, err := exec.LookPath("./kumactl")
	if err != nil {
		e := GetKumactl(version)
		if e != nil {
			return "", ErrFetchManifest(e, e.Error())
		}
	}

	command := exec.Command(Executable, "install", "control-plane")
	command.Stdout = &out
	command.Stderr = &er
	err = command.Run()
	if err != nil {
		return "", ErrFetchManifest(err, er.String())
	}

	return out.String(), nil
}

func (kuma *Kuma) applyManifest(del bool, namespace string, contents []byte) error {

	err := kuma.MesheryKubeclient.ApplyManifest(contents, mesherykube.ApplyOptions{
		Namespace: namespace,
		Update:    true,
		Delete:    del,
	})
	if err != nil {
		return err
	}

	return nil
}
