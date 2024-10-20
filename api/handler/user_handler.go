package handler

import (
	db "go_simplebank/db/sqlc"
	"go_simplebank/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	UserHandlerService interface {
		createUser(ctx *gin.Context)
		loginUser(ctx *gin.Context)
	}

	userhandler struct {
		userUsecase usecase.UserUseCaseService
	}
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Fullname string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username  string    `json:"username"`
	Fullname  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(User db.User) userResponse {
	return userResponse{
		Username:  User.Username,
		Fullname:  User.FullName,
		Email:     User.Email,
		CreatedAt: User.CreatedAt,
	}
}

func NewUserHandlerService(userUsecase usecase.UserUseCaseService) UserHandlerService {
	return &userhandler{userUsecase: userUsecase}
}
func (u *userhandler) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, (err))
		return
	}
	u.userUsecase.CreateUser(ctx, usecase.CreateRequest{Username: req.Username,
		Password: req.Password,
		Fullname: req.Fullname,
		Email:    req.Email,
	})
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
}

type loginUserResponse struct {
	SessionID            uuid.UUID    `json:"session_id"`
	AccessToken          string       `json:"access_token"`
	AccessTokenExpireAt  time.Time    `json:"access_token_expire_at"`
	RefreshToken         string       `json:"refresh_token"`
	RefreshTokenExpireAt time.Time    `json:"refresh_token_expire_at"`
	User                 userResponse `json:"user"`
}

func (u *userhandler) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	u.userUsecase.Login(ctx, usecase.LoginRequest{Username: req.Username, Password: req.Password})
}
