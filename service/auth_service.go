package service

import (
	"errors"
	"github.com/katerji/bank/db"
	"github.com/katerji/bank/db/query"
	"github.com/katerji/bank/db/queryrow"
	"github.com/katerji/bank/input"
	"github.com/katerji/bank/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func (service AuthService) Register(input input.AuthInput) (int, error) {
	password, err := hashPassword(input.Password)
	if err != nil {
		return 0, err
	}
	input.Password = password
	return db.GetDbInstance().Insert(query.InsertUserQuery, input.Email, input.Password)
}

func (service AuthService) Login(input input.AuthInput) (model.User, error) {
	result := queryrow.UserQueryRow{}
	client := db.GetDbInstance()
	row := client.QueryRow(query.GetUserByEmailQuery, input.Email)
	err := row.Scan(&result.ID, &result.Email, &result.Password)
	if err != nil {
		return model.User{}, errors.New("email does not exist")
	}

	if !validPassword(result.Password, input.Password) {
		return model.User{}, errors.New("incorrect password")
	}

	return model.User{
		ID:    result.ID,
		Email: result.Email,
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
