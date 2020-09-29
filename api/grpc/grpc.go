package grpc

import (
	"fmt"
	"net"
	"time"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"google.golang.org/grpc"

	"github.com/layer5io/meshery-kuma/kuma"
	"github.com/layer5io/meshery-kuma/meshes"
)

// Service object holds all the information about the server parameters.
type Service struct {
	Name      string    `json:"name"`
	Port      string    `json:"port"`
	Version   string    `json:"version"`
	StartedAt time.Time `json:"startedat,string"`
	TraceURL  string    `json:"traceurl"`
	Handler   kuma.Handler
	Channel   chan interface{}
}

// Options holds the dependency optional objects for the server.
type Options struct{}

// panicHandler is the handler function to handle panic errors
func panicHandler(r interface{}) error {
	fmt.Println("600 Error")
	return ErrPanic(r)
}

// Start starts grpc server
func Start(s *Service, options ...Options) error {

	address := fmt.Sprintf(":%s", s.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return ErrGrpcListener(err)
	}

	middlewares := middleware.ChainUnaryServer(
		grpc_recovery.UnaryServerInterceptor(
			grpc_recovery.WithRecoveryHandler(panicHandler),
		),
	)
	// if tr != nil {
	// 	middlewares = middleware.ChainUnaryServer(
	// 		grpctrace.UnaryServerInterceptor(tr.Tracer(s.Name).(apitrace.Tracer)),
	// 	)
	// }

	server := grpc.NewServer(
		grpc.UnaryInterceptor(middlewares),
	)

	//Register Proto
	meshes.RegisterMeshServiceServer(server, s)

	// Start serving requests
	if err = server.Serve(listener); err != nil {
		return ErrGrpcServer(err)
	}
	return nil

}
