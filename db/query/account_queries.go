package query

const InsertAccountQuery = "INSERT INTO account (name, customer_id) VALUES (?, ?)"
const DepositQuery = "UPDATE account SET balance = balance + ? WHERE id = ?"
const FetchAccountQuery = "SELECT id, name, balance, customer_id FROM account WHERE id = ?"
const WithdrawQuery = "UPDATE account SET balance = balance - ? WHERE id = ?"
