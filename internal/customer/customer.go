package customer

import "fmt"

// Greeting interface exports Hello behavior
type Greeting interface {
	Hello() string
}

// Customer defines a customer name and greeting
type Customer struct {
	name, greeting string
}

/*
NewCustomer factory function returns a point to a Customer and nil

It receives a name and greeting strings and both strings cannot be empty, if empty it will return an error.
*/
func NewCustomer(name, greeting string) (Greeting, error) {
	customer := Customer{}
	if name == "" {
		return &customer, fmt.Errorf("Customer cannot have an empty name")
	}
	if greeting == "" {
		return &customer, fmt.Errorf("Customer cannot have an empty greeting")
	}
	customer.name = name
	customer.greeting = greeting
	return &customer, nil
}

/*
Hello function will return a greeting from a Customer
*/
func (c Customer) Hello() string {
	return c.greeting
}

/*
String returned a string version of customer
*/
func (c Customer) String() string {
	return fmt.Sprint("name: ", c.name, "greeting: ", c.greeting)
}
