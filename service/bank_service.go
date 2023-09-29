package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/katerji/bank/db"
	"github.com/katerji/bank/db/query"
	proto "github.com/katerji/bank/generated"
	"github.com/katerji/bank/utils"
)

type BankService struct {
	proto.UnimplementedBankServiceServer
}

func (b BankService) CreateAccount(ctx context.Context, request *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	customer := ctx.Value(utils.Customer).(*proto.Customer)
	name := request.GetName()
	if name == "" {
		return nil, errors.New("bad request")
	}

	accountID, err := db.GetDbInstance().Insert(query.InsertAccountQuery, name, customer.GetId())
	if err != nil {
		return nil, fmt.Errorf("account with name %s already exists", request.GetName())
	}

	return &proto.CreateAccountResponse{
		Account: &proto.Account{
			Id:         int32(accountID),
			Name:       request.Name,
			CustomerId: customer.Id,
			Balance:    0,
		},
	}, nil
}
