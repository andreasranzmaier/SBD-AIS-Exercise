package client

import (
	"context"
	"exc8/pb"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcClient struct {
	client pb.OrderServiceClient
}

func NewGrpcClient() (*GrpcClient, error) {
	conn, err := grpc.NewClient(":4000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewOrderServiceClient(conn)
	return &GrpcClient{client: client}, nil
}

func (c *GrpcClient) Run() error {
	ctx := context.Background()

	// 1. List drinks
	fmt.Println("Requesting drinks 🍹🍺☕")
	drinksResp, err := c.client.GetDrinks(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	fmt.Println("Available drinks:")
	for _, d := range drinksResp.Drinks {
		fmt.Printf("\t> id:%d  name:\"%s\"  price:%g  description:\"%s\"\n", d.Id, d.Name, d.Price, d.Description)
	}

	// 2. Order a few drinks
	fmt.Println("Ordering drinks 👨‍🍳⏱️🍻🍻")
	orders1 := []struct {
		id   int32
		name string
		qty  int32
	}{
		{1, "Spritzer", 2},
		{2, "Beer", 2},
		{3, "Coffee", 2},
	}
	for _, o := range orders1 {
		fmt.Printf("\t> Ordering: %d x %s\n", o.qty, o.name)
		_, err := c.client.OrderDrink(ctx, &pb.OrderDrinkRequest{DrinkId: o.id, Quantity: o.qty})
		if err != nil {
			return err
		}
	}

	// 3. Order more drinks
	fmt.Println("Ordering another round of drinks 👨‍🍳⏱️🍻🍻")
	orders2 := []struct {
		id   int32
		name string
		qty  int32
	}{
		{1, "Spritzer", 6},
		{2, "Beer", 6},
		{3, "Coffee", 6},
	}
	for _, o := range orders2 {
		fmt.Printf("\t> Ordering: %d x %s\n", o.qty, o.name)
		_, err := c.client.OrderDrink(ctx, &pb.OrderDrinkRequest{DrinkId: o.id, Quantity: o.qty})
		if err != nil {
			return err
		}
	}

	// 4. Get order total
	fmt.Println("Getting the bill Stonks 💹💹💹")
	ordersResp, err := c.client.GetOrders(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	for _, o := range ordersResp.Orders {
		fmt.Printf("\t> Total: %d x %s\n", o.Quantity, o.Drink.Name)
	}

	return nil
}
