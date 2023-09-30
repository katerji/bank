package query

const InsertAccountQuery = "INSERT INTO account (name, customer_id) VALUES (?, ?)"
const FetchAccountOwnerQuery = "SELECT customer_id FROM account WHERE id = ?"
const DepositQuery = "UPDATE account SET balance = balance + ? WHERE id = ?"
