package customer

import "testing"

func TestCustomerHello(t *testing.T) {
	greeting := "Hello, 世界"
	t.Log("Given a new customer with a new greeting")
	{
		c1, err := NewCustomer("A", greeting)
		if err != nil {
			t.Fatalf("\tNewCustomer factory function error %v", err)
		}
		if greeting != c1.Hello() {
			t.Fatalf("\tNewCustomer hello mismatch %s should be equal %s", c1.Hello(), greeting)
		}
	}
	t.Log("Given a new customer with a empty greeting")
	{
		_, err := NewCustomer("B", "")
		if err == nil {
			t.Fatalf("\tNewCustomer factory should return a error %v", err)
		}
	}
	t.Log("Given a empty customer with a new greeting")
	{
		_, err := NewCustomer("", greeting)
		if err == nil {
			t.Fatalf("\tNewCustomer factory should return a error %v", err)
		}
	}

}

func TestCustomerString(t *testing.T) {
	t.Log("Given a new customer with a new greeting")
	{
		c1 := Customer{"A", "Hello tester"}
		if c1.String() == "" {
			t.Fatalf("\tString function not working properly")
		}
	}
}
