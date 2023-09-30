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
	account, err := getAccount(accountID)
	if err != nil {
		return nil, err
	}
	if account.CustomerId != customer.Id {
		return nil, errors.New("unauthorized transaction")
	}
	ok := db.GetDbInstance().Exec(query.DepositQuery, amount)

	return &proto.GenericResponse{
		Success: ok,
	}, nil
}

func (b BankService) Withdraw(ctx context.Context, request *proto.WithdrawRequest) (*proto.GenericResponse, error) {
	customer := ctx.Value(utils.Customer).(*proto.Customer)
	amount := request.GetAmount()
	accountID := int(request.GetAccountId())
	if amount <= 0 || accountID == 0 {
		return nil, errors.New("bad request")
	}
	account, err := getAccount(accountID)
	if err != nil {
		return nil, err
	}
	if account.Id != customer.Id {
		return nil, errors.New("unauthorized transaction")
	}
	if float32(account.Balance) < amount {
		return nil, errors.New("insufficient funds")
	}
	ok := db.GetDbInstance().Exec(query.WithdrawQuery, amount)

	return &proto.GenericResponse{
		Success: ok,
	}, nil
}

func getAccount(accountID int) (*proto.Account, error) {
	account := &proto.Account{}
	row := db.GetDbInstance().FetchOne(query.FetchAccountQuery)
	if err := row.Scan(&account.Id, &account.Name, &account.Balance, &account.CustomerId); err != nil {
		return nil, fmt.Errorf("account id %d not found", accountID)
	}
	return account, nil
}
