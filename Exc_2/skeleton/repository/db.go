package repository

import (
	"ordersystem/model"
	"time"
)

type DatabaseHandler struct {
	// drinks represent all available drinks
	drinks []model.Drink
	// orders serves as order history
	orders []model.Order
}

// todo
func NewDatabaseHandler() *DatabaseHandler {
	// Init the drinks slice with some test data
	// drinks := ...
	drinks := []model.Drink{
		{ID: 1, Name: "cola", Price: 2, Description: "sadly now without cocaine"},
		{ID: 2, Name: "mate", Price: 3, Description: "Tschickwossa"},
		{ID: 3, Name: "coffee", Price: 4, Description: "Hot drink to wake you up"},
	}

	// Init orders slice with some test data
	orders := []model.Order{
		{DrinkID: 3, CreatedAt: time.Date(2025, time.January, 7, 9, 12, 33, 0, time.UTC), Amount: 6},
		{DrinkID: 1, CreatedAt: time.Date(2025, time.January, 2, 14, 58, 12, 0, time.UTC), Amount: 9},
		{DrinkID: 2, CreatedAt: time.Date(2025, time.January, 6, 18, 20, 40, 0, time.UTC), Amount: 3},
		{DrinkID: 1, CreatedAt: time.Date(2025, time.January, 4, 11, 5, 55, 0, time.UTC), Amount: 12},
		{DrinkID: 3, CreatedAt: time.Date(2025, time.January, 3, 16, 45, 22, 0, time.UTC), Amount: 4},
		{DrinkID: 2, CreatedAt: time.Date(2025, time.January, 5, 8, 10, 15, 0, time.UTC), Amount: 5},
		{DrinkID: 1, CreatedAt: time.Date(2025, time.January, 8, 20, 30, 5, 0, time.UTC), Amount: 7},
		{DrinkID: 2, CreatedAt: time.Date(2025, time.January, 1, 13, 42, 19, 0, time.UTC), Amount: 1},
		{DrinkID: 3, CreatedAt: time.Date(2025, time.January, 9, 10, 25, 48, 0, time.UTC), Amount: 8},
	}

	return &DatabaseHandler{
		drinks: drinks,
		orders: orders,
	}
}

func (db *DatabaseHandler) GetDrinks() []model.Drink {
	return db.drinks
}

func (db *DatabaseHandler) GetOrders() []model.Order {
	return db.orders
}

func (db *DatabaseHandler) GetTotalledOrders() map[uint64]uint64 {
	// calculate total orders
	// key = DrinkID, value = Amount of orders
	allOrders := make(map[uint64]uint64)
	for _, o := range db.orders {
		allOrders[o.DrinkID] += uint64(o.Amount) // because Amount is actually of a different data type, we convert here
	}
	return allOrders
}

func (db *DatabaseHandler) AddOrder(order *model.Order) {
	// append order to db.orders slice
	db.orders = append(db.orders, *order)
}
