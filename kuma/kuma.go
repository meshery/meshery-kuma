package kuma

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/user"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/layer5io/gokit/logger"
	"github.com/layer5io/gokit/models"
	"github.com/layer5io/meshery-kuma/internal/config"
	"gopkg.in/yaml.v2"
)

// Handler provides the methods supported by the adaptor
type Handler interface {
	GetName() string
	CreateInstance([]byte, string, *chan interface{}) error
	ApplyOperation(context.Context, string, string, bool) error
	ListOperations() (Operations, error)

	StreamErr(*Event, error)
	StreamInfo(*Event)
}

// handler holds the dependencies for kuma-adaptor
type handler struct {
	config  config.Handler
	log     logger.Handler
	channel *chan interface{}

	kubeClient *kubernetes.Clientset
}

// New initializes email handler.
func New(c config.Handler, l logger.Handler) Handler {
	return &handler{
		config: c,
		log:    l,
	}
}

// newClient creates a new client
func (h *handler) CreateInstance(kubeconfig []byte, contextName string, ch *chan interface{}) error {

	h.channel = ch

	config, err := clientConfig(kubeconfig, contextName)
	if err != nil {
		return ErrClientConfig(err)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return ErrClientSet(err)
	}

	h.kubeClient = clientset

	return nil
}

// configClient creates a config client
func clientConfig(kubeconfig []byte, contextName string) (*rest.Config, error) {
	if len(kubeconfig) > 0 {
		ccfg, err := clientcmd.Load(kubeconfig)
		if err != nil {
			return nil, err
		}
		if contextName != "" {
			ccfg.CurrentContext = contextName
		}
		writeKubeconfig(kubeconfig, contextName)
		return clientcmd.NewDefaultClientConfig(*ccfg, &clientcmd.ConfigOverrides{}).ClientConfig()
	}
	return rest.InClusterConfig()
}

// writeKubeconfig creates kubeconfig in local container
func writeKubeconfig(kubeconfig []byte, contextName string) error {

	yamlConfig := models.Kubeconfig{}
	err := yaml.Unmarshal(kubeconfig, &yamlConfig)
	if err != nil {
		return err
	}

	yamlConfig.CurrentContext = contextName

	d, err := yaml.Marshal(yamlConfig)
	if err != nil {
		return err
	}

	user, err := user.Current()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/.kube/config", user.HomeDir), d, 0600)
	if err != nil {
		return err
	}

	return nil
}
