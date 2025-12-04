package server

import (
	"context"
	"exc8/pb"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type GRPCService struct {
	pb.UnimplementedOrderServiceServer
	drinks []*pb.Drink
	orders []*pb.Order
	mu     sync.Mutex // locking when accessing orders
}

func StartGrpcServer() error {
	// Create a new gRPC server.
	srv := grpc.NewServer()
	// Create grpc service
	grpcService := &GRPCService{
		drinks: []*pb.Drink{ // predefined drinks as list of pb.Drink
			{Id: 1, Name: "Spritzer", Price: 2.0, Description: "Wine with soda"},
			{Id: 2, Name: "Beer", Price: 3.0, Description: "Hagenberger Gold"},
			{Id: 3, Name: "Coffee", Price: 2.5, Description: "Mifare isn't that secure"},
		},
		orders: []*pb.Order{}, // start with empty orders
	}
	// Register our service implementation with the gRPC server.
	pb.RegisterOrderServiceServer(srv, grpcService)
	// Serve gRPC server on port 4000.
	lis, err := net.Listen("tcp", ":4000")
	if err != nil {
		return err
	}
	return srv.Serve(lis)
}
