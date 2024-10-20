package usecase

import (
	"context"
	"database/sql"
	db "go_simplebank/db/sqlc"
	"go_simplebank/token"
	usecase_error "go_simplebank/usecase/error"
	"go_simplebank/util"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type UserUseCaseService interface {
	CreateUser(ctx context.Context, req CreateRequest) (db.User, usecase_error.IUseCaseError)
	Login(ctx context.Context, req LoginRequest) (loginUserResponse, usecase_error.IUseCaseError)
}

type userUseCaseService struct {
	config            util.Config
	userRepository    db.UserRepository
	sessionRepository db.SessionRepository
	tokenMaker        token.Maker
}

func NewUserUseCaseService(userRepository db.UserRepository, sessionRepository db.SessionRepository, tokenMaker token.Maker) UserUseCaseService {
	return &userUseCaseService{
		tokenMaker:        tokenMaker,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
	}
}

type CreateRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Fullname string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
}

type userResponse struct {
	Username  string    `json:"username"`
	Fullname  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type loginUserResponse struct {
	SessionID            uuid.UUID    `json:"session_id"`
	AccessToken          string       `json:"access_token"`
	AccessTokenExpireAt  time.Time    `json:"access_token_expire_at"`
	RefreshToken         string       `json:"refresh_token"`
	RefreshTokenExpireAt time.Time    `json:"refresh_token_expire_at"`
	User                 userResponse `json:"user"`
}

func newUserResponse(User db.User) userResponse {
	return userResponse{
		Username:  User.Username,
		Fullname:  User.FullName,
		Email:     User.Email,
		CreatedAt: User.CreatedAt,
	}
}

func (s *userUseCaseService) CreateUser(ctx context.Context, req CreateRequest) (db.User, usecase_error.IUseCaseError) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return db.User{}, usecase_error.NewUseCaseError(http.StatusInternalServerError, err.Error())
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.Fullname,
		Email:          req.Email,
	}

	user, err := s.userRepository.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return db.User{}, usecase_error.NewUseCaseError(http.StatusForbidden, err.Error())
		}
		return db.User{}, usecase_error.NewUseCaseError(http.StatusInternalServerError, err.Error())
	}
	return user, nil
}

func (s *userUseCaseService) Login(ctx context.Context, req LoginRequest) (loginUserResponse, usecase_error.IUseCaseError) {
	user, err := s.userRepository.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return loginUserResponse{}, usecase_error.NewUseCaseError(http.StatusNotFound, "user not found")
		}
		return loginUserResponse{}, usecase_error.NewUseCaseError(http.StatusInternalServerError, err.Error())
	}

	err = util.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		return loginUserResponse{}, usecase_error.NewUseCaseError(http.StatusUnauthorized, "incorrect password")
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		return loginUserResponse{}, usecase_error.NewUseCaseError(http.StatusInternalServerError, err.Error())
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.RefreshTokenDuration)
	if err != nil {
		return loginUserResponse{}, usecase_error.NewUseCaseError(http.StatusInternalServerError, err.Error())
	}

	session, err := s.sessionRepository.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "", // You need to pass the user agent from the request
		ClientIp:     "", // You need to pass the client IP from the request
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return loginUserResponse{}, usecase_error.NewUseCaseError(http.StatusInternalServerError, err.Error())
	}

	response := loginUserResponse{
		SessionID:            session.ID,
		AccessToken:          accessToken,
		AccessTokenExpireAt:  accessPayload.ExpiredAt,
		RefreshToken:         refreshToken,
		RefreshTokenExpireAt: refreshPayload.ExpiredAt,
		User:                 newUserResponse(user),
	}
	return response, nil
}
