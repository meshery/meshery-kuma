package grpc

import (
	"context"
	"time"

	"github.com/layer5io/meshery-kuma/kuma"
	"github.com/layer5io/meshery-kuma/meshes"
)

// CreateMeshInstance is the handler function for the method CreateMeshInstance.
func (s *Service) CreateMeshInstance(ctx context.Context, req *meshes.CreateMeshInstanceRequest) (*meshes.CreateMeshInstanceResponse, error) {
	err := s.Handler.CreateInstance(req.K8SConfig, req.ContextName, &s.Channel)
	if err != nil {
		return nil, err
	}
	return &meshes.CreateMeshInstanceResponse{}, nil
}

// MeshName is the handler function for the method MeshName.
func (s *Service) MeshName(ctx context.Context, req *meshes.MeshNameRequest) (*meshes.MeshNameResponse, error) {
	return &meshes.MeshNameResponse{
		Name: s.Handler.GetName(),
	}, nil
}

// ApplyOperation is the handler function for the method ApplyOperation.
func (s *Service) ApplyOperation(ctx context.Context, req *meshes.ApplyRuleRequest) (*meshes.ApplyRuleResponse, error) {

	if req == nil {
		return nil, ErrRequestInvalid
	}

	err := s.Handler.ApplyOperation(ctx, req.OpName, req.OperationId, req.DeleteOp)
	if err != nil {
		return nil, err
	}

	return &meshes.ApplyRuleResponse{
		OperationId: req.OperationId,
	}, nil
}

// SupportedOperations is the handler function for the method SupportedOperations.
func (s *Service) SupportedOperations(ctx context.Context, req *meshes.SupportedOperationsRequest) (*meshes.SupportedOperationsResponse, error) {
	result, err := s.Handler.ListOperations()
	if err != nil {
		return nil, err
	}

	operations := make([]*meshes.SupportedOperation, 0)
	for key, val := range result {
		operations = append(operations, &meshes.SupportedOperation{
			Key:      key,
			Value:    val.Properties["description"],
			Category: meshes.OpCategory(val.Type),
		})
	}

	return &meshes.SupportedOperationsResponse{
		Ops:   operations,
		Error: "none",
	}, nil
}

// StreamEvents is the handler function for the method StreamEvents.
func (s *Service) StreamEvents(ctx *meshes.EventsRequest, srv meshes.MeshService_StreamEventsServer) error {
	for {
		select {
		case data := <-s.Channel:
			event := &meshes.EventsResponse{
				OperationId: data.(*kuma.Event).Operationid,
				EventType:   meshes.EventType(data.(*kuma.Event).EType),
				Summary:     data.(*kuma.Event).Summary,
				Details:     data.(*kuma.Event).Details,
			}
			if err := srv.Send(event); err != nil {
				// to prevent loosing the event, will re-add to the channel
				go func() {
					s.Channel <- data
				}()
				return err
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}
