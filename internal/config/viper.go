package config

import (
	"github.com/layer5io/gokit/utils"

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
	err := v.WriteConfig()
	if err != nil {
		return nil, ErrViper(err)
	}

	return &Viper{
		instance: v,
	}, nil
}

// SetKey sets a key value in viper
func (v *Viper) SetKey(key string, value string) {
	v.instance.Set(key, value)
}

// GetKey gets a key value from viper
func (v *Viper) GetKey(key string) (string, error) {
	err := v.instance.ReadInConfig()
	if err != nil {
		return " ", ErrViper(err)
	}

	s, err := utils.Marshal(v.instance.Get(key))
	if err != nil {
		return " ", ErrViper(err)
	}

	return s, nil
}

// Server provides server specific configuration
func (v *Viper) Server(result interface{}) error {
	s, err := v.GetKey("server")
	if err != nil {
		return ErrViper(err)
	}

	return utils.Unmarshal(s, &result)
}

// MeshSpec provides mesh specific configuration
func (v *Viper) MeshSpec(result interface{}) error {
	err := v.instance.ReadInConfig()
	if err != nil {
		return ErrViper(err)
	}
	return v.instance.Unmarshal(&result)
}

// MeshInstance provides mesh specific configuration
func (v *Viper) MeshInstance(result interface{}) error {
	err := v.instance.ReadInConfig()
	if err != nil {
		return ErrViper(err)
	}
	return v.instance.Unmarshal(&result)
}

// Operations provides list of operations available
func (v *Viper) Operations(result interface{}) error {
	err := v.instance.ReadInConfig()
	if err != nil {
		return ErrViper(err)
	}
	return v.instance.Unmarshal(&result)
}
