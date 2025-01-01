package main

import (
	"context"
	"fmt"
	"net"

	"github.com/dvlpr-shivendra/qoi-decoder/service" // Adjust the import path accordingly
	"google.golang.org/grpc"
)

type grpcHandler struct {
	service.UnimplementedBackendServiceServer
}

// Implement TestConnection method
func (h *grpcHandler) TestConnection(ctx context.Context, in *service.Empty) (*service.Message, error) {
	return &service.Message{Message: "Connection successful"}, nil
}

func main() {
	// Initialize gRPC server
	grpcServer := grpc.NewServer()

	// Listen on localhost:2000
	l, err := net.Listen("tcp", "localhost:2000")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	// Register BackendServiceServer with the gRPC server
	service.RegisterBackendServiceServer(grpcServer, &grpcHandler{})

	fmt.Println("Starting gRPC server on localhost:2000")

	// Start serving gRPC requests
	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}

	// QoifImage processing
	qoifImage := NewQoif("./qoi_test_images/kodim23.qoi")
	qoifImage.Process()
}
