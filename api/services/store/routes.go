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
	optionalAuthRouter.HandleFunc("/", h.getStores).Methods("GET")
	optionalAuthRouter.HandleFunc("/pages", h.getStoresPages).Methods("GET")
	optionalAuthRouter.HandleFunc("/{storeId}", h.getStore).Methods("GET")
	optionalAuthRouter.HandleFunc("/settings/{storeId}", h.getStoreSettings).Methods("GET")
	optionalAuthRouter.HandleFunc("/address/{storeId}", h.getStoreAddresses).Methods("GET")
	optionalAuthRouter.HandleFunc("/phonenumber/{storeId}", h.getStorePhoneNumbers).Methods("GET")
	optionalAuthRouter.Use(h.authHandler.WithJWTAuthOptional(h.db))

	registerRouter := router.PathPrefix("/register").Methods("POST").Subrouter()
	registerRouter.HandleFunc("/",
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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		}

		return
	}

	createdStore, err := h.db.CreateStore(types.CreateStorePayload{
		Name:        store.Name,
		Description: store.Description,
		OwnerId:     user.Id,
	})
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusInternalServerError,
			types.ErrInternalServer,
		)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, map[string]int{"storeId": createdStore}, nil)
}

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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
		utils.WriteErrorInResponse(
			w,
			http.StatusInternalServerError,
			types.ErrInternalServer,
		)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, map[string]int{"addressId": addrId}, nil)
}

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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
		utils.WriteErrorInResponse(
			w,
			http.StatusInternalServerError,
			types.ErrInternalServer,
		)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, map[string]int{"phoneNumberId": phoneId}, nil)
}

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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			err,
		)
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
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, addresses, nil)
}

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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			err,
		)
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
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, phones, nil)
}

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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		}

		return
	}

	store, err := h.db.GetStoreById(storeAddr.StoreId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrUpdateAddress)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		}

		return
	}

	store, err := h.db.GetStoreById(storePhone.StoreId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrUpdatePhoneNumber)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		}

		return
	}

	store, err := h.db.GetStoreById(storeAddr.StoreId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessAddress)
		return
	}

	err = h.db.DeleteStoreAddress(addrId, store.Id)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrDeleteAddress)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		}

		return
	}

	store, err := h.db.GetStoreById(storePhone.StoreId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessPhoneNumber)
		return
	}

	err = h.db.DeleteStorePhoneNumber(phoneId, store.Id)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrDeletePhoneNumber)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			err,
		)
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
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			err,
		)
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
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			err,
		)
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
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	if !userHasFullAccess {
		query.OwnerId = nil
	}

	count, err := h.db.GetStoresCount(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxStoresInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]int32{
		"pages": pageCount,
	}, nil)
}

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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
			return
		}

		userHasFullAccess = loggedUserHasFullAccess
	}

	store, err := h.db.GetStoreWithSettingsById(storeId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	query.OwnerId = utils.Ptr(userId.(int))

	stores, err := h.db.GetStores(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
			utils.WriteErrorInResponse(
				w,
				http.StatusInternalServerError,
				types.ErrInternalServer,
			)
			return
		}
	}

	err = h.db.UpdateStore(storeId, types.UpdateStorePayload{
		Name:        payload.Name,
		Description: payload.Description,
	})
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusInternalServerError,
			types.ErrInternalServer,
		)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	err = h.db.UpdateStoreSettings(store.Id, payload)
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusInternalServerError,
			types.ErrInternalServer,
		)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		}

		return
	}

	filteredSettings := utils.FilterStruct(settings, map[string]bool{
		"public": true,
	})

	utils.WriteJSONInResponse(w, http.StatusOK, filteredSettings, nil)
}

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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	err = h.db.DeleteStore(storeId)
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusInternalServerError,
			types.ErrInternalServer,
		)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}
