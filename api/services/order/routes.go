package order

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
	return &Handler{
		db:          db,
		authHandler: authHandler,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	withAuthRouter := router.Methods("GET", "POST", "PUT", "PATCH", "DELETE").Subrouter()
	withAuthRouter.HandleFunc("/", h.authHandler.WithResourcePermissionAuth(
		h.getOrders,
		h.db,
		[]types.Resource{types.ResourceOrdersFullAccess},
	)).Methods("GET")
	withAuthRouter.HandleFunc("/pages", h.authHandler.WithResourcePermissionAuth(
		h.getOrdersPages,
		h.db,
		[]types.Resource{types.ResourceOrdersFullAccess},
	)).Methods("GET")
	withAuthRouter.HandleFunc("/{orderId}", h.authHandler.WithResourcePermissionAuth(
		h.getOrder,
		h.db,
		[]types.Resource{types.ResourceOrdersFullAccess},
	)).Methods("GET")
	withAuthRouter.HandleFunc("/{orderId}/products", h.authHandler.WithResourcePermissionAuth(
		h.getOrderProducts,
		h.db,
		[]types.Resource{types.ResourceOrdersFullAccess},
	)).Methods("GET")
	withAuthRouter.HandleFunc("/store/{storeId}", h.authHandler.WithResourcePermissionAuth(
		h.getStoreOrders,
		h.db,
		[]types.Resource{types.ResourceOrdersFullAccess},
	)).Methods("GET")
	withAuthRouter.HandleFunc("/store/{storeId}/pages", h.authHandler.WithResourcePermissionAuth(
		h.getStoreOrdersPages,
		h.db,
		[]types.Resource{types.ResourceOrdersFullAccess},
	)).Methods("GET")
	withAuthRouter.HandleFunc("/user/{userId}", h.authHandler.WithResourcePermissionAuth(
		h.getUserOrders,
		h.db,
		[]types.Resource{types.ResourceOrdersFullAccess},
	)).Methods("GET")
	withAuthRouter.HandleFunc("/user/{userId}/pages", h.authHandler.WithResourcePermissionAuth(
		h.getUserOrdersPages,
		h.db,
		[]types.Resource{types.ResourceOrdersFullAccess},
	)).Methods("GET")
	withAuthRouter.HandleFunc("/store/me/{storeId}", h.getMyStoreOrders).Methods("GET")
	withAuthRouter.HandleFunc("/store/me/{storeId}/pages", h.getMyStoreOrdersPages).Methods("GET")
	withAuthRouter.HandleFunc("/store/me/{storeId}/{orderId}", h.getMyStoreOrder).Methods("GET")
	withAuthRouter.HandleFunc("/store/me/{storeId}/{orderId}/products", h.getMyStoreOrderProducts).
		Methods("GET")
	withAuthRouter.HandleFunc("/me", h.getMyOrders).Methods("GET")
	withAuthRouter.HandleFunc("/me/pages", h.getMyOrdersPages).Methods("GET")
	withAuthRouter.HandleFunc("/me/{orderId}", h.getMyOrder).Methods("GET")
	withAuthRouter.HandleFunc("/me/{orderId}/products", h.getMyOrderProducts).Methods("GET")
	withAuthRouter.HandleFunc("/", h.authHandler.WithActionPermissionAuth(
		h.createOrder,
		h.db,
		[]types.Action{types.ActionCanCreateOrder},
	)).Methods("POST")
	withAuthRouter.HandleFunc("/complete/{orderId}", h.authHandler.WithActionPermissionAuth(
		h.completeOrderPayment,
		h.db,
		[]types.Action{types.ActionCanCompleteOrderPayment},
	)).Methods("PATCH")
	withAuthRouter.HandleFunc("/cancel/{orderId}", h.authHandler.WithActionPermissionAuth(
		h.cancelOrderPayment,
		h.db,
		[]types.Action{types.ActionCanCancelOrderPayment},
	)).Methods("PATCH")
	withAuthRouter.HandleFunc("/shipment/{orderId}", h.authHandler.WithActionPermissionAuth(
		h.updateOrderShipment,
		h.db,
		[]types.Action{types.ActionCanUpdateOrderShipment},
	)).Methods("PATCH")
	withAuthRouter.HandleFunc("/{orderId}", h.authHandler.WithActionPermissionAuth(
		h.deleteOrder,
		h.db,
		[]types.Action{types.ActionCanDeleteOrder},
	)).Methods("DELETE")
	withAuthRouter.Use(h.authHandler.WithJWTAuth(h.db))
	withAuthRouter.Use(h.authHandler.WithCSRFToken())
	withAuthRouter.Use(h.authHandler.WithVerifiedEmail(h.db))
	withAuthRouter.Use(h.authHandler.WithUnbannedProfile(h.db))
}

func (h *Handler) getOrders(w http.ResponseWriter, r *http.Request) {
	query := types.OrderSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"user":     &query.UserId,
		"store":    &query.StoreId,
		"paystat":  &query.PaymentStatus,
		"shipstat": &query.ShipmentStatus,
		"calt":     &query.CreatedAtLessThan,
		"camt":     &query.CreatedAtMoreThan,
		"p":        &page,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if query.PaymentStatus != nil {
		if !query.PaymentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderPaymentStatusEnum,
			)
			return
		}
	}

	if query.ShipmentStatus != nil {
		if !query.ShipmentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderShipmentStatusEnum,
			)
			return
		}
	}

	query.Limit = utils.Ptr(int(config.Env.MaxOrdersInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	orders, err := h.db.GetOrders(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, orders, nil)
}

func (h *Handler) getOrdersPages(w http.ResponseWriter, r *http.Request) {
	query := types.OrderSearchQuery{}

	queryMapping := map[string]any{
		"user":     &query.UserId,
		"store":    &query.StoreId,
		"paystat":  &query.PaymentStatus,
		"shipstat": &query.ShipmentStatus,
		"calt":     &query.CreatedAtLessThan,
		"camt":     &query.CreatedAtMoreThan,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if query.PaymentStatus != nil {
		if !query.PaymentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderPaymentStatusEnum,
			)
			return
		}
	}

	if query.ShipmentStatus != nil {
		if !query.ShipmentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderShipmentStatusEnum,
			)
			return
		}
	}

	count, err := h.db.GetOrdersCount(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxOrdersInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]int32{
		"pages": pageCount,
	}, nil)
}

func (h *Handler) getStoreOrders(w http.ResponseWriter, r *http.Request) {
	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query := types.OrderSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"user":     &query.UserId,
		"paystat":  &query.PaymentStatus,
		"shipstat": &query.ShipmentStatus,
		"calt":     &query.CreatedAtLessThan,
		"camt":     &query.CreatedAtMoreThan,
		"p":        &page,
	}

	queryValues := r.URL.Query()

	err = utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if query.PaymentStatus != nil {
		if !query.PaymentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderPaymentStatusEnum,
			)
			return
		}
	}

	if query.ShipmentStatus != nil {
		if !query.ShipmentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderShipmentStatusEnum,
			)
			return
		}
	}

	query.Limit = utils.Ptr(int(config.Env.MaxOrdersInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	orders, err := h.db.GetOrders(types.OrderSearchQuery{
		UserId:            query.UserId,
		StoreId:           &storeId,
		PaymentStatus:     query.PaymentStatus,
		ShipmentStatus:    query.ShipmentStatus,
		CreatedAtLessThan: query.CreatedAtLessThan,
		CreatedAtMoreThan: query.CreatedAtMoreThan,
		Limit:             query.Limit,
		Offset:            query.Offset,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, orders, nil)
}

func (h *Handler) getStoreOrdersPages(w http.ResponseWriter, r *http.Request) {
	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query := types.OrderSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"user":     &query.UserId,
		"paystat":  &query.PaymentStatus,
		"shipstat": &query.ShipmentStatus,
		"calt":     &query.CreatedAtLessThan,
		"camt":     &query.CreatedAtMoreThan,
		"p":        &page,
	}

	queryValues := r.URL.Query()

	err = utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if query.PaymentStatus != nil {
		if !query.PaymentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderPaymentStatusEnum,
			)
			return
		}
	}

	if query.ShipmentStatus != nil {
		if !query.ShipmentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderShipmentStatusEnum,
			)
			return
		}
	}

	count, err := h.db.GetOrdersCount(types.OrderSearchQuery{
		UserId:            query.UserId,
		StoreId:           &storeId,
		PaymentStatus:     query.PaymentStatus,
		ShipmentStatus:    query.ShipmentStatus,
		CreatedAtLessThan: query.CreatedAtLessThan,
		CreatedAtMoreThan: query.CreatedAtMoreThan,
		Limit:             query.Limit,
		Offset:            query.Offset,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxOrdersInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]int32{
		"pages": pageCount,
	}, nil)
}

func (h *Handler) getUserOrders(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseIntURLParam("userId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query := types.OrderSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"store":    &query.StoreId,
		"paystat":  &query.PaymentStatus,
		"shipstat": &query.ShipmentStatus,
		"calt":     &query.CreatedAtLessThan,
		"camt":     &query.CreatedAtMoreThan,
		"p":        &page,
	}

	queryValues := r.URL.Query()

	err = utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if query.PaymentStatus != nil {
		if !query.PaymentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderPaymentStatusEnum,
			)
			return
		}
	}

	if query.ShipmentStatus != nil {
		if !query.ShipmentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderShipmentStatusEnum,
			)
			return
		}
	}

	query.Limit = utils.Ptr(int(config.Env.MaxOrdersInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	orders, err := h.db.GetOrders(types.OrderSearchQuery{
		UserId:            &userId,
		StoreId:           query.StoreId,
		PaymentStatus:     query.PaymentStatus,
		ShipmentStatus:    query.ShipmentStatus,
		CreatedAtLessThan: query.CreatedAtLessThan,
		CreatedAtMoreThan: query.CreatedAtMoreThan,
		Limit:             query.Limit,
		Offset:            query.Offset,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, orders, nil)
}

func (h *Handler) getUserOrdersPages(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseIntURLParam("userId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query := types.OrderSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"store":    &query.StoreId,
		"paystat":  &query.PaymentStatus,
		"shipstat": &query.ShipmentStatus,
		"calt":     &query.CreatedAtLessThan,
		"camt":     &query.CreatedAtMoreThan,
		"p":        &page,
	}

	queryValues := r.URL.Query()

	err = utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if query.PaymentStatus != nil {
		if !query.PaymentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderPaymentStatusEnum,
			)
			return
		}
	}

	if query.ShipmentStatus != nil {
		if !query.ShipmentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderShipmentStatusEnum,
			)
			return
		}
	}

	count, err := h.db.GetOrdersCount(types.OrderSearchQuery{
		UserId:            &userId,
		StoreId:           query.StoreId,
		PaymentStatus:     query.PaymentStatus,
		ShipmentStatus:    query.ShipmentStatus,
		CreatedAtLessThan: query.CreatedAtLessThan,
		CreatedAtMoreThan: query.CreatedAtMoreThan,
		Limit:             query.Limit,
		Offset:            query.Offset,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxOrdersInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]int32{
		"pages": pageCount,
	}, nil)
}

func (h *Handler) getMyStoreOrders(w http.ResponseWriter, r *http.Request) {
	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
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

	store, err := h.db.GetStoreById(storeId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	query := types.OrderSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"user":     &query.UserId,
		"paystat":  &query.PaymentStatus,
		"shipstat": &query.ShipmentStatus,
		"calt":     &query.CreatedAtLessThan,
		"camt":     &query.CreatedAtMoreThan,
		"p":        &page,
	}

	queryValues := r.URL.Query()

	err = utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if query.PaymentStatus != nil {
		if !query.PaymentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderPaymentStatusEnum,
			)
			return
		}
	}

	if query.ShipmentStatus != nil {
		if !query.ShipmentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderShipmentStatusEnum,
			)
			return
		}
	}

	query.Limit = utils.Ptr(int(config.Env.MaxOrdersInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	orders, err := h.db.GetOrders(types.OrderSearchQuery{
		UserId:            query.UserId,
		StoreId:           &storeId,
		PaymentStatus:     query.PaymentStatus,
		ShipmentStatus:    query.ShipmentStatus,
		CreatedAtLessThan: query.CreatedAtLessThan,
		CreatedAtMoreThan: query.CreatedAtMoreThan,
		Limit:             query.Limit,
		Offset:            query.Offset,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, orders, nil)
}

func (h *Handler) getMyStoreOrdersPages(w http.ResponseWriter, r *http.Request) {
	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
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

	store, err := h.db.GetStoreById(storeId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	query := types.OrderSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"user":     &query.UserId,
		"paystat":  &query.PaymentStatus,
		"shipstat": &query.ShipmentStatus,
		"calt":     &query.CreatedAtLessThan,
		"camt":     &query.CreatedAtMoreThan,
		"p":        &page,
	}

	queryValues := r.URL.Query()

	err = utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if query.PaymentStatus != nil {
		if !query.PaymentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderPaymentStatusEnum,
			)
			return
		}
	}

	if query.ShipmentStatus != nil {
		if !query.ShipmentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderShipmentStatusEnum,
			)
			return
		}
	}

	count, err := h.db.GetOrdersCount(types.OrderSearchQuery{
		UserId:            query.UserId,
		StoreId:           &storeId,
		PaymentStatus:     query.PaymentStatus,
		ShipmentStatus:    query.ShipmentStatus,
		CreatedAtLessThan: query.CreatedAtLessThan,
		CreatedAtMoreThan: query.CreatedAtMoreThan,
		Limit:             query.Limit,
		Offset:            query.Offset,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxOrdersInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]int32{
		"pages": pageCount,
	}, nil)
}

func (h *Handler) getMyOrders(w http.ResponseWriter, r *http.Request) {
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

	query := types.OrderSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"store":    &query.StoreId,
		"paystat":  &query.PaymentStatus,
		"shipstat": &query.ShipmentStatus,
		"calt":     &query.CreatedAtLessThan,
		"camt":     &query.CreatedAtMoreThan,
		"p":        &page,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if query.PaymentStatus != nil {
		if !query.PaymentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderPaymentStatusEnum,
			)
			return
		}
	}

	if query.ShipmentStatus != nil {
		if !query.ShipmentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderShipmentStatusEnum,
			)
			return
		}
	}

	query.Limit = utils.Ptr(int(config.Env.MaxOrdersInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	orders, err := h.db.GetOrders(types.OrderSearchQuery{
		UserId:            &userId,
		StoreId:           query.StoreId,
		PaymentStatus:     query.PaymentStatus,
		ShipmentStatus:    query.ShipmentStatus,
		CreatedAtLessThan: query.CreatedAtLessThan,
		CreatedAtMoreThan: query.CreatedAtMoreThan,
		Limit:             query.Limit,
		Offset:            query.Offset,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, orders, nil)
}

func (h *Handler) getMyOrdersPages(w http.ResponseWriter, r *http.Request) {
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

	query := types.OrderSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"store":    &query.StoreId,
		"paystat":  &query.PaymentStatus,
		"shipstat": &query.ShipmentStatus,
		"calt":     &query.CreatedAtLessThan,
		"camt":     &query.CreatedAtMoreThan,
		"p":        &page,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if query.PaymentStatus != nil {
		if !query.PaymentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderPaymentStatusEnum,
			)
			return
		}
	}

	if query.ShipmentStatus != nil {
		if !query.ShipmentStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidOrderShipmentStatusEnum,
			)
			return
		}
	}

	count, err := h.db.GetOrdersCount(types.OrderSearchQuery{
		UserId:            &userId,
		StoreId:           query.StoreId,
		PaymentStatus:     query.PaymentStatus,
		ShipmentStatus:    query.ShipmentStatus,
		CreatedAtLessThan: query.CreatedAtLessThan,
		CreatedAtMoreThan: query.CreatedAtMoreThan,
		Limit:             query.Limit,
		Offset:            query.Offset,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxOrdersInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]int32{
		"pages": pageCount,
	}, nil)
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	orderId, err := utils.ParseIntURLParam("orderId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	order, err := h.db.GetOrderWithFullInfoById(orderId)
	if err != nil {
		if err == types.ErrOrderNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, order, nil)
}

func (h *Handler) getOrderProducts(w http.ResponseWriter, r *http.Request) {
	orderId, err := utils.ParseIntURLParam("orderId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	products, err := h.db.GetOrderProductVariantsInfo(orderId)
	if err != nil {
		if err == types.ErrOrderNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, products, nil)
}

func (h *Handler) getMyStoreOrder(w http.ResponseWriter, r *http.Request) {
	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
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

	store, err := h.db.GetStoreById(storeId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	orderId, err := utils.ParseIntURLParam("orderId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	isStorePartOfOrder, err := h.db.IsStoreHasParticipationInOrder(orderId, storeId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}
	if !isStorePartOfOrder {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessOrder)
		return
	}

	order, err := h.db.GetOrderWithFullInfoById(orderId)
	if err != nil {
		if err == types.ErrOrderNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, order, nil)
}

func (h *Handler) getMyStoreOrderProducts(w http.ResponseWriter, r *http.Request) {
	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
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

	store, err := h.db.GetStoreById(storeId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	orderId, err := utils.ParseIntURLParam("orderId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	isStorePartOfOrder, err := h.db.IsStoreHasParticipationInOrder(orderId, storeId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}
	if !isStorePartOfOrder {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessOrder)
		return
	}

	products, err := h.db.GetOrderProductVariantsInfo(orderId)
	if err != nil {
		if err == types.ErrOrderNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, products, nil)
}

func (h *Handler) getMyOrder(w http.ResponseWriter, r *http.Request) {
	orderId, err := utils.ParseIntURLParam("orderId", mux.Vars(r))
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

	order, err := h.db.GetOrderWithFullInfoById(orderId)
	if err != nil {
		if err == types.ErrOrderNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if order.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessOrder)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, order, nil)
}

func (h *Handler) getMyOrderProducts(w http.ResponseWriter, r *http.Request) {
	orderId, err := utils.ParseIntURLParam("orderId", mux.Vars(r))
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

	order, err := h.db.GetOrderById(orderId)
	if err != nil {
		if err == types.ErrOrderNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if order.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessOrder)
		return
	}

	products, err := h.db.GetOrderProductVariantsInfo(orderId)
	if err != nil {
		if err == types.ErrOrderNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, products, nil)
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateOrderPayload
	err := utils.ParseRequestPayload(r, &payload)
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

	address, err := h.db.GetUserAddressById(payload.ReceiverAddressId)
	if err != nil {
		if err == types.ErrUserAddressNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if address.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessAddress)
		return
	}

	createdOrder, err := h.db.CreateOrder(types.CreateOrderPayload{
		UserId:            userId,
		ArrivalDate:       payload.ArrivalDate,
		ProductVariants:   payload.ProductVariants,
		ReceiverAddressId: payload.ReceiverAddressId,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, map[string]int{"orderId": createdOrder}, nil)
}

func (h *Handler) completeOrderPayment(w http.ResponseWriter, r *http.Request) {
	// TODO: handle payment validation

	orderId, err := utils.ParseIntURLParam("orderId", mux.Vars(r))
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

	order, err := h.db.GetOrderById(orderId)
	if err != nil {
		if err == types.ErrOrderNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if order.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessOrder)
		return
	}

	err = h.db.UpdateOrderPayment(orderId, types.UpdateOrderPaymentPayload{
		Status: utils.Ptr(types.OrderPaymentStatusSuccessful),
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) cancelOrderPayment(w http.ResponseWriter, r *http.Request) {
	orderId, err := utils.ParseIntURLParam("orderId", mux.Vars(r))
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

	order, err := h.db.GetOrderById(orderId)
	if err != nil {
		if err == types.ErrOrderNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if order.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessOrder)
		return
	}

	err = h.db.UpdateOrderPayment(orderId, types.UpdateOrderPaymentPayload{
		Status: utils.Ptr(types.OrderPaymentStatusFailed),
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) updateOrderShipment(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateOrderShipmentPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	orderId, err := utils.ParseIntURLParam("orderId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.UpdateOrderShipment(orderId, types.UpdateOrderShipmentPayload{
		Status:      payload.Status,
		ArrivalDate: payload.ArrivalDate,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) deleteOrder(w http.ResponseWriter, r *http.Request) {
	orderId, err := utils.ParseIntURLParam("orderId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.DeleteOrder(orderId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}
