package usecase

import (
	"context"
	"database/sql"
	"errors"
	db "go_simplebank/db/sqlc"
	"go_simplebank/token"

	"github.com/lib/pq"
)

type AccountUsecase struct {
	store db.Store
}

func NewAccountUsecase(store db.Store) *AccountUsecase {
	return &AccountUsecase{store: store}
}

type CreateAccountParam struct {
	Owner    string
	Currency string
}

func (uc *AccountUsecase) CreateAccount(ctx context.Context, authPayload *token.Payload, req CreateAccountParam) (db.Account, error) {
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := uc.store.CreateAccount(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				return db.Account{}, err
			}
		}
		return db.Account{}, err
	}

	return account, nil
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (uc *AccountUsecase) GetAccount(ctx context.Context, authPayload *token.Payload, req GetAccountRequest) (db.Account, error) {
	account, err := uc.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.Account{}, err
		}
		return db.Account{}, err
	}

	if account.Owner != authPayload.Username {
		return db.Account{}, errors.New("account does not belong to the authenticated user")
	}

	return account, nil
}

type ListAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (uc *AccountUsecase) ListAccounts(ctx context.Context, authPayload *token.Payload, req ListAccountsRequest) ([]db.Account, error) {
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := uc.store.ListAccounts(ctx, arg)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}
