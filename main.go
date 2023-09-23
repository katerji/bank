package main

import (
	"fmt"
	"github.com/katerji/bank/db"
	proto "github.com/katerji/bank/generated"
	"github.com/katerji/bank/service"
	"google.golang.org/grpc"
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
	authService := service.AuthService{}
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 88))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, authService)

	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}
