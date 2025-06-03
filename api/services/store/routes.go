package store

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
	optionalAuthRouter := router.Methods("GET", "POST", "PATCH", "DELETE").Subrouter()
	optionalAuthRouter.HandleFunc("", h.getStores).Methods("GET")
	optionalAuthRouter.HandleFunc("/pages", h.getStoresPages).Methods("GET")
	optionalAuthRouter.HandleFunc("/{storeId}", h.getStore).Methods("GET")
	optionalAuthRouter.HandleFunc("/settings/{storeId}", h.getStoreSettings).Methods("GET")
	optionalAuthRouter.HandleFunc("/address/{storeId}", h.getStoreAddresses).Methods("GET")
	optionalAuthRouter.HandleFunc("/phonenumber/{storeId}", h.getStorePhoneNumbers).Methods("GET")
	optionalAuthRouter.Use(h.authHandler.WithJWTAuthOptional(h.db))

	registerRouter := router.PathPrefix("/register").Methods("POST").Subrouter()
	registerRouter.HandleFunc("",
		h.authHandler.WithActionPermissionAuth(
			h.register,
			h.db,
			[]types.Action{types.ActionCanAddStore},
		),
	).Methods("POST")
	registerRouter.Use(h.authHandler.WithJWTAuth(h.db))
	registerRouter.Use(h.authHandler.WithCSRFToken())
	registerRouter.Use(h.authHandler.WithVerifiedEmail(h.db))
	registerRouter.Use(h.authHandler.WithUnbannedProfile(h.db))

	withAuthRouter := router.Methods("GET", "POST", "PATCH", "DELETE").Subrouter()
	withAuthRouter.HandleFunc("/me", h.getMyStores).Methods("GET")
	withAuthRouter.HandleFunc("/me/{storeId}", h.getMyStore).Methods("GET")
	withAuthRouter.HandleFunc("/{storeId}", h.authHandler.WithActionPermissionAuth(
		h.updateStore,
		h.db,
		[]types.Action{types.ActionCanUpdateStore},
	)).Methods("PATCH")
	withAuthRouter.HandleFunc("/{storeId}", h.authHandler.WithActionPermissionAuth(
		h.deleteStore,
		h.db,
		[]types.Action{types.ActionCanDeleteStore},
	)).Methods("DELETE")
	withAuthRouter.Use(h.authHandler.WithJWTAuth(h.db))
	withAuthRouter.Use(h.authHandler.WithCSRFToken())
	withAuthRouter.Use(h.authHandler.WithVerifiedEmail(h.db))
	withAuthRouter.Use(h.authHandler.WithUnbannedProfile(h.db))

	settingsRouter := withAuthRouter.PathPrefix("/settings").Subrouter()
	settingsRouter.HandleFunc("/me/{storeId}", h.getMyStoreSettings).Methods("GET")
	settingsRouter.HandleFunc("/{storeId}", h.authHandler.WithActionPermissionAuth(
		h.updateSettings,
		h.db,
		[]types.Action{types.ActionCanUpdateStore},
	)).Methods("PATCH")

	addressHandlerRouter := withAuthRouter.PathPrefix("/address").Subrouter()
	addressHandlerRouter.HandleFunc("/{storeId}", h.authHandler.WithActionPermissionAuth(
		h.createAddress,
		h.db,
		[]types.Action{types.ActionCanUpdateStore},
	)).Methods("POST")
	addressHandlerRouter.HandleFunc("/me/{storeId}", h.getMyStoreAddresses).Methods("GET")
	addressHandlerRouter.HandleFunc("/{addrId}", h.authHandler.WithActionPermissionAuth(
		h.updateAddress,
		h.db,
		[]types.Action{types.ActionCanUpdateStore},
	)).Methods("PATCH")
	addressHandlerRouter.HandleFunc("/{addrId}", h.authHandler.WithActionPermissionAuth(
		h.deleteAddress,
		h.db,
		[]types.Action{types.ActionCanUpdateStore},
	)).Methods("DELETE")

	phoneNumberHandlerRouter := withAuthRouter.PathPrefix("/phonenumber").Subrouter()
	phoneNumberHandlerRouter.HandleFunc("/{storeId}", h.authHandler.WithActionPermissionAuth(
		h.createPhoneNumber,
		h.db,
		[]types.Action{types.ActionCanUpdateStore},
	)).Methods("POST")
	phoneNumberHandlerRouter.HandleFunc("/me/{storeId}", h.getMyStorePhoneNumbers).Methods("GET")
	phoneNumberHandlerRouter.HandleFunc("/{phoneId}", h.authHandler.WithActionPermissionAuth(
		h.updatePhoneNumber,
		h.db,
		[]types.Action{types.ActionCanUpdateStore},
	)).Methods("PATCH")
	phoneNumberHandlerRouter.HandleFunc("/{phoneId}", h.authHandler.WithActionPermissionAuth(
		h.deletePhoneNumber,
		h.db,
		[]types.Action{types.ActionCanUpdateStore},
	)).Methods("DELETE")
}

// register godoc
// @Summary      Register a new store
// @Description  Creates a new store with the authenticated user as owner. Requires permission to add stores.
// @Tags         store
// @Accept       json
// @Produce      json
// @Param        store  body      types.CreateStorePayload  true  "Store details"
// @Success      201    {object}  types.NewStoreResponse
// @Failure      400    {object}  types.HTTPError
// @Failure      401    {object}  types.HTTPError
// @Failure      403    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/register [post]
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var store types.CreateStorePayload
	err := utils.ParseRequestPayload(r, &store)
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

	user, err := h.db.GetUserById(userId.(int))
	if err != nil {
		if err == types.ErrUserNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	createdStore, err := h.db.CreateStore(types.CreateStorePayload{
		Name:        store.Name,
		Description: store.Description,
		OwnerId:     user.Id,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	res := types.NewStoreResponse{
		StoreId: createdStore,
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, res, nil)
}

// createAddress godoc
// @Summary      Create a new store address
// @Description  Creates a new address for the specified store. Requires permission to update stores.
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        storeId  path      int                             true  "Store ID"
// @Param        address  body      types.CreateStoreAddressPayload  true  "Address details"
// @Success      201      {object}  types.NewAddressResponse
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/address/{storeId} [post]
func (h *Handler) createAddress(w http.ResponseWriter, r *http.Request) {
	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
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

	store, err := h.db.GetStoreById(storeId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId.(int) {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	var payload types.CreateStoreAddressPayload
	err = utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	addrId, err := h.db.CreateStoreAddress(types.CreateStoreAddressPayload{
		City:    payload.City,
		State:   payload.State,
		Street:  payload.Street,
		Zipcode: payload.Zipcode,
		Details: payload.Details,
		StoreId: store.Id,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	res := types.NewAddressResponse{
		AddressId: addrId,
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, res, nil)
}

// createPhoneNumber godoc
// @Summary      Create a new store phone number
// @Description  Creates a new phone number for the specified store. Requires permission to update stores.
// @Tags         phone
// @Accept       json
// @Produce      json
// @Param        storeId  path      int                               true  "Store ID"
// @Param        phone    body      types.CreateStorePhoneNumberPayload  true  "Phone number details"
// @Success      201      {object}  types.NewPhoneNumberResponse
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/phonenumber/{storeId} [post]
func (h *Handler) createPhoneNumber(w http.ResponseWriter, r *http.Request) {
	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
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

	store, err := h.db.GetStoreById(storeId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId.(int) {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	var payload types.CreateStorePhoneNumberPayload
	err = utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	phoneId, err := h.db.CreateStorePhoneNumber(types.CreateStorePhoneNumberPayload{
		CountryCode: payload.CountryCode,
		Number:      payload.Number,
		StoreId:     store.Id,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	res := types.NewPhoneNumberResponse{
		PhoneNumberId: phoneId,
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, res, nil)
}

// getMyStoreAddresses godoc
// @Summary      Get current user's store addresses
// @Description  Returns all addresses belonging to the specified store owned by current user with optional filters
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        storeId     path      int     true   "Store ID"
// @Param        visibility  query     string  false  "Filter by visibility status (public/private)"
// @Success      200         {array}   types.StoreAddress
// @Failure      400         {object}  types.HTTPError
// @Failure      401         {object}  types.HTTPError
// @Failure      403         {object}  types.HTTPError
// @Failure      404         {object}  types.HTTPError
// @Failure      500         {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/address/me/{storeId} [get]
func (h *Handler) getMyStoreAddresses(w http.ResponseWriter, r *http.Request) {
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

	query := types.StoreAddressSearchQuery{}

	queryMapping := map[string]any{
		"visibility": &query.VisibilityStatus,
	}

	queryValues := r.URL.Query()

	err = utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if query.VisibilityStatus != nil {
		if !query.VisibilityStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidVisibilityStatusOption,
			)
			return
		}
	}

	addresses, err := h.db.GetStoreAddresses(store.Id, query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, addresses, nil)
}

// getMyStorePhoneNumbers godoc
// @Summary      Get current user's store phone numbers
// @Description  Returns all phone numbers belonging to the specified store owned by current user with optional filters
// @Tags         phone
// @Accept       json
// @Produce      json
// @Param        storeId     path      int     true   "Store ID"
// @Param        visibility  query     string  false  "Filter by visibility status (public/private)"
// @Param        verified    query     string  false  "Filter by verification status (verified/unverified)"
// @Success      200         {array}   types.StorePhoneNumber
// @Failure      400         {object}  types.HTTPError
// @Failure      401         {object}  types.HTTPError
// @Failure      403         {object}  types.HTTPError
// @Failure      404         {object}  types.HTTPError
// @Failure      500         {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/phonenumber/me/{storeId} [get]
func (h *Handler) getMyStorePhoneNumbers(w http.ResponseWriter, r *http.Request) {
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

	query := types.StorePhoneNumberSearchQuery{}

	queryMapping := map[string]any{
		"visibility": &query.VisibilityStatus,
		"verified":   &query.VerificationStatus,
	}

	queryValues := r.URL.Query()

	err = utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if query.VisibilityStatus != nil {
		if !query.VisibilityStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidVisibilityStatusOption,
			)
			return
		}
	}

	if query.VerificationStatus != nil {
		if !query.VerificationStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidVerificationStatusOption,
			)
			return
		}
	}

	phones, err := h.db.GetStorePhoneNumbers(store.Id, query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, phones, nil)
}

// updateAddress godoc
// @Summary      Update a store address
// @Description  Updates an existing address belonging to the specified store. Requires permission to update stores.
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        addrId  path      int                             true  "Address ID"
// @Param        address body      types.UpdateStoreAddressPayload  true  "Address update payload"
// @Success      200			"Address updated"
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/address/{addrId} [patch]
func (h *Handler) updateAddress(w http.ResponseWriter, r *http.Request) {
	addrId, err := utils.ParseIntURLParam("addrId", mux.Vars(r))
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

	var payload types.UpdateStoreAddressPayload
	err = utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	storeAddr, err := h.db.GetStoreAddressById(addrId)
	if err != nil {
		if err == types.ErrStoreAddressNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	store, err := h.db.GetStoreById(storeAddr.StoreId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessAddress)
		return
	}

	err = h.db.UpdateStoreAddress(addrId, store.Id, types.UpdateStoreAddressPayload{
		State:    payload.State,
		City:     payload.City,
		Street:   payload.Street,
		Zipcode:  payload.Zipcode,
		Details:  payload.Details,
		IsPublic: payload.IsPublic,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// updatePhoneNumber godoc
// @Summary      Update a store phone number
// @Description  Updates an existing phone number belonging to the specified store. Requires permission to update stores.
// @Tags         phone
// @Accept       json
// @Produce      json
// @Param        phoneId  path      int                               true  "Phone number ID"
// @Param        phone    body      types.UpdateStorePhoneNumberPayload  true  "Phone number update payload"
// @Success      200			"Phone number updated"
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/phonenumber/{phoneId} [patch]
func (h *Handler) updatePhoneNumber(w http.ResponseWriter, r *http.Request) {
	phoneId, err := utils.ParseIntURLParam("phoneId", mux.Vars(r))
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

	var payload types.UpdateStorePhoneNumberPayload
	err = utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	storePhone, err := h.db.GetStorePhoneNumberById(phoneId)
	if err != nil {
		if err == types.ErrStorePhoneNumberNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	store, err := h.db.GetStoreById(storePhone.StoreId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessPhoneNumber)
		return
	}

	err = h.db.UpdateStorePhoneNumber(phoneId, store.Id, types.UpdateStorePhoneNumberPayload{
		CountryCode: payload.CountryCode,
		Number:      payload.Number,
		IsPublic:    payload.IsPublic,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// deleteAddress godoc
// @Summary      Delete a store address
// @Description  Deletes an existing address belonging to the specified store. Requires permission to update stores.
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        addrId  path      int  true  "Address ID"
// @Success      200			"Address deleted"
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/address/{addrId} [delete]
func (h *Handler) deleteAddress(w http.ResponseWriter, r *http.Request) {
	addrId, err := utils.ParseIntURLParam("addrId", mux.Vars(r))
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

	storeAddr, err := h.db.GetStoreAddressById(addrId)
	if err != nil {
		if err == types.ErrStoreAddressNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	store, err := h.db.GetStoreById(storeAddr.StoreId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessAddress)
		return
	}

	err = h.db.DeleteStoreAddress(addrId, store.Id)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// deletePhoneNumber godoc
// @Summary      Delete a store phone number
// @Description  Deletes an existing phone number belonging to the specified store. Requires permission to update stores.
// @Tags         phone
// @Accept       json
// @Produce      json
// @Param        phoneId  path      int  true  "Phone number ID"
// @Success      200			"Phone number deleted"
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/phonenumber/{phoneId} [delete]
func (h *Handler) deletePhoneNumber(w http.ResponseWriter, r *http.Request) {
	phoneId, err := utils.ParseIntURLParam("phoneId", mux.Vars(r))
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

	storePhone, err := h.db.GetStorePhoneNumberById(phoneId)
	if err != nil {
		if err == types.ErrStorePhoneNumberNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	store, err := h.db.GetStoreById(storePhone.StoreId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessPhoneNumber)
		return
	}

	err = h.db.DeleteStorePhoneNumber(phoneId, store.Id)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// getStoreAddresses godoc
// @Summary      Get store addresses
// @Description  Returns all addresses belonging to the specified store with optional filters. Response fields are filtered based on store privacy settings and requester's permissions.
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        storeId     path      int     true   "Store ID"
// @Param        visibility  query     string  false  "Filter by visibility status (public/private)"
// @Success      200         {array}   types.StoreAddress
// @Failure      400         {object}  types.HTTPError
// @Failure      404         {object}  types.HTTPError
// @Failure      500         {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/address/{storeId} [get]
func (h *Handler) getStoreAddresses(w http.ResponseWriter, r *http.Request) {
	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	cLoggedUserRoleId := ctx.Value("userRoleId")
	userHasFullAccess := false

	if cLoggedUserRoleId != nil {
		loggedUserRoleId := cLoggedUserRoleId.(int)

		loggedUserHasFullAccess, err := h.db.IsRoleHasAllResourcePermissions(
			[]types.Resource{types.ResourceStoresFullAccess},
			loggedUserRoleId,
		)
		if err != nil {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
			return
		}

		userHasFullAccess = loggedUserHasFullAccess
	}

	query := types.StoreAddressSearchQuery{}

	queryMapping := map[string]any{
		"visibility": &query.VisibilityStatus,
	}

	queryValues := r.URL.Query()

	err = utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if query.VisibilityStatus != nil {
		if !query.VisibilityStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidVisibilityStatusOption,
			)
			return
		}
	}

	if !userHasFullAccess {
		query.VisibilityStatus = utils.Ptr(types.SettingVisibilityStatusPublic)
	}

	addresses, err := h.db.GetStoreAddresses(storeId, query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	filteredAddresses := []map[string]any{}

	for _, a := range addresses {
		f := utils.FilterStruct(a, map[string]bool{
			"public":         true,
			"isPublic":       a.IsPublic,
			"needPermission": userHasFullAccess,
		})

		filteredAddresses = append(filteredAddresses, f)
	}

	utils.WriteJSONInResponse(w, http.StatusOK, filteredAddresses, nil)
}

// getStorePhoneNumbers godoc
// @Summary      Get store phone numbers
// @Description  Returns all phone numbers belonging to the specified store with optional filters. Response fields are filtered based on store privacy settings and requester's permissions.
// @Tags         phone
// @Accept       json
// @Produce      json
// @Param        storeId     path      int     true   "Store ID"
// @Param        visibility  query     string  false  "Filter by visibility status (public/private)"
// @Param        verified    query     string  false  "Filter by verification status (verified/unverified)"
// @Success      200         {array}   types.StorePhoneNumber
// @Failure      400         {object}  types.HTTPError
// @Failure      404         {object}  types.HTTPError
// @Failure      500         {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/phonenumber/{storeId} [get]
func (h *Handler) getStorePhoneNumbers(w http.ResponseWriter, r *http.Request) {
	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	cLoggedUserRoleId := ctx.Value("userRoleId")
	userHasFullAccess := false

	if cLoggedUserRoleId != nil {
		loggedUserRoleId := cLoggedUserRoleId.(int)

		loggedUserHasFullAccess, err := h.db.IsRoleHasAllResourcePermissions(
			[]types.Resource{types.ResourceStoresFullAccess},
			loggedUserRoleId,
		)
		if err != nil {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
			return
		}

		userHasFullAccess = loggedUserHasFullAccess
	}

	query := types.StorePhoneNumberSearchQuery{}

	queryMapping := map[string]any{
		"visibility": &query.VisibilityStatus,
		"verified":   &query.VerificationStatus,
	}

	queryValues := r.URL.Query()

	err = utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if query.VisibilityStatus != nil {
		if !query.VisibilityStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidVisibilityStatusOption,
			)
			return
		}
	}

	if query.VerificationStatus != nil {
		if !query.VerificationStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidVerificationStatusOption,
			)
			return
		}
	}

	if !userHasFullAccess {
		query.VisibilityStatus = utils.Ptr(types.SettingVisibilityStatusPublic)
		query.VerificationStatus = utils.Ptr(types.CredentialVerificationStatusVerified)
	}

	phones, err := h.db.GetStorePhoneNumbers(storeId, query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	filteredPhones := []map[string]any{}

	for _, p := range phones {
		f := utils.FilterStruct(p, map[string]bool{
			"public":         true,
			"isPublic":       p.IsPublic,
			"needPermission": userHasFullAccess,
		})

		filteredPhones = append(filteredPhones, f)
	}

	utils.WriteJSONInResponse(w, http.StatusOK, filteredPhones, nil)
}

// getStores godoc
// @Summary      Get stores list
// @Description  Retrieves a paginated list of stores with optional filtering. The response fields are filtered based on store privacy settings and requester's permissions.
// @Tags         store
// @Produce      json
// @Param        name   query     string  false  "Filter stores by name (partial match)"
// @Param        owner  query     int     false  "Filter stores by owner ID (requires permissions)"
// @Param        p      query     int     false  "Page number (default: 1)"
// @Success      200    {array}   types.Store     "List of stores with filtered fields"
// @Failure      400    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store [get]
func (h *Handler) getStores(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cLoggedUserRoleId := ctx.Value("userRoleId")
	userHasFullAccess := false

	if cLoggedUserRoleId != nil {
		loggedUserRoleId := cLoggedUserRoleId.(int)

		loggedUserHasFullAccess, err := h.db.IsRoleHasAllResourcePermissions(
			[]types.Resource{types.ResourceStoresFullAccess},
			loggedUserRoleId,
		)
		if err != nil {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
			return
		}

		userHasFullAccess = loggedUserHasFullAccess
	}

	query := types.StoreSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"name":  &query.Name,
		"owner": &query.OwnerId,
		"p":     &page,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query.Limit = utils.Ptr(int(config.Env.MaxStoresInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	if !userHasFullAccess {
		query.OwnerId = nil
	}

	stores, err := h.db.GetStoresWithSettings(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	filteredStores := []map[string]any{}

	for _, s := range stores {
		f := utils.FilterStruct(s.Store, map[string]bool{
			"public":         true,
			"publicOwner":    s.PublicOwner,
			"needPermission": userHasFullAccess,
		})

		filteredStores = append(filteredStores, f)
	}

	utils.WriteJSONInResponse(w, http.StatusOK, filteredStores, nil)
}

// getStoresPages godoc
// @Summary      Get total stores pages count
// @Description  Calculates the total number of pages available for stores listing based on filters and pagination settings
// @Tags         store
// @Produce      json
// @Param        name   query     string  false  "Filter stores by name (partial match)"
// @Param        owner  query     int     false  "Filter stores by owner ID (requires permissions)"
// @Success      200    {object}  types.TotalPageCountResponse  "Returns total page count"
// @Failure      400    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/pages [get]
func (h *Handler) getStoresPages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cLoggedUserRoleId := ctx.Value("userRoleId")
	userHasFullAccess := false

	if cLoggedUserRoleId != nil {
		loggedUserRoleId := cLoggedUserRoleId.(int)

		loggedUserHasFullAccess, err := h.db.IsRoleHasAllResourcePermissions(
			[]types.Resource{types.ResourceStoresFullAccess},
			loggedUserRoleId,
		)
		if err != nil {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
			return
		}

		userHasFullAccess = loggedUserHasFullAccess
	}

	query := types.StoreSearchQuery{}

	queryMapping := map[string]any{
		"name": &query.Name,
		"role": &query.OwnerId,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if !userHasFullAccess {
		query.OwnerId = nil
	}

	count, err := h.db.GetStoresCount(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxStoresInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, types.TotalPageCountResponse{
		Pages: pageCount,
	}, nil)
}

// getStore godoc
// @Summary      Get store details
// @Description  Returns details for a specific store. The response fields are filtered based on store privacy settings and requester's permissions.
// @Tags         store
// @Produce      json
// @Param        storeId  path      int     true  "Store ID"
// @Success      200      {object}  types.Store
// @Failure      400      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/{storeId} [get]
func (h *Handler) getStore(w http.ResponseWriter, r *http.Request) {
	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	cLoggedUserRoleId := ctx.Value("userRoleId")
	userHasFullAccess := false

	if cLoggedUserRoleId != nil {
		loggedUserRoleId := cLoggedUserRoleId.(int)

		loggedUserHasFullAccess, err := h.db.IsRoleHasAllResourcePermissions(
			[]types.Resource{types.ResourceStoresFullAccess},
			loggedUserRoleId,
		)
		if err != nil {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
			return
		}

		userHasFullAccess = loggedUserHasFullAccess
	}

	store, err := h.db.GetStoreWithSettingsById(storeId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	filteredStore := utils.FilterStruct(store.Store, map[string]bool{
		"public":         true,
		"publicOwner":    store.PublicOwner,
		"needPermission": userHasFullAccess,
	})

	utils.WriteJSONInResponse(w, http.StatusOK, filteredStore, nil)
}

// getMyStores godoc
// @Summary      Get current user's stores
// @Description  Returns all stores owned by the current user with optional name filtering
// @Tags         store
// @Produce      json
// @Param        name  query     string  false  "Filter stores by name (partial match)"
// @Success      200   {array}   types.Store
// @Failure      400   {object}  types.HTTPError
// @Failure      401   {object}  types.HTTPError
// @Failure      500   {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/me [get]
func (h *Handler) getMyStores(w http.ResponseWriter, r *http.Request) {
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

	query := types.StoreSearchQuery{}

	queryMapping := map[string]any{
		"name": &query.Name,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query.OwnerId = utils.Ptr(userId.(int))

	stores, err := h.db.GetStores(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	filteredStores := []map[string]any{}

	for _, s := range stores {
		f := utils.FilterStruct(s, map[string]bool{
			"public":         true,
			"publicOwner":    true,
			"needPermission": true,
		})

		filteredStores = append(filteredStores, f)
	}

	utils.WriteJSONInResponse(w, http.StatusOK, filteredStores, nil)
}

// getMyStore godoc
// @Summary      Get current user's store details
// @Description  Returns details for a specific store owned by the current user
// @Tags         store
// @Produce      json
// @Param        storeId  path      int     true  "Store ID"
// @Success      200      {object}  types.Store
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/me/{storeId} [get]
func (h *Handler) getMyStore(w http.ResponseWriter, r *http.Request) {
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

	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	store, err := h.db.GetStoreById(storeId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId.(int) {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	filteredStore := utils.FilterStruct(store, map[string]bool{
		"public":         true,
		"publicOwner":    true,
		"needPermission": true,
	})

	utils.WriteJSONInResponse(w, http.StatusOK, filteredStore, nil)
}

// updateStore godoc
// @Summary      Update store details
// @Description  Updates details for a specific store. Requires permission to update stores.
// @Tags         store
// @Accept       json
// @Produce      json
// @Param        storeId  path      int                       true  "Store ID"
// @Param        store    body      types.UpdateStorePayload  true  "Store update payload"
// @Success      200      "Store updated"
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/{storeId} [patch]
func (h *Handler) updateStore(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateStorePayload
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

	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

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

	if payload.Name != nil {
		existingStore, err := h.db.GetStoreByName(*payload.Name)
		if err == nil && existingStore.Id != -1 {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrDuplicateStoreName,
			)
			return
		}
		if err != types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	err = h.db.UpdateStore(storeId, types.UpdateStorePayload{
		Name:        payload.Name,
		Description: payload.Description,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// updateSettings godoc
// @Summary      Update store settings
// @Description  Updates privacy settings for a specific store. Requires permission to update stores.
// @Tags         settings
// @Accept       json
// @Produce      json
// @Param        storeId  path      int                               true  "Store ID"
// @Param        settings body      types.UpdateStoreSettingsPayload  true  "Settings update payload"
// @Success      200      "Store settings updated"
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/settings/{storeId} [patch]
func (h *Handler) updateSettings(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateStoreSettingsPayload
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

	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

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

	err = h.db.UpdateStoreSettings(store.Id, payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// getStoreSettings godoc
// @Summary      Get store settings
// @Description  Returns the settings for a specific store. Response fields are filtered based on privacy settings.
// @Tags         settings
// @Produce      json
// @Param        storeId  path      int  true  "Store ID"
// @Success      200      {object}  types.StoreSettings
// @Failure      400      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/settings/{storeId} [get]
func (h *Handler) getStoreSettings(w http.ResponseWriter, r *http.Request) {
	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	settings, err := h.db.GetStoreSettings(storeId)
	if err != nil {
		if err == types.ErrStoreSettingsNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	filteredSettings := utils.FilterStruct(settings, map[string]bool{
		"public": true,
	})

	utils.WriteJSONInResponse(w, http.StatusOK, filteredSettings, nil)
}

// getMyStoreSettings godoc
// @Summary      Get current user's store settings
// @Description  Returns all settings for a store owned by the current user (includes private settings)
// @Tags         settings
// @Produce      json
// @Param        storeId  path      int  true  "Store ID"
// @Success      200      {object}  types.StoreSettings
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/settings/me/{storeId} [get]
func (h *Handler) getMyStoreSettings(w http.ResponseWriter, r *http.Request) {
	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
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

	store, err := h.db.GetStoreById(storeId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId.(int) {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	settings, err := h.db.GetStoreSettings(store.Id)
	if err != nil {
		if err == types.ErrStoreSettingsNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	settingsRes := utils.FilterStruct(settings, map[string]bool{
		"public":         true,
		"private":        true,
		"needPermission": true,
	})

	utils.WriteJSONInResponse(w, http.StatusOK, settingsRes, nil)
}

// deleteStore godoc
// @Summary      Delete a store
// @Description  Permanently deletes a store and all its associated data. Requires permission to delete stores.
// @Tags         store
// @Produce      json
// @Param        storeId  path      int  true  "Store ID"
// @Success      200      "Store deleted"
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /store/{storeId} [delete]
func (h *Handler) deleteStore(w http.ResponseWriter, r *http.Request) {
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

	storeId, err := utils.ParseIntURLParam("storeId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

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

	err = h.db.DeleteStore(storeId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}
