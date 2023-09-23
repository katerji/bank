package query

const (
	InsertCustomerQuery = "INSERT INTO customer (name, phone_number, email, password) VALUES (?, ?)"
	GetCustomerByEmail  = "SELECT id, email, password, name, phone_number FROM customer WHERE email = ?"
)
