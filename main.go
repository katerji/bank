package main

import (
	"fmt"
	"github.com/katerji/bank/db"
	"github.com/katerji/bank/server"
	"log"
	"net"
)

func main() {
	initDB()
	initGRPCServer()
}

func initDB() {
	client := db.GetDbInstance()
	err := client.Ping()
	if err != nil {
		panic(err)
	}
}

func initGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 88))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := server.NewGRPCServer()
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
}
