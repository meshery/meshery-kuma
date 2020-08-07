package kuma

import (
	"fmt"
	"os"
	"os/exec"
)

// MeshInstance holds the information of the instance of the mesh
type MeshInstance struct {
	InstallMode     string `json:"installmode,omitempty"`
	InstallPlatform string `json:"installplatform,omitempty"`
	InstallZone     string `json:"installzone,omitempty"`
	InstallVersion  string `json:"installversion,omitempty"`
	MgmtAddr        string `json:"mgmtaddr,omitempty"`
	Kumaaddr        string `json:"kumaaddr,omitempty"`
}

// CreateInstance installs and creates a mesh environment up and running
func (h *handler) CreateInstance(k8sconfig []byte, kubecontext string) error {
	meshinstance := &MeshInstance{}
	err := h.config.MeshInstance(meshinstance)
	if err != nil {
		return ErrMeshConfig(err)
	}

	h.log.Info("Installing Kuma.......................")
	err = meshinstance.installMesh()
	if err != nil {
		h.log.Err("Kuma installation failed!!", ErrInstallMesh(err).Error())
		return ErrInstallMesh(err)
	}

	h.log.Info("Port forwarding!!!")
	err = meshinstance.portForward()
	if err != nil {
		h.log.Err("Kuma portforwarding failed!!", ErrPortForward(err).Error())
		return ErrPortForward(err)
	}
	return nil
}

// installMesh installs the mesh in the cluster or the target location
func (m *MeshInstance) installMesh() error {
	Executable, err := exec.LookPath("./scripts/installer.sh")
	if err != nil {
		return err
	}

	cmd := &exec.Cmd{
		Path:   Executable,
		Args:   []string{Executable},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("KUMA_VERSION=%s", m.InstallVersion),
		fmt.Sprintf("KUMA_MODE=%s", m.InstallMode),
		fmt.Sprintf("KUMA_PLATFORM=%s", m.InstallPlatform),
		fmt.Sprintf("KUMA_ZONE=%s", m.InstallZone),
	)

	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (m *MeshInstance) portForward() error {
	// Needs implementation
	return nil
}
