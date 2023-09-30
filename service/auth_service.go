package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/katerji/bank/db"
	"github.com/katerji/bank/db/query"
	proto "github.com/katerji/bank/generated"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	proto.UnimplementedAuthServiceServer
}

func (service AuthService) Register(_ context.Context, request *proto.RegisterRequest) (*proto.Customer, error) {
	if request.GetName() == "" || request.GetEmail() == "" || request.GetPassword() == "" || request.GetPhoneNumber() == "" {
		return nil, errors.New("invalid request")
	}

	hashedPassword, err := hashPassword(request.GetPassword())
	if err != nil {
		return nil, err
	}
	customerID, err := db.GetDbInstance().Insert(query.InsertCustomerQuery, request.GetName(), request.GetPhoneNumber(), request.GetEmail(), hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("email %s already exists", request.Email)
	}

	customer := &proto.Customer{
		Id:          int32(customerID),
		Name:        request.GetName(),
		Email:       request.GetEmail(),
		PhoneNumber: request.GetPhoneNumber(),
	}

	return customer, nil
}

func (service AuthService) Login(_ context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
	customer := &proto.Customer{}
	client := db.GetDbInstance()
	row := client.FetchOne(query.GetCustomerByEmail, request.GetEmail())
	err := row.Scan(&customer.Id, &customer.Email, &customer.Password, &customer.Name, &customer.PhoneNumber)
	if err != nil {
		return nil, errors.New("email does not exist")
	}

	if !validPassword(customer.Password, request.Password) {
		return nil, errors.New("incorrect password")
	}

	customer.Password = ""

	jwtService := JWTService{}
	pair, err := jwtService.CreateJWTPair(customer)
	if err != nil {
		return nil, err
	}
	return &proto.LoginResponse{
		Customer: customer,
		JwtPair:  pair,
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
