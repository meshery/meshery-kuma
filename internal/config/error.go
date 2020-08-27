package config

import (
	"fmt"

	"github.com/layer5io/gokit/errors"
)

var (
	// ErrEmptyConfig is the error object for empty config
	ErrEmptyConfig = errors.New("700", "Config not initialized")
)

// ErrViper is the error object for viper
func ErrViper(err error) error {
	return errors.New(errors.ErrViper, fmt.Sprintf("Viper initialization failed with error: %s", err.Error()))
}
