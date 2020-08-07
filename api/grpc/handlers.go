package grpc

import (
	"context"

	"github.com/layer5io/meshery-kuma/meshes"
)

// CreateMeshInstance is the handler function for the method CreateMeshInstance.
func (s *Service) CreateMeshInstance(ctx context.Context, req *meshes.CreateMeshInstanceRequest) (*meshes.CreateMeshInstanceResponse, error) {
	s.Handler.CreateInstance(req.K8SConfig, req.ContextName)
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
	return &meshes.ApplyRuleResponse{
		Error:       " ",
		OperationId: " ",
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
			Category: meshes.OpCategory(meshes.OpCategory_value[val.Type]),
		})
	}

	return &meshes.SupportedOperationsResponse{
		Ops:   operations,
		Error: "none",
	}, nil
}

// StreamEvents is the handler function for the method StreamEvents.
func (s *Service) StreamEvents(ctx *meshes.EventsRequest, srv meshes.MeshService_StreamEventsServer) error {
	return nil
}
