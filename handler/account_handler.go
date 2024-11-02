package handler

import (
	"database/sql"
	"go_simplebank/token"
	"go_simplebank/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	AccountHandlerService interface {
		createAccount(ctx *gin.Context)
		getAccount(ctx *gin.Context)
		listAccounts(ctx *gin.Context)
	}

	accounthandler struct {
		accountUsecase usecase.AccountUsecase
	}
)

type CreateAccountParam struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func NewAccountHandlerService(accountUsecase usecase.AccountUsecase) AccountHandlerService {
	return &accounthandler{accountUsecase: accountUsecase}
}

func (h *accounthandler) createAccount(ctx *gin.Context) {
	var req CreateAccountParam
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	usecaseReq := usecase.CreateAccountParam{
		Owner:    req.Owner,
		Currency: req.Currency,
	}
	account, err := h.accountUsecase.CreateAccount(ctx, authPayload, usecaseReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (h *accounthandler) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	usecaseReq := usecase.GetAccountRequest{ID: req.ID}
	account, err := h.accountUsecase.GetAccount(ctx, authPayload, usecaseReq)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (h *accounthandler) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	usecaseReq := usecase.ListAccountsRequest{
		PageID:   req.PageID,
		PageSize: req.PageSize,
	}
	accounts, err := h.accountUsecase.ListAccounts(ctx, authPayload, usecaseReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
