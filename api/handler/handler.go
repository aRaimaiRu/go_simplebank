package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	db "go_simplebank/db/sqlc"
	"go_simplebank/token"
	"go_simplebank/usecase"
	"go_simplebank/util"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	Router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker, config: config}
	gin.SetMode(gin.ReleaseMode)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	userUsecase := usecase.NewUserUseCaseService(server.store, server.store, server.tokenMaker)
	UserHandlerService := NewUserHandlerService(userUsecase)
	router := gin.Default()
	router.POST("/users", UserHandlerService.createUser)
	router.POST("/users/login", UserHandlerService.loginUser)
	// router.POST("/tokens/renew_access", server.renewToken)

	// authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// authRoutes.POST("/accounts", server.createAccount)
	// authRoutes.GET("/accounts/:id", server.getAccount)
	// authRoutes.GET("/accounts", server.listAccounts)
	// authRoutes.POST("/transfers", server.createTransfer)
	server.Router = router
}

func (server *Server) Start(addr string) error {
	return server.Router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSuportedCurrency(currency)
	}
	return false
}
