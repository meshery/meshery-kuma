package config

// Handler is the handler interface for config
type Handler interface {

	// SetKey sets a key value in the config
	SetKey(key string, value string)

	// GetKey gets a key value from the config
	GetKey(key string) (string, error)

	// Server provides the server specific configuration
	Server(result interface{}) error

	// MeshSpec provides the mesh specific configuration
	MeshSpec(result interface{}) error

	// MeshInstance provides the mesh specific configuration
	MeshInstance(result interface{}) error

	// Operations provides the list of operations available
	Operations(result interface{}) error
}

// New returns the interface of the config handler
func New(name string) (Handler, error) {
	switch name {
	case "local":
		return NewLocal()
	case "viper":
		return NewViper()
	}
	return nil, ErrEmptyConfig
}
