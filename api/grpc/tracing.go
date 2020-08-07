package grpc

// import (
// 	"context"

// 	"github.com/layer5io/meshery-kuma/internal/tracing"
// 	"github.com/layer5io/meshery-kuma/meshes"
// )

// type tracingService struct {
// 	Tracer tracing.Handler
// 	Next   Service
// }

// func AddTracer(tracer tracing.Handler, s Service) Service {
// 	return &tracingService{
// 		Tracer: tracer,
// 		Next:   s,
// 	}
// }

// // CreateMeshInstance is the handler function for the method CreateMeshInstance.
// func (tr *tracingService) CreateMeshInstance(ctx context.Context, req *meshes.CreateMeshInstanceRequest) (*meshes.CreateMeshInstanceResponse, error) {

// 	tr.Tracer.Span(ctx)
// 	atts := make([]*tracing.KeyValue, 0)

// 	result, err := tr.Next.CreateMeshInstance(ctx, req)
// 	if err != nil {
// 		append(atts, &tracing.KeyValue{
// 			Key:   "error",
// 			Value: err.Error(),
// 		})
// 	}
// 	append(atts, &tracing.KeyValue{
// 		Key:   "result",
// 		Value: "success",
// 	})
// 	tr.Tracer.AddEvent("CreateMeshInstance", atts)
// 	return result, err
// }

// // MeshName is the handler function for the method MeshName.
// func (tr *tracingService) MeshName(ctx context.Context, req *meshes.MeshNameRequest) (*meshes.MeshNameResponse, error) {
// 	tr.Tracer.Span(ctx)
// 	atts := make([]*tracing.KeyValue, 0)

// 	result, err := tr.Next.MeshName(ctx, req)
// 	if err != nil {
// 		append(atts, &tracing.KeyValue{
// 			Key:   "error",
// 			Value: err.Error(),
// 		})
// 	}
// 	append(atts, &tracing.KeyValue{
// 		Key:   "result",
// 		Value: "success",
// 	})
// 	tr.Tracer.AddEvent("CreateMeshInstance", atts)
// 	return result, err
// }

// // ApplyOperation is the handler function for the method ApplyOperation.
// func (tr *tracingService) ApplyOperation(ctx context.Context, req *meshes.ApplyRuleRequest) (*meshes.ApplyRuleResponse, error) {
// 	tr.Tracer.Span(ctx)
// 	atts := make([]*tracing.KeyValue, 0)

// 	result, err := tr.Next.ApplyOperation(ctx, req)
// 	if err != nil {
// 		append(atts, &tracing.KeyValue{
// 			Key:   "error",
// 			Value: err.Error(),
// 		})
// 	}
// 	append(atts, &tracing.KeyValue{
// 		Key:   "result",
// 		Value: "success",
// 	})
// 	tr.Tracer.AddEvent("CreateMeshInstance", atts)
// 	return result, err
// }

// // SupportedOperations is the handler function for the method SupportedOperations.
// func (tr *tracingService) SupportedOperations(ctx context.Context, req *meshes.SupportedOperationsRequest) (*meshes.SupportedOperationsResponse, error) {
// 	tr.Tracer.Span(ctx)
// 	atts := make([]*tracing.KeyValue, 0)

// 	result, err := tr.Next.SupportedOperations(ctx, req)
// 	if err != nil {
// 		append(atts, &tracing.KeyValue{
// 			Key:   "error",
// 			Value: err.Error(),
// 		})
// 	}
// 	append(atts, &tracing.KeyValue{
// 		Key:   "result",
// 		Value: "success",
// 	})
// 	tr.Tracer.AddEvent("CreateMeshInstance", atts)
// 	return result, err
// }

// // StreamEvents is the handler function for the method StreamEvents.
// func (tr *tracingService) StreamEvents(ctx *meshes.EventsRequest, srv meshes.MeshService_StreamEventsServer) error {
// 	tr.Tracer.Span(ctx)
// 	atts := make([]*tracing.KeyValue, 0)

// 	result, err := tr.Next.StreamEvents(ctx, srv)
// 	if err != nil {
// 		append(atts, &tracing.KeyValue{
// 			Key:   "error",
// 			Value: err.Error(),
// 		})
// 	}
// 	append(atts, &tracing.KeyValue{
// 		Key:   "result",
// 		Value: "success",
// 	})
// 	tr.Tracer.AddEvent("CreateMeshInstance", atts)
// 	return result, err
// }
