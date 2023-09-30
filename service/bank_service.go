package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/katerji/bank/cache"
	"github.com/katerji/bank/db"
	"github.com/katerji/bank/db/query"
	proto "github.com/katerji/bank/generated"
	"github.com/katerji/bank/utils"
	"github.com/redis/go-redis/v9"
)

const accountLockKey = "account_lock_"

type BankService struct {
	proto.UnimplementedBankServiceServer
}

func (b BankService) CreateAccount(ctx context.Context, request *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	customer := ctx.Value(utils.Customer).(*proto.Customer)
	name := request.GetName()
	if name == "" {
		return nil, errors.New("bad request")
	}

	accountID, err := db.GetDbInstance().Insert(query.CreateAccountQuery, name, customer.GetId())
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

func (b BankService) GetAccount(ctx context.Context, request *proto.GetAccountRequest) (*proto.GetAccountResponse, error) {
	customer := ctx.Value(utils.Customer).(*proto.Customer)
	accountID := int(request.GetAccountId())
	account, err := getAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("account %d not found", account.Id)
	}
	if customer.Id != account.Id {
		return nil, errors.New("unauthorized request")
	}
	return &proto.GetAccountResponse{
		Account: account,
	}, nil
}

func (b BankService) CloseAccount(ctx context.Context, request *proto.CloseAccountRequest) (*proto.GenericResponse, error) {
	customer := ctx.Value(utils.Customer).(*proto.Customer)
	accountID := int(request.GetAccountId())
	account, err := getAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("account %d not found", account.Id)
	}
	if customer.Id != account.Id {
		return nil, errors.New("unauthorized request")
	}
	if locked, err := isAccountLocked(accountID); locked || err != nil {
		return nil, errors.New("another transaction in progress")
	}
	lockAccount(accountID)
	defer unlockAccount(accountID)
	if account.Balance != 0 {
		return nil, errors.New("unable to close, withdraw balance before")
	}
	success := db.GetDbInstance().Exec(query.CloseAccountQuery, accountID)

	return &proto.GenericResponse{
		Success: success,
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
	if locked, err := isAccountLocked(accountID); locked || err != nil {
		return nil, errors.New("another transaction in progress")
	}
	lockAccount(accountID)
	defer unlockAccount(accountID)
	ok := db.GetDbInstance().Exec(query.DepositQuery, amount, accountID)

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
	if locked, err := isAccountLocked(accountID); locked || err != nil {
		return nil, errors.New("another transaction in progress")
	}
	lockAccount(accountID)
	defer unlockAccount(accountID)
	if float32(account.Balance) < amount {
		return nil, errors.New("insufficient funds")
	}
	ok := db.GetDbInstance().Exec(query.WithdrawQuery, amount, accountID)

	return &proto.GenericResponse{
		Success: ok,
	}, nil
}

func getAccount(accountID int) (*proto.Account, error) {
	account := &proto.Account{}
	row := db.GetDbInstance().FetchOne(query.FetchAccountQuery, accountID)
	if err := row.Scan(&account.Id, &account.Name, &account.Balance, &account.CustomerId); err != nil {
		return nil, fmt.Errorf("account id %d not found", accountID)
	}
	return account, nil
}

func isAccountLocked(accountID int) (bool, error) {
	key := getAccountLockRedisKey(accountID)
	locked, err := cache.GetRedisInstance().GetBool(key)
	if err == redis.Nil {
		return false, nil
	}
	return locked, err
}

func lockAccount(accountID int) {
	key := getAccountLockRedisKey(accountID)
	cache.GetRedisInstance().SetWithDefaultExpiry(key, true)
}

func unlockAccount(accountID int) {
	key := getAccountLockRedisKey(accountID)
	cache.GetRedisInstance().Delete(key)
}

func getAccountLockRedisKey(accountID int) string {
	return fmt.Sprintf("%s%d", accountLockKey, accountID)
}
