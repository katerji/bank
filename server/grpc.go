package server

import (
	"context"
	"errors"
	proto "github.com/katerji/bank/generated"
	"github.com/katerji/bank/service"
	"github.com/katerji/bank/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func withServerUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(serverInterceptor)
}
func NewGRPCServer() *grpc.Server {
	s := grpc.NewServer(withServerUnaryInterceptor())

	authService := service.AuthService{}
	bankService := service.BankService{}
	proto.RegisterAuthServiceServer(s, authService)
	proto.RegisterBankServiceServer(s, bankService)
	return s
}

func serverInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	protectedMethods := []string{
		"/BankService/CreateAccount",
	}
	method := info.FullMethod
	if utils.InSlice(protectedMethods, method) {
		customer, err := authorize(ctx)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, utils.Customer, customer)
	}

	h, err := handler(ctx, req)

	return h, err
}

func authorize(ctx context.Context) (*proto.Customer, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("retrieving metadata is failed")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return nil, errors.New("authorization token is not supplied")
	}
	token := authHeader[0][7:]

	jwtService := &service.JWTService{}
	customer, err := jwtService.VerifyToken(token)

	if err != nil {
		return nil, errors.New("invalid token")
	}
	return customer, nil
}
