package config

import (
	"fmt"
	"os"

	"github.com/layer5io/gokit/utils"
)

// Local instance for configuration
type Local struct {
	store map[string]string
}

// NewLocal intializes a local instance and dependencies
func NewLocal() (Handler, error) {
	return &Local{
		store: make(map[string]string),
	}, nil
}

// -------------------------------------------Application config methods----------------------------------------------------------------

// SetKey sets a key value in local store
func (l *Local) SetKey(key string, value string) {
	l.store[key] = value
}

// GetKey gets a key value from local store
func (l *Local) GetKey(key string) (string, error) {
	return l.store[key], nil
}

// Server provides server specific configuration
func (l *Local) Server(result interface{}) error {

	d := fmt.Sprintf(`{
		"name":    "kuma-adapter",
		"port":    "10007",
		"traceurl": "%s",
		"version": "v1.0.0"
	}`, os.Getenv("TRACING_ENDPOINT"))
	return utils.Unmarshal(d, result)
}

// MeshSpec provides mesh specific configuration
func (l *Local) MeshSpec(result interface{}) error {
	d := `{
		"name":    "Kuma",
		"status":  "not installed",
		"version": "none"
	}`
	return utils.Unmarshal(d, result)
}

// MeshInstance provides mesh specific configuration
func (l *Local) MeshInstance(result interface{}) error {
	d := `{
		"installmode": "flat",
		"installplatform": "kubernetes",
		"installzone": " ",
		"mgmtaddr": "0.0.0.0:8000",
		"kumaaddr": "5681"
	}`
	return utils.Unmarshal(d, result)
}

// Operations provides operations in the mesh
func (l *Local) Operations(result interface{}) error {
	d, err := utils.Marshal(operations)
	if err != nil {
		return err
	}
	return utils.Unmarshal(d, result)
}
