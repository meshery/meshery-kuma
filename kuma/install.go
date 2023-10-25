package kuma

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sync"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshery-kuma/internal/config"
	"github.com/layer5io/meshkit/models"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"gopkg.in/yaml.v2"
)

const (
	kumaRepository = "https://kumahq.github.io/charts"
	kumaChartName  = "kuma"
)

// CreateKubeconfigs creates and writes passed kubeconfig onto the filesystem
func (kuma *Kuma) CreateKubeconfigs(kubeconfigs []string) error {
	var errs = make([]error, 0)
	for _, kubeconfig := range kubeconfigs {
		kconfig := models.Kubeconfig{}
		err := yaml.Unmarshal([]byte(kubeconfig), &kconfig)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// To have control over what exactly to take in on kubeconfig
		kuma.KubeconfigHandler.SetKey("kind", kconfig.Kind)
		kuma.KubeconfigHandler.SetKey("apiVersion", kconfig.APIVersion)
		kuma.KubeconfigHandler.SetKey("current-context", kconfig.CurrentContext)
		err = kuma.KubeconfigHandler.SetObject("preferences", kconfig.Preferences)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = kuma.KubeconfigHandler.SetObject("clusters", kconfig.Clusters)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = kuma.KubeconfigHandler.SetObject("users", kconfig.Users)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = kuma.KubeconfigHandler.SetObject("contexts", kconfig.Contexts)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return mergeErrors(errs)
}

func (kuma *Kuma) installKuma(del bool, useManifest bool, namespace string, version string, kubeconfigs []string) (string, error) {
	st := status.Installing

	if del {
		st = status.Removing
	}

	err := kuma.Config.GetObject(adapter.MeshSpecKey, kuma)
	if err != nil {
		return st, ErrMeshConfig(err)
	}
	if useManifest {
		return kuma.installUsingManifests(del, st, namespace, version, kubeconfigs)
	}
	kuma.Log.Info("Installing kuma using helm charts...")
	err = kuma.applyHelmChart(del, version, namespace, kubeconfigs)
	if err != nil {
		kuma.Log.Info("Failed helm installation. ", err)
		kuma.Log.Info("Trying installing from manifests...")
		return kuma.installUsingManifests(del, st, namespace, version, kubeconfigs)
	}
	return status.Installed, nil
}
func (kuma *Kuma) installUsingManifests(del bool, st string, namespace string, version string, kubeconfigs []string) (string, error) {
	kuma.Log.Info("Installing kuma using manifests...")
	manifest, err := kuma.fetchManifest(version)
	if err != nil {
		kuma.Log.Error(ErrInstallKuma(err))
		return st, ErrInstallKuma(err)
	}

	err = kuma.applyManifest(del, namespace, []byte(manifest), kubeconfigs)
	if err != nil {
		kuma.Log.Error(ErrInstallKuma(err))
		return st, ErrInstallKuma(err)
	}

	if del {
		return status.Removed, nil
	}
	return status.Installed, nil
}
func (kuma *Kuma) applyHelmChart(del bool, version, namespace string, kubeconfigs []string) error {
	chartVersion, err := mesherykube.HelmAppVersionToChartVersion(
		kumaRepository,
		kumaChartName,
		version,
	)
	if err != nil {
		return ErrApplyHelmChart(err)
	}
	var act mesherykube.HelmChartAction
	if del {
		act = mesherykube.UNINSTALL
	} else {
		act = mesherykube.INSTALL
	}
	var errs []error
	var wg sync.WaitGroup
	var errMx sync.Mutex
	for _, kubeconfig := range kubeconfigs {
		wg.Add(1)
		go func(kubeconfig string) {
			defer wg.Done()
			kClient, err := mesherykube.New([]byte(kubeconfig))
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
			err = kClient.ApplyHelmChart(mesherykube.ApplyHelmChartConfig{
				ChartLocation: mesherykube.HelmChartLocation{
					Repository: kumaRepository,
					Chart:      kumaChartName,
					Version:    chartVersion,
				},
				Namespace:       namespace,
				Action:          act,
				CreateNamespace: true,
				ReleaseName:     kumaChartName,
			})
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
		}(kubeconfig)
	}
	wg.Wait()

	if len(errs) != 0 {
		return ErrApplyHelmChart(mergeErrors(errs))
	}
	return nil
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

func (kuma *Kuma) applyManifest(del bool, namespace string, contents []byte, kubeconfigs []string) error {
	var errs []error
	var wg sync.WaitGroup
	var errMx sync.Mutex
	for _, kubconfig := range kubeconfigs {
		wg.Add(1)
		go func(kubconfig string) {
			defer wg.Done()
			kClient, err := mesherykube.New([]byte(kubconfig))
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
			err = kClient.ApplyManifest(contents, mesherykube.ApplyOptions{
				Namespace: namespace,
				Update:    true,
				Delete:    del,
			})
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
		}(kubconfig)
	}
	wg.Wait()
	if len(errs) != 0 {
		return mergeErrors(errs)
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
	if _, err = os.Stat(executable); err == nil {
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

	// Move binary to the right location
	err = os.Rename(path.Join(binPath, alternateBinaryName, "kuma-"+release, "bin", "kumactl"), path.Join(binPath, "kumactl"))
	if err != nil {
		return "", ErrGetKumactl(err)
	}

	// Cleanup
	kuma.Log.Info("Cleaning up...")
	if err = os.RemoveAll(path.Join(binPath, alternateBinaryName)); err != nil {
		return "", ErrGetKumactl(err)
	}

	if err = os.Rename(path.Join(binPath, "kumactl"), path.Join(binPath, alternateBinaryName)); err != nil {
		return "", ErrGetKumactl(err)
	}

	// Set permissions
	// Permsission has to be +x to be able to run the binary
	// #nosec
	if err = os.Chmod(path.Join(binPath, alternateBinaryName), 0750); err != nil {
		return "", ErrGetKumactl(err)
	}

	kuma.Log.Info("Done")
	return path.Join(binPath, alternateBinaryName), nil
}

func downloadBinary(platform, arch, release string) (*http.Response, error) {
	var url = fmt.Sprintf("https://download.konghq.com/mesh-alpine/kuma-%s-%s-%s.tar.gz", release, platform, arch)

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

	err := os.MkdirAll(location, 0750)
	if err != nil {
		return err
	}

	switch platform {
	case "darwin":
		fallthrough
	case "linux":
		uncompressedStream, err := gzip.NewReader(res.Body)
		if err != nil {
			return err
		}

		tarReader := tar.NewReader(uncompressedStream)

		for {
			header, err := tarReader.Next()

			if err == io.EOF {
				break
			}

			if err != nil {
				return ErrInstallBinary(err)
			}

			switch header.Typeflag {
			case tar.TypeDir:
				// File traversal is required to store the binary at the right place
				if err = os.MkdirAll(path.Join(location, header.Name), 0750); err != nil {
					return ErrInstallBinary(err)
				}
			case tar.TypeReg:
				// File traversal is required to store the binary at the right place
				outFile, errCreateFile := os.Create(path.Join(location, header.Name))
				if errCreateFile != nil {
					return ErrInstallBinary(err)
				}
				// Trust kuma tar
				if _, err = io.Copy(outFile, tarReader); err != nil {
					return ErrInstallBinary(err)
				}
				if err = outFile.Close(); err != nil {
					return ErrInstallBinary(err)
				}
			default:
				return ErrInstallBinary(err)
			}
		}
	}
	return nil
}
