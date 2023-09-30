package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/katerji/bank/envs"
	"time"
)

type Client struct {
	db *sql.DB
}

var instance *Client

func GetDbInstance() *Client {
	if instance == nil {
		instance, _ = getDbClient()
	}
	return instance
}

func getDbClient() (*Client, error) {
	dbHost := envs.GetInstance().GetDbHost()
	dbUser := envs.GetInstance().GetDbUser()
	dbPort := envs.GetInstance().GetDbPort()
	dbPass := envs.GetInstance().GetDbPassword()
	dbName := envs.GetInstance().GetDbName()

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return &Client{}, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &Client{
		db,
	}, nil
}

func (c *Client) Fetch(query string, args ...any) []any {
	rows, err := c.db.Query(query, args...)
	if err != nil {
		fmt.Println(err.Error())
	}
	results := make([]interface{}, 0)
	for rows.Next() {
		var row interface{}

		err = rows.Scan(&row)
		results = append(results, row)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
	}
	return results
}

func (c *Client) FetchOne(query string, args ...any) *sql.Row {
	return c.db.QueryRow(query, args...)
}

func (c *Client) Insert(query string, args ...any) (int, error) {
	rows, err := c.db.Prepare(query)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	exec, err := rows.Exec(args...)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	insertId, err := exec.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	return int(insertId), nil
}

func (c *Client) Exec(query string, args ...any) bool {
	rows, err := c.db.Prepare(query)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = rows.Exec(args...)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (c *Client) Ping() error {
	return c.db.Ping()
}
