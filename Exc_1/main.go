package main

import "fmt"

type Car struct {
	Name       string
	Year       int
	Color      string
	Horsepower int
}

// Define a method with a receiver to print the struct
func (c Car) Greet() {
	fmt.Printf("struct1: %v\n", c)
}

func main() {

	car := Car{
		Name:       "Ferrari",
		Year:       2000,
		Color:      "Red",
		Horsepower: 380,
	}

	car.Greet()
}
