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
	withAuthRouter.HandleFunc("", h.authHandler.WithResourcePermissionAuth(
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
	withAuthRouter.HandleFunc("", h.authHandler.WithActionPermissionAuth(
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

// getOrders godoc
// @Summary      Get orders list
// @Description  Retrieves a paginated list of orders with optional filtering. Requires orders full access permission.
// @Tags         order
// @Produce      json
// @Param        user      query     int     false  "Filter by user ID"
// @Param        store     query     int     false  "Filter by store ID"
// @Param        paystat   query     string  false  "Filter by payment status"
// @Param        shipstat  query     string  false  "Filter by shipment status"
// @Param        calt      query     string  false  "Filter by created before date (RFC3339)"
// @Param        camt      query     string  false  "Filter by created after date (RFC3339)"
// @Param        p         query     int     false  "Page number (default: 1)"
// @Success      200       {array}   types.Order
// @Failure      400       {object}  types.HTTPError
// @Failure      401       {object}  types.HTTPError
// @Failure      403       {object}  types.HTTPError
// @Failure      500       {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order [get]
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

// getOrdersPages godoc
// @Summary      Get total orders pages count
// @Description  Calculates the total number of pages available for orders listing based on filters. Requires orders full access permission.
// @Tags         order
// @Produce      json
// @Param        user      query     int     false  "Filter by user ID"
// @Param        store     query     int     false  "Filter by store ID"
// @Param        paystat   query     string  false  "Filter by payment status"
// @Param        shipstat  query     string  false  "Filter by shipment status"
// @Param        calt      query     string  false  "Filter by created before date (RFC3339)"
// @Param        camt      query     string  false  "Filter by created after date (RFC3339)"
// @Success      200       {object}  types.TotalPageCountResponse
// @Failure      400       {object}  types.HTTPError
// @Failure      401       {object}  types.HTTPError
// @Failure      403       {object}  types.HTTPError
// @Failure      500       {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/pages [get]
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

	utils.WriteJSONInResponse(w, http.StatusOK, types.TotalPageCountResponse{
		Pages: pageCount,
	}, nil)
}

// getStoreOrders godoc
// @Summary      Get store's orders
// @Description  Retrieves a paginated list of orders for a specific store with optional filtering. Requires orders full access permission.
// @Tags         order
// @Produce      json
// @Param        storeId   path      int     true   "Store ID"
// @Param        user      query     int     false  "Filter by user ID"
// @Param        paystat   query     string  false  "Filter by payment status"
// @Param        shipstat  query     string  false  "Filter by shipment status"
// @Param        calt      query     string  false  "Filter by created before date (RFC3339)"
// @Param        camt      query     string  false  "Filter by created after date (RFC3339)"
// @Param        p         query     int     false  "Page number (default: 1)"
// @Success      200       {array}   types.Order
// @Failure      400       {object}  types.HTTPError
// @Failure      401       {object}  types.HTTPError
// @Failure      403       {object}  types.HTTPError
// @Failure      500       {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/store/{storeId} [get]
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

// getStoreOrdersPages godoc
// @Summary      Get store's orders pages count
// @Description  Calculates the total number of pages available for a store's orders listing. Requires orders full access permission.
// @Tags         order
// @Produce      json
// @Param        storeId   path      int     true   "Store ID"
// @Param        user      query     int     false  "Filter by user ID"
// @Param        paystat   query     string  false  "Filter by payment status"
// @Param        shipstat  query     string  false  "Filter by shipment status"
// @Param        calt      query     string  false  "Filter by created before date (RFC3339)"
// @Param        camt      query     string  false  "Filter by created after date (RFC3339)"
// @Success      200       {object}  types.TotalPageCountResponse
// @Failure      400       {object}  types.HTTPError
// @Failure      401       {object}  types.HTTPError
// @Failure      403       {object}  types.HTTPError
// @Failure      500       {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/store/{storeId}/pages [get]
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

	utils.WriteJSONInResponse(w, http.StatusOK, types.TotalPageCountResponse{
		Pages: pageCount,
	}, nil)
}

// getUserOrders godoc
// @Summary      Get user's orders
// @Description  Retrieves a paginated list of orders for a specific user with optional filtering. Requires orders full access permission.
// @Tags         order
// @Produce      json
// @Param        userId    path      int     true   "User ID"
// @Param        store     query     int     false  "Filter by store ID"
// @Param        paystat   query     string  false  "Filter by payment status"
// @Param        shipstat  query     string  false  "Filter by shipment status"
// @Param        calt      query     string  false  "Filter by created before date (RFC3339)"
// @Param        camt      query     string  false  "Filter by created after date (RFC3339)"
// @Param        p         query     int     false  "Page number (default: 1)"
// @Success      200       {array}   types.Order
// @Failure      400       {object}  types.HTTPError
// @Failure      401       {object}  types.HTTPError
// @Failure      403       {object}  types.HTTPError
// @Failure      500       {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/user/{userId} [get]
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

// getUserOrdersPages godoc
// @Summary      Get user's orders pages count
// @Description  Calculates the total number of pages available for a user's orders listing. Requires orders full access permission.
// @Tags         order
// @Produce      json
// @Param        userId    path      int     true   "User ID"
// @Param        store     query     int     false  "Filter by store ID"
// @Param        paystat   query     string  false  "Filter by payment status"
// @Param        shipstat  query     string  false  "Filter by shipment status"
// @Param        calt      query     string  false  "Filter by created before date (RFC3339)"
// @Param        camt      query     string  false  "Filter by created after date (RFC3339)"
// @Success      200       {object}  types.TotalPageCountResponse
// @Failure      400       {object}  types.HTTPError
// @Failure      401       {object}  types.HTTPError
// @Failure      403       {object}  types.HTTPError
// @Failure      500       {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/user/{userId}/pages [get]
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

	utils.WriteJSONInResponse(w, http.StatusOK, types.TotalPageCountResponse{
		Pages: pageCount,
	}, nil)
}

// getMyStoreOrders godoc
// @Summary      Get current user's store orders
// @Description  Retrieves a paginated list of orders for the current user's store with optional filtering.
// @Tags         order
// @Produce      json
// @Param        storeId   path      int     true   "Store ID"
// @Param        user      query     int     false  "Filter by user ID"
// @Param        paystat   query     string  false  "Filter by payment status"
// @Param        shipstat  query     string  false  "Filter by shipment status"
// @Param        calt      query     string  false  "Filter by created before date (RFC3339)"
// @Param        camt      query     string  false  "Filter by created after date (RFC3339)"
// @Param        p         query     int     false  "Page number (default: 1)"
// @Success      200       {array}   types.Order
// @Failure      400       {object}  types.HTTPError
// @Failure      401       {object}  types.HTTPError
// @Failure      403       {object}  types.HTTPError
// @Failure      404       {object}  types.HTTPError
// @Failure      500       {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/store/me/{storeId} [get]
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

// getMyStoreOrdersPages godoc
// @Summary      Get current user's store orders pages count
// @Description  Calculates the total number of pages available for the current user's store orders listing.
// @Tags         order
// @Produce      json
// @Param        storeId   path      int     true   "Store ID"
// @Param        user      query     int     false  "Filter by user ID"
// @Param        paystat   query     string  false  "Filter by payment status"
// @Param        shipstat  query     string  false  "Filter by shipment status"
// @Param        calt      query     string  false  "Filter by created before date (RFC3339)"
// @Param        camt      query     string  false  "Filter by created after date (RFC3339)"
// @Success      200       {object}  types.TotalPageCountResponse
// @Failure      400       {object}  types.HTTPError
// @Failure      401       {object}  types.HTTPError
// @Failure      403       {object}  types.HTTPError
// @Failure      404       {object}  types.HTTPError
// @Failure      500       {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/store/me/{storeId}/pages [get]
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

	utils.WriteJSONInResponse(w, http.StatusOK, types.TotalPageCountResponse{
		Pages: pageCount,
	}, nil)
}

// getMyOrders godoc
// @Summary      Get current user's orders
// @Description  Retrieves a paginated list of orders for the current user with optional filtering.
// @Tags         order
// @Produce      json
// @Param        store     query     int     false  "Filter by store ID"
// @Param        paystat   query     string  false  "Filter by payment status"
// @Param        shipstat  query     string  false  "Filter by shipment status"
// @Param        calt      query     string  false  "Filter by created before date (RFC3339)"
// @Param        camt      query     string  false  "Filter by created after date (RFC3339)"
// @Param        p         query     int     false  "Page number (default: 1)"
// @Success      200       {array}   types.Order
// @Failure      400       {object}  types.HTTPError
// @Failure      401       {object}  types.HTTPError
// @Failure      500       {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/me [get]
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

// getMyOrdersPages godoc
// @Summary      Get current user's orders pages count
// @Description  Calculates the total number of pages available for the current user's orders listing.
// @Tags         order
// @Produce      json
// @Param        store     query     int     false  "Filter by store ID"
// @Param        paystat   query     string  false  "Filter by payment status"
// @Param        shipstat  query     string  false  "Filter by shipment status"
// @Param        calt      query     string  false  "Filter by created before date (RFC3339)"
// @Param        camt      query     string  false  "Filter by created after date (RFC3339)"
// @Success      200       {object}  types.TotalPageCountResponse
// @Failure      400       {object}  types.HTTPError
// @Failure      401       {object}  types.HTTPError
// @Failure      500       {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/me/pages [get]
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

	utils.WriteJSONInResponse(w, http.StatusOK, types.TotalPageCountResponse{
		Pages: pageCount,
	}, nil)
}

// getOrder godoc
// @Summary      Get order details
// @Description  Retrieves full details for a specific order. Requires orders full access permission.
// @Tags         order
// @Produce      json
// @Param        orderId  path      int  true  "Order ID"
// @Success      200      {object}  types.Order
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/{orderId} [get]
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

// getOrderProducts godoc
// @Summary      Get order products
// @Description  Retrieves product variants for a specific order. Requires orders full access permission.
// @Tags         order
// @Produce      json
// @Param        orderId  path      int  true  "Order ID"
// @Success      200      {array}   types.OrderProductVariantInfo
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/{orderId}/products [get]
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

// getMyStoreOrder godoc
// @Summary      Get current user's store order
// @Description  Retrieves full details for a specific order belonging to the current user's store.
// @Tags         order
// @Produce      json
// @Param        storeId  path      int  true  "Store ID"
// @Param        orderId  path      int  true  "Order ID"
// @Success      200      {object}  types.Order
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/store/me/{storeId}/{orderId} [get]
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

// getMyStoreOrderProducts godoc
// @Summary      Get current user's store order products
// @Description  Retrieves product variants for a specific order belonging to the current user's store.
// @Tags         order
// @Produce      json
// @Param        storeId  path      int  true  "Store ID"
// @Param        orderId  path      int  true  "Order ID"
// @Success      200      {array}   types.OrderProductVariantInfo
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/store/me/{storeId}/{orderId}/products [get]
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

// getMyOrder godoc
// @Summary      Get current user's order
// @Description  Retrieves full details for a specific order belonging to the current user.
// @Tags         order
// @Produce      json
// @Param        orderId  path      int  true  "Order ID"
// @Success      200      {object}  types.Order
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/me/{orderId} [get]
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

// getMyOrderProducts godoc
// @Summary      Get current user's order products
// @Description  Retrieves product variants for a specific order belonging to the current user.
// @Tags         order
// @Produce      json
// @Param        orderId  path      int  true  "Order ID"
// @Success      200      {array}   types.OrderProductVariantInfo
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/me/{orderId}/products [get]
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

// createOrder godoc
// @Summary      Create a new order
// @Description  Creates a new order with the provided details. Requires create order permission.
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        order  body      types.CreateOrderPayload  true  "Order details"
// @Success      201    {object}  types.NewOrderResponse
// @Failure      400    {object}  types.HTTPError
// @Failure      401    {object}  types.HTTPError
// @Failure      403    {object}  types.HTTPError
// @Failure      404    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order [post]
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

	res := types.NewOrderResponse{
		OrderId: createdOrder,
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, res, nil)
}

// completeOrderPayment godoc
// @Summary      Complete order payment
// @Description  Marks an order's payment as completed. Requires complete order payment permission.
// @Tags         order
// @Produce      json
// @Param        orderId  path      int  true  "Order ID"
// @Success      200      "Order payment completed"
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/complete/{orderId} [patch]
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

// cancelOrderPayment godoc
// @Summary      Cancel order payment
// @Description  Marks an order's payment as cancelled. Requires cancel order payment permission.
// @Tags         order
// @Produce      json
// @Param        orderId  path      int  true  "Order ID"
// @Success      200			"Order payment cancelled"
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/cancel/{orderId} [patch]
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

// updateOrderShipment godoc
// @Summary      Update order shipment
// @Description  Updates shipment details for an order. Requires update order shipment permission.
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        orderId  path      int                               true  "Order ID"
// @Param        shipment body      types.UpdateOrderShipmentPayload  true  "Shipment details"
// @Success      200			"Order shipment updated"
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/shipment/{orderId} [patch]
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

// deleteOrder godoc
// @Summary      Delete an order
// @Description  Permanently deletes an order. Requires delete order permission.
// @Tags         order
// @Produce      json
// @Param        orderId  path      int  true  "Order ID"
// @Success      200    	"Order deleted"
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /order/{orderId} [delete]
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
