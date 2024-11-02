package handler

import (
	"go_simplebank/token"
	"go_simplebank/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	transferRequest struct {
		FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
		ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
		Amount        int64  `json:"amount" binding:"required,gt=0"`
		Currency      string `json:"currency" binding:"required,currency"`
	}

	transferHandler struct {
		transferUsecase usecase.TransferUsecase
	}

	TransferHandlerService interface {
		createTransfer(ctx *gin.Context)
	}
)

func NewTransferHandlerService(transferUsecase usecase.TransferUsecase) TransferHandlerService {
	return &transferHandler{transferUsecase: transferUsecase}
}

func (h *transferHandler) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	usecaseReq := usecase.TransferRequest{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
		Currency:      req.Currency,
	}
	transfer, err := h.transferUsecase.CreateTransfer(ctx, authPayload, usecaseReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}
