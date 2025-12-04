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

func (s *GRPCService) GetDrinks(ctx context.Context, _ *emptypb.Empty) (*pb.GetDrinksResponse, error) {
	//   ^ 1. Receiver    		^ 2. Input Parameters                   ^ 3. Return Values
	return &pb.GetDrinksResponse{Drinks: s.drinks}, nil
}

func (s *GRPCService) OrderDrink(ctx context.Context, req *pb.OrderDrinkRequest) (*wrapperspb.BoolValue, error) {
	s.mu.Lock()         // locking when modifying orders
	defer s.mu.Unlock() // unlocking after modification

	// Validate drink exists
	var foundDrink *pb.Drink
	for _, d := range s.drinks {
		if d.Id == req.DrinkId {
			foundDrink = d
			break
		}
	}
	if foundDrink == nil {
		return &wrapperspb.BoolValue{Value: false}, fmt.Errorf("drink with id %d not found", req.DrinkId)
	}

	// Check if we already have an order for this drink
	var existingOrder *pb.Order
	for _, o := range s.orders {
		if o.Drink.Id == req.DrinkId {
			existingOrder = o
			break // found existing order, stop searching
		}
	}

	// Update existing order or add new one
	if existingOrder != nil {
		existingOrder.Quantity += req.Quantity
	} else {
		s.orders = append(s.orders, &pb.Order{
			Drink:    foundDrink,
			Quantity: req.Quantity,
		})
	}

	return &wrapperspb.BoolValue{Value: true}, nil
}

func (s *GRPCService) GetOrders(ctx context.Context, _ *emptypb.Empty) (*pb.GetOrdersResponse, error) {
	//   ^ 1. Receiver    		^ 2. Input Parameters                   ^ 3. Return Values
	s.mu.Lock()         // locking when accessing orders
	defer s.mu.Unlock() // unlocking after access
	return &pb.GetOrdersResponse{Orders: s.orders}, nil
}