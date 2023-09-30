package query

const CreateAccountQuery = "INSERT INTO account (name, customer_id) VALUES (?, ?)"
const CloseAccountQuery = "UPDATE account SET deleted = 1 WHERE id = ?"
const DepositQuery = "UPDATE account SET balance = balance + ? WHERE id = ?"
const FetchAccountQuery = "SELECT id, name, balance, customer_id FROM account WHERE id = ? AND deleted = 0"
const WithdrawQuery = "UPDATE account SET balance = balance - ? WHERE id = ?"
