package config

import (
	"github.com/spf13/viper"
)

// Viper instance for configuration
type Viper struct {
	instance *viper.Viper
}

// NewViper intializes a viper instance and dependencies
func NewViper() (Handler, error) {
	v := viper.New()
	v.AddConfigPath(filepath)
	v.SetConfigType(filetype)
	v.SetConfigName(filename)
	v.AutomaticEnv()

	v.SetDefault("server", server)
	v.SetDefault("mesh", mesh)

	err := v.ReadInConfig()
	return &Viper{
		instance: v,
	}, ErrViper(err)
}

// SetKey sets a key value in viper
func (v *Viper) SetKey(key string, value string) {
	v.instance.Set(key, value)
}

// GetKey gets a key value from viper
func (v *Viper) GetKey(key string) string {
	return v.instance.Get(key).(string)
}

// Server provides server specific configuration
func (v *Viper) Server(result interface{}) error {
	return v.instance.Unmarshal(&result)
}

// MeshSpec provides mesh specific configuration
func (v *Viper) MeshSpec(result interface{}) error {
	return v.instance.Unmarshal(&result)
}

// MeshInstance provides mesh specific configuration
func (v *Viper) MeshInstance(result interface{}) error {
	return v.instance.Unmarshal(&result)
}

// Operations provides list of operations available
func (v *Viper) Operations(result interface{}) error {
	return v.instance.Unmarshal(&result)
}
