package kuma

import (
	"context"
	"fmt"

	"github.com/layer5io/gokit/errors"
	"github.com/layer5io/gokit/logger"
)

type loggingService struct {
	log  logger.Handler
	next Handler
}

func AddLogger(logger logger.Handler, h Handler) Handler {
	return &loggingService{
		log:  logger,
		next: h,
	}
}

func (s *loggingService) GetName() string {
	return s.next.GetName()
}

func (s *loggingService) CreateInstance(b []byte, st string, c *chan interface{}) error {
	s.log.Info("Creating instance")
	err := s.next.CreateInstance(b, st, c)
	if err != nil {
		s.log.Err(errors.GetCode(err), err.Error())
	}
	return err
}

func (s *loggingService) ApplyOperation(ctx context.Context, op string, id string, del bool) error {
	s.log.Info(fmt.Sprintf("Applying operation %s", op))
	err := s.next.ApplyOperation(ctx, op, id, del)
	if err != nil {
		s.log.Err(errors.GetCode(err), err.Error())
	}
	return err
}

func (s *loggingService) ListOperations() (Operations, error) {
	s.log.Info("Listing Operations")
	ops, err := s.next.ListOperations()
	if err != nil {
		s.log.Err(errors.GetCode(err), err.Error())
	}
	return ops, err
}

func (s *loggingService) StreamErr(e *Event, err error) {
	s.log.Err("Sending error event", err.Error())
}

func (s *loggingService) StreamInfo(*Event) {
	s.log.Info("Sending event")
}
