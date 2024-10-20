package gapi

import (
	"fmt"
	db "go_simplebank/db/sqlc"
	"go_simplebank/token"
	"go_simplebank/util"

	"go_simplebank/pb"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{store: store, tokenMaker: tokenMaker, config: config}
	gin.SetMode(gin.ReleaseMode)

	return server, nil
}

// func (server *Server) CreateUsers(ctx *context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
// 	hashedPassword, err := util.HashPassword(req.Password)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	arg := db.CreateUserParams{
// 		Username:       req.Username,
// 		HashedPassword: hashedPassword,
// 		FullName:       req.Fullname,
// 		Email:          req.Email,
// 	}
// 	return nil, nil
// }

// func errorResponse(err error) gin.H {
// 	return gin.H{"error": err.Error()}
// }
