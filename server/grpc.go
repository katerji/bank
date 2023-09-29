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

func NewGRPCServer() *grpc.Server {
	s := grpc.NewServer(grpc.UnaryInterceptor(getAuthMiddleware))

	authService := service.AuthService{}
	bankService := service.BankService{}
	proto.RegisterAuthServiceServer(s, authService)
	proto.RegisterBankServiceServer(s, bankService)

	return s
}

func getAuthMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if !utils.InSlice(getUnProtectedMethods(), info.FullMethod) {
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
		return nil, errors.New("retrieving metadata has failed")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return nil, errors.New("authorization token missing")
	}
	token := authHeader[0][7:]

	jwtService := &service.JWTService{}
	customer, err := jwtService.VerifyToken(token)

	if err != nil {
		return nil, errors.New("invalid token")
	}

	return customer, nil
}

func getUnProtectedMethods() []string {
	return []string{
		"/AuthService/Register",
		"/AuthService/Login",
	}
}
