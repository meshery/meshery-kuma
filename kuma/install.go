package kuma

import (
	"bytes"
	"os/exec"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

func (kuma *Kuma) installKuma(del bool, version string) (string, error) {
	st := status.Installing

	if del {
		st = status.Removing
	}

	err := kuma.Config.GetObject(adapter.MeshSpecKey, kuma)
	if err != nil {
		return st, ErrMeshConfig(err)
	}

	manifest, err := kuma.fetchManifest()
	if err != nil {
		kuma.Log.Error(ErrInstallKuma(err))
		return st, ErrInstallKuma(err)
	}

	err = kuma.applyManifest([]byte(manifest))
	if err != nil {
		kuma.Log.Error(ErrInstallKuma(err))
		return st, ErrInstallKuma(err)
	}

	if del {
		return status.Removed, nil
	}
	return status.Installed, nil
}

func (kuma *Kuma) fetchManifest() (string, error) {

	var (
		out bytes.Buffer
		er  bytes.Buffer
	)

	Executable, err := exec.LookPath("kumactl")
	if err != nil {
		return "", ErrFetchManifest(err, "")
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

func (kuma *Kuma) applyManifest(contents []byte) error {
	kclient, err := mesherykube.New(kuma.KubeClient, kuma.RestConfig)
	if err != nil {
		return err
	}

	err = kclient.ApplyManifest(contents, mesherykube.ApplyOptions{})
	if err != nil {
		return err
	}

	return nil
}
