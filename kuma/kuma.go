package kuma

import (
	"github.com/kumarabd/gokit/logger"
	"github.com/layer5io/meshery-kuma/internal/config"
)

// Handler provides the methods supported by the adaptor
type Handler interface {
	GetName() string
	CreateInstance([]byte, string) error
	ApplyOperation() error
	ListOperations() (Operations, error)
	Stream() error
}

// handler holds the dependencies for kuma-adaptor
type handler struct {
	config config.Handler
	log    logger.Handler
}

// New initializes email handler.
func New(c config.Handler, l logger.Handler) Handler {
	return &handler{
		config: c,
		log:    l,
	}
}
