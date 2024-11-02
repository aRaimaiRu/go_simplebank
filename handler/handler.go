package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	db "go_simplebank/db/sqlc"
	"go_simplebank/token"
	"go_simplebank/usecase"
	"go_simplebank/util"

	usecase_error "go_simplebank/usecase/error"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	TokenMaker token.Maker
	Router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{store: store, TokenMaker: tokenMaker, config: config}
	gin.SetMode(gin.ReleaseMode)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	userUsecase := usecase.NewUserUseCaseService(server.store, server.store, server.TokenMaker)
	accountUsecase := usecase.NewAccountUsecase(server.store)
	AccountHandlerService := NewAccountHandlerService(*accountUsecase)
	tokenUsecase := usecase.NewTokenUsecase(server.store, server.TokenMaker, usecase.Config{AccessTokenDuration: server.config.AccessTokenDuration})
	TokenHandlerService := NewTokenHandlerService(*tokenUsecase)
	TransferUsecase := usecase.NewTransferUsecase(server.store)
	TransferHandlerService := NewTransferHandlerService(*TransferUsecase)
	UserHandlerService := NewUserHandlerService(userUsecase)
	router := gin.Default()
	router.POST("/users", UserHandlerService.createUser)
	router.POST("/users/login", UserHandlerService.loginUser)
	router.POST("/tokens/renew_access", TokenHandlerService.renewToken)

	authRoutes := router.Group("/").Use(AuthMiddleware(server.TokenMaker))

	authRoutes.POST("/accounts", AccountHandlerService.createAccount)
	authRoutes.GET("/accounts/:id", AccountHandlerService.getAccount)
	authRoutes.GET("/accounts", AccountHandlerService.listAccounts)
	authRoutes.POST("/transfers", TransferHandlerService.createTransfer)
	server.Router = router
}

func (server *Server) Start(addr string) error {
	return server.Router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func errorResponseFromUsecase(err usecase_error.IUseCaseError) gin.H {
	return gin.H{"error": err.Message()}
}

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSuportedCurrency(currency)
	}
	return false
}
