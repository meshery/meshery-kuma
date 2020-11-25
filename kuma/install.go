package kuma

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshery-kuma/internal/config"
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

	Executable, err := kuma.getExecutable(version)
	if err != nil {
		return "", ErrFetchManifest(err, err.Error())
	}
	// We need variable executable hence
	// #nosec
	command := exec.Command(Executable, "install", "control-plane")
	command.Stdout = &out
	command.Stderr = &er
	err = command.Run()
	if err != nil {
		kuma.Log.Info(out.String())
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

// getExecutable looks for the executable in
// 1. $PATH
// 2. Root config path
//
// If it doesn't find the executable in the path then it proceeds
// to download the binary from github releases and installs it
// in the root config path
func (kuma *Kuma) getExecutable(release string) (string, error) {
	const binaryName = "kumactl"
	alternateBinaryName := "kumactl-" + release

	// Look for the executable in the path
	kuma.Log.Info("Looking for kuma in the path...")
	executable, err := exec.LookPath(binaryName)
	if err == nil {
		return executable, nil
	}
	executable, err = exec.LookPath(alternateBinaryName)
	if err == nil {
		return executable, nil
	}

	// Look for config in the root path
	binPath := path.Join(config.RootPath(), "bin")
	kuma.Log.Info("Looking for kuma in", binPath, "...")
	executable = path.Join(binPath, alternateBinaryName)
	if _, err := os.Stat(executable); err == nil {
		return executable, nil
	}

	// Proceed to download the binary in the config root path
	kuma.Log.Info("kuma not found in the path, downloading...")
	res, err := downloadBinary(os.Getenv("DISTRO"), runtime.GOARCH, release)
	if err != nil {
		return "", ErrGetKumactl(err)
	}
	// Install the binary
	kuma.Log.Info("Installing...")
	if err = installBinary(path.Join(binPath, alternateBinaryName), runtime.GOOS, res); err != nil {
		return "", ErrGetKumactl(err)
	}

	kuma.Log.Info("Done")
	return path.Join(binPath, alternateBinaryName), nil
}

func downloadBinary(platform, arch, release string) (*http.Response, error) {
	var url = fmt.Sprintf("https://kong.bintray.com/kuma/kuma-%s-%s-%s.tar.gz", release, platform, arch)

	// We need variable url hence
	// #nosec
	resp, err := http.Get(url)
	if err != nil {
		return nil, ErrDownloadBinary(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrDownloadBinary(fmt.Errorf("binary not found, possibly the operating system is not supported"))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrDownloadBinary(fmt.Errorf("bad status: %s", resp.Status))
	}

	return resp, nil
}

func installBinary(location, platform string, res *http.Response) error {
	// Close the response body
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	out, err := os.Create(location)
	if err != nil {
		return err
	}
	defer func() {
		if err := out.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	switch platform {
	case "darwin":
		fallthrough
	case "linux":
		r, err := gzip.NewReader(res.Body)
		if err != nil {
			return err
		}

		// Trust Kuma tar
		// #nosec
		_, err = io.Copy(out, r)
		if err != nil {
			return ErrInstallBinary(err)
		}

		if err = r.Close(); err != nil {
			return ErrInstallBinary(err)
		}

		if err = out.Chmod(0750); err != nil {
			return ErrInstallBinary(err)
		}
	case "windows":
	}
	return nil
}
