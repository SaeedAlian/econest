package wallet

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/SaeedAlian/econest/api/config"
	db_manager "github.com/SaeedAlian/econest/api/db/manager"
	"github.com/SaeedAlian/econest/api/services/auth"
	"github.com/SaeedAlian/econest/api/types"
	"github.com/SaeedAlian/econest/api/utils"
)

type Handler struct {
	db          *db_manager.Manager
	authHandler *auth.AuthHandler
}

func NewHandler(
	db *db_manager.Manager,
	authHandler *auth.AuthHandler,
) *Handler {
	return &Handler{db: db, authHandler: authHandler}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	withAuthRouter := router.Methods("GET", "POST", "PATCH", "DELETE").Subrouter()
	withAuthRouter.HandleFunc("/me", h.getMyWallet).Methods("GET")
	withAuthRouter.HandleFunc("/me/transaction", h.getMyTransactions).Methods("GET")
	withAuthRouter.HandleFunc("/me/transaction/pages", h.getMyTransactionsPages).Methods("GET")
	withAuthRouter.HandleFunc("/me/transaction/{txId}", h.getMyTransaction).Methods("GET")
	withAuthRouter.HandleFunc("/user/{userId}", h.authHandler.WithResourcePermissionAuth(
		h.getUserWallet,
		h.db,
		[]types.Resource{types.ResourceWalletTransactionsFullAccess},
	)).Methods("GET")
	withAuthRouter.HandleFunc("/user/{userId}/transaction", h.authHandler.WithResourcePermissionAuth(
		h.getUserTransactions,
		h.db,
		[]types.Resource{types.ResourceWalletTransactionsFullAccess},
	)).
		Methods("GET")
	withAuthRouter.HandleFunc("/user/{userId}/transaction/pages", h.authHandler.WithResourcePermissionAuth(
		h.getUserTransactionsPages,
		h.db,
		[]types.Resource{types.ResourceWalletTransactionsFullAccess},
	)).
		Methods("GET")
	withAuthRouter.HandleFunc("/user/{userId}/transaction/{txId}", h.authHandler.WithResourcePermissionAuth(
		h.getUserTransaction,
		h.db,
		[]types.Resource{types.ResourceWalletTransactionsFullAccess},
	)).
		Methods("GET")
	withAuthRouter.Use(h.authHandler.WithJWTAuth(h.db))
	withAuthRouter.Use(h.authHandler.WithCSRFToken())
	withAuthRouter.Use(h.authHandler.WithVerifiedEmail(h.db))
	withAuthRouter.Use(h.authHandler.WithUnbannedProfile(h.db))

	withdrawRouter := withAuthRouter.PathPrefix("/withdraw").Subrouter()
	withdrawRouter.HandleFunc("", h.createWithdrawTransaction).Methods("POST")
	withdrawRouter.HandleFunc("/complete/{txId}", h.authHandler.WithActionPermissionAuth(
		h.completeWithdrawTransaction,
		h.db,
		[]types.Action{types.ActionCanApproveWithdrawTransaction},
	)).
		Methods("PATCH")
	withdrawRouter.HandleFunc("/cancel/{txId}", h.authHandler.WithActionPermissionAuth(
		h.cancelWithdrawTransaction,
		h.db,
		[]types.Action{types.ActionCanCancelWithdrawTransaction},
	)).
		Methods("PATCH")

	depositRouter := withAuthRouter.PathPrefix("/deposit").Subrouter()
	depositRouter.HandleFunc("", h.createDepositTransaction).Methods("POST")
	depositRouter.HandleFunc("/complete/{txId}", h.completeDepositTransaction).Methods("PATCH")
	depositRouter.HandleFunc("/cancel/{txId}", h.cancelDepositTransaction).Methods("PATCH")
}

func (h *Handler) createDepositTransaction(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateWalletTransactionPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	userId := ctx.Value("userId")

	if userId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	wallet, err := h.db.GetUserWallet(userId.(int))
	if err != nil {
		if err == types.ErrWalletNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	createdTx, err := h.db.CreateWalletTransaction(types.CreateWalletTransactionPayload{
		Amount:   payload.Amount,
		TxType:   types.TransactionTypeDeposit,
		WalletId: wallet.Id,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, map[string]int{"txId": createdTx}, nil)
}

func (h *Handler) createWithdrawTransaction(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateWalletTransactionPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	userId := ctx.Value("userId")

	if userId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	wallet, err := h.db.GetUserWallet(userId.(int))
	if err != nil {
		if err == types.ErrWalletNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	createdTx, err := h.db.CreateWalletTransaction(types.CreateWalletTransactionPayload{
		Amount:   payload.Amount,
		TxType:   types.TransactionTypeWithdraw,
		WalletId: wallet.Id,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, map[string]int{"txId": createdTx}, nil)
}

func (h *Handler) getMyWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	wallet, err := h.db.GetUserWallet(userId)
	if err != nil {
		if err == types.ErrWalletNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, wallet, nil)
}

func (h *Handler) getMyTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	query := types.WalletTransactionSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"typ":  &query.TxType,
		"stat": &query.Status,
		"aftd": &query.AfterDate,
		"befd": &query.BeforeDate,
		"p":    &page,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query.Limit = utils.Ptr(int(config.Env.MaxWalletTransactionsInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	txs, err := h.db.GetWalletTransactions(types.WalletTransactionSearchQuery{
		Status:     query.Status,
		TxType:     query.TxType,
		BeforeDate: query.BeforeDate,
		AfterDate:  query.AfterDate,
		UserId:     &userId,
		Limit:      query.Limit,
		Offset:     query.Offset,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, txs, nil)
}

func (h *Handler) getMyTransactionsPages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	query := types.WalletTransactionSearchQuery{}

	queryMapping := map[string]any{
		"typ":  &query.TxType,
		"stat": &query.Status,
		"aftd": &query.AfterDate,
		"befd": &query.BeforeDate,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	count, err := h.db.GetWalletTransactionsCount(types.WalletTransactionSearchQuery{
		Status:     query.Status,
		TxType:     query.TxType,
		BeforeDate: query.BeforeDate,
		AfterDate:  query.AfterDate,
		UserId:     &userId,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxWalletTransactionsInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]int32{
		"pages": pageCount,
	}, nil)
}

func (h *Handler) getMyTransaction(w http.ResponseWriter, r *http.Request) {
	txId, err := utils.ParseIntURLParam("txId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	wallet, err := h.db.GetUserWallet(userId)
	if err != nil {
		if err == types.ErrWalletNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	tx, err := h.db.GetWalletTransactionById(txId)
	if err != nil {
		if err == types.ErrWalletTransactionNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if tx.WalletId != wallet.Id {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessWalletTransaction)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, tx, nil)
}

func (h *Handler) getUserWallet(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseIntURLParam("userId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	wallet, err := h.db.GetUserWallet(userId)
	if err != nil {
		if err == types.ErrWalletNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, wallet, nil)
}

func (h *Handler) getUserTransactions(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseIntURLParam("userId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query := types.WalletTransactionSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"typ":  &query.TxType,
		"stat": &query.Status,
		"aftd": &query.AfterDate,
		"befd": &query.BeforeDate,
		"p":    &page,
	}

	queryValues := r.URL.Query()

	err = utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query.Limit = utils.Ptr(int(config.Env.MaxWalletTransactionsInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	txs, err := h.db.GetWalletTransactions(types.WalletTransactionSearchQuery{
		Status:     query.Status,
		TxType:     query.TxType,
		BeforeDate: query.BeforeDate,
		AfterDate:  query.AfterDate,
		UserId:     &userId,
		Limit:      query.Limit,
		Offset:     query.Offset,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, txs, nil)
}

func (h *Handler) getUserTransactionsPages(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseIntURLParam("userId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query := types.WalletTransactionSearchQuery{}

	queryMapping := map[string]any{
		"typ":  &query.TxType,
		"stat": &query.Status,
		"aftd": &query.AfterDate,
		"befd": &query.BeforeDate,
	}

	queryValues := r.URL.Query()

	err = utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	count, err := h.db.GetWalletTransactionsCount(types.WalletTransactionSearchQuery{
		Status:     query.Status,
		TxType:     query.TxType,
		BeforeDate: query.BeforeDate,
		AfterDate:  query.AfterDate,
		UserId:     &userId,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxWalletTransactionsInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]int32{
		"pages": pageCount,
	}, nil)
}

func (h *Handler) getUserTransaction(w http.ResponseWriter, r *http.Request) {
	txId, err := utils.ParseIntURLParam("txId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	userId, err := utils.ParseIntURLParam("userId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	wallet, err := h.db.GetUserWallet(userId)
	if err != nil {
		if err == types.ErrWalletNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	tx, err := h.db.GetWalletTransactionById(txId)
	if err != nil {
		if err == types.ErrWalletTransactionNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if tx.WalletId != wallet.Id {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrTransactionIsNotForWallet)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, tx, nil)
}

func (h *Handler) completeDepositTransaction(w http.ResponseWriter, r *http.Request) {
	// TODO: handle payment validation

	txId, err := utils.ParseIntURLParam("txId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	wallet, err := h.db.GetUserWallet(userId)
	if err != nil {
		if err == types.ErrWalletNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	tx, err := h.db.GetWalletTransactionById(txId)
	if err != nil {
		if err == types.ErrWalletTransactionNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if tx.WalletId != wallet.Id || wallet.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessWalletTransaction)
		return
	}

	if tx.TxType != types.TransactionTypeDeposit {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrInvalidTransactionTypeEnum)
		return
	}

	err = h.db.UpdateWalletTransaction(tx.WalletId, txId,
		types.UpdateWalletTransactionPayload{
			Status: utils.Ptr(types.TransactionStatusSuccessful),
		},
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) completeWithdrawTransaction(w http.ResponseWriter, r *http.Request) {
	// TODO: handle payment validation

	txId, err := utils.ParseIntURLParam("txId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := h.db.GetWalletTransactionById(txId)
	if err != nil {
		if err == types.ErrWalletTransactionNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if tx.TxType != types.TransactionTypeWithdraw {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrInvalidTransactionTypeEnum)
		return
	}

	err = h.db.UpdateWalletTransaction(tx.WalletId, txId,
		types.UpdateWalletTransactionPayload{
			Status: utils.Ptr(types.TransactionStatusSuccessful),
		},
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) cancelDepositTransaction(w http.ResponseWriter, r *http.Request) {
	txId, err := utils.ParseIntURLParam("txId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	wallet, err := h.db.GetUserWallet(userId)
	if err != nil {
		if err == types.ErrWalletNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	tx, err := h.db.GetWalletTransactionById(txId)
	if err != nil {
		if err == types.ErrWalletTransactionNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if tx.WalletId != wallet.Id || wallet.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessWalletTransaction)
		return
	}

	if tx.TxType != types.TransactionTypeDeposit {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrInvalidTransactionTypeEnum)
		return
	}

	err = h.db.UpdateWalletTransaction(tx.WalletId, txId,
		types.UpdateWalletTransactionPayload{
			Status: utils.Ptr(types.TransactionStatusFailed),
		},
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) cancelWithdrawTransaction(w http.ResponseWriter, r *http.Request) {
	txId, err := utils.ParseIntURLParam("txId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	tx, err := h.db.GetWalletTransactionById(txId)
	if err != nil {
		if err == types.ErrWalletTransactionNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if tx.TxType != types.TransactionTypeWithdraw {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrInvalidTransactionTypeEnum)
		return
	}

	err = h.db.UpdateWalletTransaction(tx.WalletId, txId,
		types.UpdateWalletTransactionPayload{
			Status: utils.Ptr(types.TransactionStatusFailed),
		},
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}
