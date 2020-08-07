package config

import (
	"fmt"

	"github.com/kumarabd/gokit/errors"
)

var (
	ErrEmptyConfig = errors.New("700", "Config not initialized")
)

func ErrViper(err error) error {
	return errors.New("701", fmt.Sprintf("Viper initialization failed with error: %s", err.Error()))
}
