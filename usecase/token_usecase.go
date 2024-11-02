package usecase

import (
	"context"
	"database/sql"
	"fmt"
	db "go_simplebank/db/sqlc"
	"go_simplebank/token"
	"time"

	"github.com/google/uuid"
)

type TokenUsecase struct {
	store      db.Store
	tokenMaker token.Maker
	config     Config
}

type Config struct {
	AccessTokenDuration time.Duration
}

type RenewTokenResponse struct {
	SessionID            uuid.UUID    `json:"session_id"`
	AccessToken          string       `json:"access_token"`
	AccessTokenExpireAt  time.Time    `json:"access_token_expire_at"`
	RefreshToken         string       `json:"refresh_token"`
	RefreshTokenExpireAt time.Time    `json:"refresh_token_expire_at"`
	User                 userResponse `json:"user"`
}

func NewTokenUsecase(store db.Store, tokenMaker token.Maker, config Config) *TokenUsecase {
	return &TokenUsecase{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
}

func (uc *TokenUsecase) RenewToken(ctx context.Context, refreshToken string) (RenewTokenResponse, error) {
	refreshPayload, err := uc.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		return RenewTokenResponse{}, err
	}

	session, err := uc.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return RenewTokenResponse{}, err
		}
		return RenewTokenResponse{}, err
	}

	if session.IsBlocked {
		return RenewTokenResponse{}, fmt.Errorf("session is blocked")
	}

	if session.Username != refreshPayload.Username {
		return RenewTokenResponse{}, fmt.Errorf("incorrect session user")
	}

	if time.Now().After(session.ExpiresAt) {
		return RenewTokenResponse{}, fmt.Errorf("session is expired")
	}

	accessToken, accessPayload, err := uc.tokenMaker.CreateToken(refreshPayload.Username, uc.config.AccessTokenDuration)
	if err != nil {
		return RenewTokenResponse{}, err
	}

	response := RenewTokenResponse{
		SessionID:            session.ID,
		AccessToken:          accessToken,
		AccessTokenExpireAt:  accessPayload.ExpiredAt,
		RefreshToken:         refreshToken,
		RefreshTokenExpireAt: session.ExpiresAt,
		User:                 userResponse{Username: session.Username},
	}

	return response, nil
}
