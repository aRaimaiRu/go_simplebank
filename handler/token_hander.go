package handler

import (
	"go_simplebank/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	TokenHandlerService interface {
		renewToken(ctx *gin.Context)
	}

	tokenHandler struct {
		tokenUsecase usecase.TokenUsecase
	}

	renewTokenRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
)

func NewTokenHandlerService(tokenUsecase usecase.TokenUsecase) TokenHandlerService {
	return &tokenHandler{tokenUsecase: tokenUsecase}
}

func (h *tokenHandler) renewToken(ctx *gin.Context) {
	var req renewTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response, err := h.tokenUsecase.RenewToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}
