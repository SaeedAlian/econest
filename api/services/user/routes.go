package user

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
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

func NewHandler(db *db_manager.Manager, authHandler *auth.AuthHandler) *Handler {
	return &Handler{db: db, authHandler: authHandler}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	registerRouter := router.PathPrefix("/register").Methods("POST").Subrouter()
	registerRouter.HandleFunc("/vendor", h.register("Vendor"))
	registerRouter.HandleFunc("/customer", h.register("Customer"))

	authRouter := router.Methods("POST").Subrouter()
	authRouter.HandleFunc("/login", h.login)
	authRouter.HandleFunc("/refresh", h.refresh).Methods("POST")

	logoutRouter := router.Methods("POST").Subrouter()
	logoutRouter.HandleFunc("/logout", h.logout)
	logoutRouter.Use(h.authHandler.WithJWTAuth(h.db))
	logoutRouter.Use(h.authHandler.WithCSRFToken())

	withAuthRouter := router.Methods("GET", "POST", "PATCH", "DELETE").Subrouter()
	withAuthRouter.HandleFunc("/me", h.getMe)
	withAuthRouter.Use(h.authHandler.WithJWTAuth(h.db))
	withAuthRouter.Use(h.authHandler.WithCSRFToken())

	addressHandlerRouter := withAuthRouter.PathPrefix("/address").Subrouter()
	addressHandlerRouter.HandleFunc("/", h.createAddress).Methods("POST")
	addressHandlerRouter.HandleFunc("/", h.getMyAddresses).Methods("GET")
	addressHandlerRouter.HandleFunc("/{addrId}", h.updateAddress).Methods("PATCH")
	addressHandlerRouter.HandleFunc("/{addrId}", h.deleteAddress).Methods("DELETE")

	phoneNumberHandlerRouter := withAuthRouter.PathPrefix("/phonenumber").Subrouter()
	phoneNumberHandlerRouter.HandleFunc("/", h.createPhoneNumber).Methods("POST")
	phoneNumberHandlerRouter.HandleFunc("/", h.getMyPhoneNumbers).Methods("GET")
	phoneNumberHandlerRouter.HandleFunc("/{phoneId}", h.updatePhoneNumber).Methods("PATCH")
	phoneNumberHandlerRouter.HandleFunc("/{phoneId}", h.deletePhoneNumber).Methods("DELETE")
}

func (h *Handler) register(roleName string) func(w http.ResponseWriter, r *http.Request) {
	var callback func(w http.ResponseWriter, r *http.Request)
	callback = func(w http.ResponseWriter, r *http.Request) {
		var user types.CreateUserPayload
		if err := utils.ParseJSONFromRequest(r, &user); err != nil {
			utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidUserPayload)
			return
		}

		if err := utils.Validator.Struct(user); err != nil {
			errors := err.(validator.ValidationErrors)
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidPayload(errors[0]),
			)
			return
		}

		if u, _ := h.db.GetUserByUsernameOrEmail(user.Username, user.Email); u != nil {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrDuplicateUsernameOrEmail,
			)
			return
		}

		hashedPassword, err := auth.HashPassword(user.Password)
		if err != nil {
			utils.WriteErrorInResponse(
				w,
				http.StatusInternalServerError,
				types.ErrInternalServer,
			)
			return
		}

		customerRole, err := h.db.GetRoleByName(roleName)
		if err != nil {
			utils.WriteErrorInResponse(
				w,
				http.StatusInternalServerError,
				err,
			)
			return
		}

		created_user, err := h.db.CreateUser(types.CreateUserPayload{
			Username:  strings.ToLower(user.Username),
			Email:     strings.ToLower(user.Email),
			FullName:  user.FullName,
			BirthDate: user.BirthDate,
			Password:  hashedPassword,
			RoleId:    customerRole.Id,
		})
		if err != nil {
			utils.WriteErrorInResponse(
				w,
				http.StatusInternalServerError,
				types.ErrInternalServer,
			)
			return
		}

		utils.WriteJSONInResponse(w, http.StatusCreated, created_user, nil)
	}

	return callback
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJSONFromRequest(r, &payload); err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidLoginPayload)
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidPayload(errors[0]))
		return
	}

	user, err := h.db.GetUserByUsername(payload.Username)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidCredentials)
		return
	}

	if isPasswordCorrect := auth.ComparePassword(payload.Password, user.Password); !isPasswordCorrect {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidCredentials)
		return
	}

	accessToken, _, err := h.authHandler.GenerateToken(
		user.Id,
		config.Env.AccessTokenExpirationInMin,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	refreshToken, refreshTokenJTI, err := h.authHandler.GenerateToken(
		user.Id,
		config.Env.RefreshTokenExpirationInMin,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	err = h.authHandler.SaveRefreshToken(
		refreshTokenJTI,
		user.Id,
		config.Env.RefreshTokenExpirationInMin,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires: time.Now().
			Add(time.Duration(config.Env.RefreshTokenExpirationInMin) * time.Minute),
	}

	http.SetCookie(w, &refreshTokenCookie)

	csrf, err := h.authHandler.GenerateCSRFToken()
	if err != nil {
		h.authHandler.RevokeRefreshToken(refreshTokenJTI)
		utils.DeleteCookie(w, &refreshTokenCookie)
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	err = h.authHandler.SaveCSRFToken(user.Id, csrf, config.Env.CSRFTokenExpirationInMin)
	if err != nil {
		h.authHandler.RevokeRefreshToken(refreshTokenJTI)
		utils.DeleteCookie(w, &refreshTokenCookie)
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	csrfTokenCookie := http.Cookie{
		Name:     "csrf_token",
		Value:    csrf,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: false,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(config.Env.CSRFTokenExpirationInMin) * time.Minute),
	}

	http.SetCookie(w, &csrfTokenCookie)

	utils.WriteJSONInResponse(w, http.StatusOK, types.LoginResponsePayload{
		AccessToken: accessToken,
	}, nil)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusUnauthorized, types.ErrRefreshTokenNotFound)
		return
	}
	claims := types.UserJWTClaims{}
	tokenRes, err := h.authHandler.ValidateToken(refreshToken.Value, &claims)
	if err != nil || !tokenRes.Valid {
		utils.WriteErrorInResponse(w, http.StatusUnauthorized, types.ErrInvalidRefreshToken)
		return
	}

	isValid, err := h.authHandler.IsRefreshTokenValid(claims.JTI)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}
	if !isValid {
		utils.WriteErrorInResponse(w, http.StatusUnauthorized, types.ErrInvalidRefreshToken)
		return
	}

	user, err := h.db.GetUserById(claims.UserId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusUnauthorized, types.ErrInvalidRefreshToken)
		return
	}

	newAccessToken, _, err := h.authHandler.GenerateToken(
		user.Id,
		config.Env.AccessTokenExpirationInMin,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	newRefreshToken, newRefreshTokenJTI, err := h.authHandler.GenerateToken(
		user.Id,
		config.Env.RefreshTokenExpirationInMin,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	err = h.authHandler.RotateRefreshToken(
		claims.JTI,
		newRefreshTokenJTI,
		user.Id,
		config.Env.RefreshTokenExpirationInMin,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires: time.Now().
			Add(time.Duration(config.Env.RefreshTokenExpirationInMin) * time.Minute),
	}

	http.SetCookie(w, &refreshTokenCookie)

	newCSRFToken, err := h.authHandler.GenerateCSRFToken()
	if err != nil {
		h.authHandler.RevokeRefreshToken(newRefreshTokenJTI)
		utils.DeleteCookie(w, &refreshTokenCookie)
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	err = h.authHandler.RotateCSRFToken(user.Id, newCSRFToken, config.Env.CSRFTokenExpirationInMin)
	if err != nil {
		h.authHandler.RevokeRefreshToken(newRefreshTokenJTI)
		utils.DeleteCookie(w, &refreshTokenCookie)
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	csrfTokenCookie := http.Cookie{
		Name:     "csrf_token",
		Value:    newCSRFToken,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: false,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(config.Env.CSRFTokenExpirationInMin) * time.Minute),
	}

	http.SetCookie(w, &csrfTokenCookie)

	utils.WriteJSONInResponse(w, http.StatusOK, types.LoginResponsePayload{
		AccessToken: newAccessToken,
	}, nil)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")

	if err == nil && refreshToken != nil {
		claims := types.UserJWTClaims{}
		_, err := h.authHandler.ValidateToken(refreshToken.Value, &claims)
		if err == nil {
			h.authHandler.RevokeRefreshToken(claims.JTI)
			h.authHandler.RevokeCSRFToken(claims.UserId)
		}
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires: time.Now().
			Add(time.Duration(config.Env.RefreshTokenExpirationInMin) * time.Minute),
	}

	csrfTokenCookie := http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: false,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(config.Env.CSRFTokenExpirationInMin) * time.Minute),
	}

	utils.DeleteCookie(w, &refreshTokenCookie)
	utils.DeleteCookie(w, &csrfTokenCookie)

	utils.WriteJSONInResponse(w, http.StatusOK, types.LogoutResponsePayload{
		Success: true,
	}, nil)
}

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {
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

	userRes := utils.FilterStruct(user, map[string]bool{
		"public":         true,
		"private":        true,
		"needPermission": true,
	})

	utils.WriteJSONInResponse(w, http.StatusOK, userRes, nil)
}

func (h *Handler) createAddress(w http.ResponseWriter, r *http.Request) {
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

	var payload types.CreateUserAddressPayload
	if err := utils.ParseJSONFromRequest(r, &payload); err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidAddressPayload)
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidPayload(errors[0]))
		return
	}

	addrId, err := h.db.CreateUserAddress(types.CreateUserAddressPayload{
		City:    payload.City,
		State:   payload.State,
		Street:  payload.Street,
		Zipcode: payload.Zipcode,
		Details: payload.Details,
		UserId:  userId.(int),
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrCreateAddress)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, map[string]int{"addressId": addrId}, nil)
}

func (h *Handler) createPhoneNumber(w http.ResponseWriter, r *http.Request) {
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

	var payload types.CreateUserPhoneNumberPayload
	if err := utils.ParseJSONFromRequest(r, &payload); err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidPhoneNumberPayload)
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidPayload(errors[0]))
		return
	}

	phoneId, err := h.db.CreateUserPhoneNumber(types.CreateUserPhoneNumberPayload{
		CountryCode: payload.CountryCode,
		Number:      payload.Number,
		UserId:      userId.(int),
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrCreatePhoneNumber)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, map[string]int{"phoneNumberId": phoneId}, nil)
}

func (h *Handler) getMyAddresses(w http.ResponseWriter, r *http.Request) {
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

	visibilityStatusQuery := r.URL.Query().Get("visibility")
	var visibilityStatus *types.SettingVisibilityStatus = nil

	if visibilityStatusQuery != "" {
		visibilityStatus = utils.Ptr(types.SettingVisibilityStatus(visibilityStatusQuery))
		if !visibilityStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidVisibilityStatusOption,
			)
			return
		}
	}

	query := types.UserAddressSearchQuery{
		VisibilityStatus: visibilityStatus,
	}

	addresses, err := h.db.GetUserAddresses(userId.(int), query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, addresses, nil)
}

func (h *Handler) getMyPhoneNumbers(w http.ResponseWriter, r *http.Request) {
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

	visibilityStatusQuery := r.URL.Query().Get("visibility")
	var visibilityStatus *types.SettingVisibilityStatus = nil

	verificationStatusQuery := r.URL.Query().Get("verified")
	var verificationStatus *types.CredentialVerificationStatus = nil

	if visibilityStatusQuery != "" {
		visibilityStatus = utils.Ptr(types.SettingVisibilityStatus(visibilityStatusQuery))
		if !visibilityStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidVisibilityStatusOption,
			)
			return
		}
	}

	if verificationStatusQuery != "" {
		verificationStatus = utils.Ptr(types.CredentialVerificationStatus(verificationStatusQuery))
		if !verificationStatus.IsValid() {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrInvalidVerificationStatusOption,
			)
			return
		}
	}

	query := types.UserPhoneNumberSearchQuery{
		VerificationStatus: verificationStatus,
		VisibilityStatus:   visibilityStatus,
	}

	phones, err := h.db.GetUserPhoneNumbers(userId.(int), query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrInternalServer)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, phones, nil)
}

func (h *Handler) updateAddress(w http.ResponseWriter, r *http.Request) {
	addrIdParam := mux.Vars(r)["addrId"]

	addrId, err := strconv.Atoi(addrIdParam)
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			types.ErrInvalidAddressId,
		)
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

	var payload types.UpdateUserAddressPayload
	if err := utils.ParseJSONFromRequest(r, &payload); err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidAddressPayload)
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidPayload(errors[0]))
		return
	}

	userAddr, err := h.db.GetUserAddressById(addrId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusNotFound, types.ErrUserAddressNotFound)
		return
	}

	if userAddr.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessAddress)
		return
	}

	err = h.db.UpdateUserAddress(addrId, userId, payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrUpdateAddress)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) updatePhoneNumber(w http.ResponseWriter, r *http.Request) {
	phoneIdParam := mux.Vars(r)["phoneId"]

	phoneId, err := strconv.Atoi(phoneIdParam)
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			types.ErrInvalidPhoneNumberId,
		)
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

	var payload types.UpdateUserPhoneNumberPayload
	if err := utils.ParseJSONFromRequest(r, &payload); err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidPhoneNumberPayload)
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidPayload(errors[0]))
		return
	}

	userPhone, err := h.db.GetUserPhoneNumberById(phoneId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusNotFound, types.ErrUserPhoneNumberNotFound)
		return
	}

	if userPhone.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessPhoneNumber)
		return
	}

	err = h.db.UpdateUserPhoneNumber(phoneId, userId, payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrUpdatePhoneNumber)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) deleteAddress(w http.ResponseWriter, r *http.Request) {
	addrIdParam := mux.Vars(r)["addrId"]

	addrId, err := strconv.Atoi(addrIdParam)
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			types.ErrInvalidAddressId,
		)
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

	userAddr, err := h.db.GetUserAddressById(addrId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusNotFound, types.ErrUserAddressNotFound)
		return
	}

	if userAddr.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessAddress)
		return
	}

	err = h.db.DeleteUserAddress(addrId, userId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrDeleteAddress)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) deletePhoneNumber(w http.ResponseWriter, r *http.Request) {
	phoneIdParam := mux.Vars(r)["phoneId"]

	phoneId, err := strconv.Atoi(phoneIdParam)
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			types.ErrInvalidPhoneNumberId,
		)
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

	userPhone, err := h.db.GetUserPhoneNumberById(phoneId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusNotFound, types.ErrUserPhoneNumberNotFound)
		return
	}

	if userPhone.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessPhoneNumber)
		return
	}

	err = h.db.DeleteUserPhoneNumber(phoneId, userId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrDeletePhoneNumber)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}
