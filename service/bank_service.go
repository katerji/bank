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

func (b BankService) Deposit(ctx context.Context, request *proto.DepositRequest) (*proto.GenericResponse, error) {
	customer := ctx.Value(utils.Customer).(*proto.Customer)
	amount := request.GetAmount()
	accountID := int(request.GetAccountId())
	if amount <= 0 || accountID == 0 {
		return nil, errors.New("bad request")
	}
	accountOwnerID := getAccountOwner(accountID)
	if accountOwnerID != int(customer.Id) {
		return nil, errors.New("unauthorized transaction")
	}
	ok := db.GetDbInstance().Exec(query.DepositQuery, amount)

	return &proto.GenericResponse{
		Success: ok,
	}, nil
}

func getAccountOwner(accountID int) int {
	row := db.GetDbInstance().FetchOne(query.FetchAccountOwnerQuery, accountID)
	var ownerID int
	if err := row.Scan(&ownerID); err != nil {
		return 0
	}
	return ownerID
}
