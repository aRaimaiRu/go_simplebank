package usecase

import (
	"context"
	"database/sql"
	"fmt"
	db "go_simplebank/db/sqlc"
	"go_simplebank/token"
)

type TransferUsecase struct {
	store db.Store
}

func NewTransferUsecase(store db.Store) *TransferUsecase {
	return &TransferUsecase{store: store}
}

type TransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (uc *TransferUsecase) CreateTransfer(ctx context.Context, authPayload *token.Payload, req TransferRequest) (db.TransferTxResult, error) {
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	fromAccount, valid := uc.validAccount(ctx, arg.FromAccountID, req.Currency)
	if !valid {
		return db.TransferTxResult{}, fmt.Errorf("invalid from account")
	}

	if fromAccount.Owner != authPayload.Username {
		return db.TransferTxResult{}, fmt.Errorf("account [%d] does not belong to user [%s]", arg.FromAccountID, authPayload.Username)
	}

	_, valid = uc.validAccount(ctx, arg.ToAccountID, req.Currency)
	if !valid {
		return db.TransferTxResult{}, fmt.Errorf("invalid to account")
	}

	transfer, err := uc.store.TransferTx(ctx, arg)
	if err != nil {
		return db.TransferTxResult{}, err
	}

	return transfer, nil
}

func (uc *TransferUsecase) validAccount(ctx context.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := uc.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return account, false
		}
		return account, false
	}

	if account.Currency != currency {
		return account, false
	}
	return account, true
}
