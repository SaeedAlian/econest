package user

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/SaeedAlian/econest/api/config"
	db_manager "github.com/SaeedAlian/econest/api/db/manager"
	"github.com/SaeedAlian/econest/api/services/auth"
	"github.com/SaeedAlian/econest/api/services/smtp"
	"github.com/SaeedAlian/econest/api/types"
	"github.com/SaeedAlian/econest/api/utils"
)

type Handler struct {
	db          *db_manager.Manager
	authHandler *auth.AuthHandler
	smtpServer  *smtp.SMTPServer
}

func NewHandler(
	db *db_manager.Manager,
	authHandler *auth.AuthHandler,
	smtpServer *smtp.SMTPServer,
) *Handler {
	return &Handler{db: db, authHandler: authHandler, smtpServer: smtpServer}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	registerRouter := router.PathPrefix("/register").Methods("POST").Subrouter()
	registerRouter.HandleFunc("/vendor", h.registerVendor).Methods("POST")
	registerRouter.HandleFunc("/customer", h.registerCustomer).Methods("POST")

	withRoleRegisterRouter := registerRouter.PathPrefix("/withrole").Subrouter()
	withRoleRegisterRouter.HandleFunc(
		"/{roleName}",
		h.authHandler.WithActionPermissionAuth(
			h.registerWithRole,
			h.db,
			[]types.Action{types.ActionCanAddUserWithRole},
		),
	).Methods("PATCH")

	authRouter := router.Methods("POST").Subrouter()
	authRouter.HandleFunc("/login", h.login).Methods("POST")
	authRouter.HandleFunc("/forgotpass", h.forgotPasswordRequest).Methods("POST")
	authRouter.HandleFunc("/forgotpass", h.resetPassword).Methods("PATCH")
	authRouter.HandleFunc("/refresh", h.refresh).Methods("POST")
	authRouter.HandleFunc("/email/verify", h.verifyEmail).Methods("PATCH")

	logoutRouter := router.Methods("POST").Subrouter()
	logoutRouter.HandleFunc("/logout", h.logout).Methods("POST")
	logoutRouter.Use(h.authHandler.WithJWTAuth(h.db))
	logoutRouter.Use(h.authHandler.WithCSRFToken())

	withAuthRouter := router.Methods("GET", "POST", "PATCH", "DELETE").Subrouter()
	withAuthRouter.HandleFunc("", h.getUsers).Methods("GET")
	withAuthRouter.HandleFunc("/pages", h.getUsersPages).Methods("GET")
	withAuthRouter.HandleFunc("/me", h.getMe).Methods("GET")
	withAuthRouter.HandleFunc("/{userId}", h.getUser).Methods("GET")
	withAuthRouter.HandleFunc("/email/verify", h.verifyEmailRequest).Methods("POST")
	withAuthRouter.HandleFunc("", h.updateProfile).Methods("PATCH")
	withAuthRouter.HandleFunc("/email", h.updateEmail).Methods("PATCH")
	withAuthRouter.HandleFunc("/password", h.updatePassword).Methods("PATCH")
	withAuthRouter.HandleFunc(
		"/ban/{userId}",
		h.authHandler.WithActionPermissionAuth(
			h.banUser,
			h.db,
			[]types.Action{types.ActionCanBanUser},
		),
	).Methods("PATCH")
	withAuthRouter.HandleFunc(
		"/unban/{userId}",
		h.authHandler.WithActionPermissionAuth(
			h.unbanUser,
			h.db,
			[]types.Action{types.ActionCanUnbanUser},
		),
	).Methods("PATCH")
	withAuthRouter.Use(h.authHandler.WithJWTAuth(h.db))
	withAuthRouter.Use(h.authHandler.WithCSRFToken())
	withAuthRouter.Use(h.authHandler.WithUnbannedProfile(h.db))

	settingsRouter := withAuthRouter.PathPrefix("/settings").Subrouter()
	settingsRouter.HandleFunc("/me", h.getMySettings).Methods("GET")
	settingsRouter.HandleFunc("/{userId}", h.getUserSettings).Methods("GET")
	settingsRouter.HandleFunc("", h.updateSettings).Methods("PATCH")

	addressHandlerRouter := withAuthRouter.PathPrefix("/address").Subrouter()
	addressHandlerRouter.HandleFunc("", h.createAddress).Methods("POST")
	addressHandlerRouter.HandleFunc("", h.getMyAddresses).Methods("GET")
	addressHandlerRouter.HandleFunc("/{userId}", h.getUserAddresses).Methods("GET")
	addressHandlerRouter.HandleFunc("/{addrId}", h.updateAddress).Methods("PATCH")
	addressHandlerRouter.HandleFunc("/{addrId}", h.deleteAddress).Methods("DELETE")

	phoneNumberHandlerRouter := withAuthRouter.PathPrefix("/phonenumber").Subrouter()
	phoneNumberHandlerRouter.HandleFunc("", h.createPhoneNumber).Methods("POST")
	phoneNumberHandlerRouter.HandleFunc("", h.getMyPhoneNumbers).Methods("GET")
	phoneNumberHandlerRouter.HandleFunc("/{userId}", h.getUserPhoneNumbers).Methods("GET")
	phoneNumberHandlerRouter.HandleFunc("/{phoneId}", h.updatePhoneNumber).Methods("PATCH")
	phoneNumberHandlerRouter.HandleFunc("/{phoneId}", h.deletePhoneNumber).Methods("DELETE")
}

func (h *Handler) register(roleName string, w *http.ResponseWriter, r **http.Request) {
	var user types.CreateUserPayload
	err := utils.ParseRequestPayload(*r, &user)
	if err != nil {
		utils.WriteErrorInResponse(*w, http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteErrorInResponse(*w, http.StatusInternalServerError, err)
		return
	}

	role, err := h.db.GetRoleByName(roleName)
	if err != nil {
		utils.WriteErrorInResponse(*w, http.StatusBadRequest, err)
		return
	}

	if role.Name == types.DefaultRoleSuperAdmin.String() {
		utils.WriteErrorInResponse(*w, http.StatusForbidden, types.ErrCannotRegisterThisUser)
		return
	}

	createdUser, err := h.db.CreateUser(types.CreateUserPayload{
		Username:  strings.ToLower(user.Username),
		Email:     strings.ToLower(user.Email),
		FullName:  user.FullName,
		BirthDate: user.BirthDate,
		Password:  hashedPassword,
		RoleId:    role.Id,
	})
	if err != nil {
		utils.WriteErrorInResponse(*w, http.StatusBadRequest, err)
		return
	}

	res := types.NewUserResponse{
		UserId: createdUser,
	}

	utils.WriteJSONInResponse(*w, http.StatusCreated, res, nil)
}

func (h *Handler) registerWithRole(w http.ResponseWriter, r *http.Request) {
	roleName, err := utils.ParseStringURLParam("roleName", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	h.register(roleName, &w, &r)
}

func (h *Handler) registerCustomer(w http.ResponseWriter, r *http.Request) {
	h.register(types.DefaultRoleCustomer.String(), &w, &r)
}

func (h *Handler) registerVendor(w http.ResponseWriter, r *http.Request) {
	h.register(types.DefaultRoleVendor.String(), &w, &r)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.db.GetUserByUsername(payload.Username)
	if err != nil {
		if err == types.ErrUserNotFound {
			utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidCredentials)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if user.IsBanned {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrUserIsBanned)
		return
	}

	userRole, err := h.db.GetRoleById(user.RoleId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if userRole.Name == types.DefaultRoleSuperAdmin.String() {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotLoginWithThisUser)
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
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	refreshToken, refreshTokenJTI, err := h.authHandler.GenerateToken(
		user.Id,
		config.Env.RefreshTokenExpirationInMin,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = h.authHandler.SaveRefreshToken(
		refreshTokenJTI,
		user.Id,
		config.Env.RefreshTokenExpirationInMin,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
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
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = h.authHandler.SaveCSRFToken(user.Id, csrf, config.Env.CSRFTokenExpirationInMin)
	if err != nil {
		h.authHandler.RevokeRefreshToken(refreshTokenJTI)
		utils.DeleteCookie(w, &refreshTokenCookie)
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
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
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}
	if !isValid {
		utils.WriteErrorInResponse(w, http.StatusUnauthorized, types.ErrInvalidRefreshToken)
		return
	}

	user, err := h.db.GetUserById(claims.UserId)
	if err != nil {
		if err == types.ErrUserNotFound {
			utils.WriteErrorInResponse(w, http.StatusUnauthorized, types.ErrInvalidRefreshToken)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	newAccessToken, _, err := h.authHandler.GenerateToken(
		user.Id,
		config.Env.AccessTokenExpirationInMin,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	newRefreshToken, newRefreshTokenJTI, err := h.authHandler.GenerateToken(
		user.Id,
		config.Env.RefreshTokenExpirationInMin,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = h.authHandler.RotateRefreshToken(
		claims.JTI,
		newRefreshTokenJTI,
		user.Id,
		config.Env.RefreshTokenExpirationInMin,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
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
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = h.authHandler.RotateCSRFToken(user.Id, newCSRFToken, config.Env.CSRFTokenExpirationInMin)
	if err != nil {
		h.authHandler.RevokeRefreshToken(newRefreshTokenJTI)
		utils.DeleteCookie(w, &refreshTokenCookie)
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
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
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
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
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
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
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	res := types.NewAddressResponse{
		AddressId: addrId,
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, res, nil)
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
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	phoneId, err := h.db.CreateUserPhoneNumber(types.CreateUserPhoneNumberPayload{
		CountryCode: payload.CountryCode,
		Number:      payload.Number,
		UserId:      userId.(int),
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

	query := types.UserAddressSearchQuery{}

	queryMapping := map[string]any{
		"visibility": &query.VisibilityStatus,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
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

	addresses, err := h.db.GetUserAddresses(userId.(int), query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
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

	query := types.UserPhoneNumberSearchQuery{}

	queryMapping := map[string]any{
		"visibility": &query.VisibilityStatus,
		"verified":   &query.VerificationStatus,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
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

	phones, err := h.db.GetUserPhoneNumbers(userId.(int), query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
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

	var payload types.UpdateUserAddressPayload
	err = utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	userAddr, err := h.db.GetUserAddressById(addrId)
	if err != nil {
		if err == types.ErrUserAddressNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if userAddr.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessAddress)
		return
	}

	err = h.db.UpdateUserAddress(addrId, userId, types.UpdateUserAddressPayload{
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

	var payload types.UpdateUserPhoneNumberPayload
	err = utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	userPhone, err := h.db.GetUserPhoneNumberById(phoneId)
	if err != nil {
		if err == types.ErrUserPhoneNumberNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if userPhone.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessPhoneNumber)
		return
	}

	err = h.db.UpdateUserPhoneNumber(phoneId, userId, types.UpdateUserPhoneNumberPayload{
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

	userAddr, err := h.db.GetUserAddressById(addrId)
	if err != nil {
		if err == types.ErrUserAddressNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if userAddr.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessAddress)
		return
	}

	err = h.db.DeleteUserAddress(addrId, userId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
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

	userPhone, err := h.db.GetUserPhoneNumberById(phoneId)
	if err != nil {
		if err == types.ErrUserPhoneNumberNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if userPhone.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessPhoneNumber)
		return
	}

	err = h.db.DeleteUserPhoneNumber(phoneId, userId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) getUserAddresses(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseIntURLParam("userId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	cLoggedUserRoleId := ctx.Value("userRoleId")

	if cLoggedUserRoleId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	query := types.UserAddressSearchQuery{}

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

	loggedUserRoleId := cLoggedUserRoleId.(int)

	loggedUserHasFullAccess, err := h.db.IsRoleHasAllResourcePermissions(
		[]types.Resource{types.ResourceUsersFullAccess},
		loggedUserRoleId,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	if !loggedUserHasFullAccess {
		query.VisibilityStatus = utils.Ptr(types.SettingVisibilityStatusPublic)
	}

	addresses, err := h.db.GetUserAddresses(userId, query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	filteredAddresses := []map[string]any{}

	for _, a := range addresses {
		f := utils.FilterStruct(a, map[string]bool{
			"public":         true,
			"isPublic":       a.IsPublic,
			"needPermission": loggedUserHasFullAccess,
		})

		filteredAddresses = append(filteredAddresses, f)
	}

	utils.WriteJSONInResponse(w, http.StatusOK, filteredAddresses, nil)
}

func (h *Handler) getUserPhoneNumbers(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseIntURLParam("userId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	cLoggedUserRoleId := ctx.Value("userRoleId")

	if cLoggedUserRoleId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	query := types.UserPhoneNumberSearchQuery{}

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

	loggedUserRoleId := cLoggedUserRoleId.(int)

	loggedUserHasFullAccess, err := h.db.IsRoleHasAllResourcePermissions(
		[]types.Resource{types.ResourceUsersFullAccess},
		loggedUserRoleId,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	if !loggedUserHasFullAccess {
		query.VisibilityStatus = utils.Ptr(types.SettingVisibilityStatusPublic)
		query.VerificationStatus = utils.Ptr(types.CredentialVerificationStatusVerified)
	}

	phones, err := h.db.GetUserPhoneNumbers(userId, query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	filteredPhones := []map[string]any{}

	for _, p := range phones {
		f := utils.FilterStruct(p, map[string]bool{
			"public":         true,
			"isPublic":       p.IsPublic,
			"needPermission": loggedUserHasFullAccess,
		})

		filteredPhones = append(filteredPhones, f)
	}

	utils.WriteJSONInResponse(w, http.StatusOK, filteredPhones, nil)
}

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cLoggedUserRoleId := ctx.Value("userRoleId")

	if cLoggedUserRoleId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	query := types.UserSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"name": &query.FullName,
		"role": &query.RoleId,
		"p":    &page,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query.Limit = utils.Ptr(int(config.Env.MaxUsersInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	loggedUserRoleId := cLoggedUserRoleId.(int)

	loggedUserHasFullAccess, err := h.db.IsRoleHasAllResourcePermissions(
		[]types.Resource{types.ResourceUsersFullAccess},
		loggedUserRoleId,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	users, err := h.db.GetUsersWithSettings(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	filteredUsers := []map[string]any{}

	for _, u := range users {
		f := utils.FilterStruct(u.User, map[string]bool{
			"public":          true,
			"publicEmail":     u.PublicEmail,
			"publicBirthDate": u.PublicBirthDate,
			"needPermission":  loggedUserHasFullAccess,
		})

		filteredUsers = append(filteredUsers, f)
	}

	utils.WriteJSONInResponse(w, http.StatusOK, filteredUsers, nil)
}

func (h *Handler) getUsersPages(w http.ResponseWriter, r *http.Request) {
	query := types.UserSearchQuery{}

	queryMapping := map[string]any{
		"name": &query.FullName,
		"role": &query.RoleId,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	count, err := h.db.GetUsersCount(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxUsersInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, types.TotalPageCountResponse{
		Pages: pageCount,
	}, nil)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseIntURLParam("userId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	cLoggedUserRoleId := ctx.Value("userRoleId")

	if cLoggedUserRoleId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	loggedUserRoleId := cLoggedUserRoleId.(int)

	loggedUserHasFullAccess, err := h.db.IsRoleHasAllResourcePermissions(
		[]types.Resource{types.ResourceUsersFullAccess},
		loggedUserRoleId,
	)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	user, err := h.db.GetUserWithSettingsById(userId)
	if err != nil {
		if err == types.ErrUserNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	filteredUser := utils.FilterStruct(user.User, map[string]bool{
		"public":          true,
		"publicEmail":     user.PublicEmail,
		"publicBirthDate": user.PublicBirthDate,
		"needPermission":  loggedUserHasFullAccess,
	})

	utils.WriteJSONInResponse(w, http.StatusOK, filteredUser, nil)
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateUserPayload
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

	if payload.Username != nil {
		existingUser, err := h.db.GetUserByUsername(*payload.Username)
		if err == nil && existingUser.Id != -1 {
			utils.WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrDuplicateUsername,
			)
			return
		}
		if err != types.ErrUserNotFound {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	err = h.db.UpdateUser(userId, types.UpdateUserPayload{
		Username:  payload.Username,
		FullName:  payload.FullName,
		BirthDate: payload.BirthDate,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) updateEmail(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateUserPayload
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

	if payload.Email == nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidPayload)
		return
	}

	existingUser, err := h.db.GetUserByEmail(*payload.Email)
	if err == nil && existingUser.Id != -1 {
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			types.ErrDuplicateUserEmail,
		)
		return
	}
	if err != types.ErrUserNotFound {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = h.db.UpdateUser(userId, types.UpdateUserPayload{
		Email:         payload.Email,
		EmailVerified: utils.Ptr(false),
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) updatePassword(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateUserPasswordPayload
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

	user, err := h.db.GetUserById(userId)
	if err != nil {
		if err == types.ErrUserNotFound {
			utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidCredentials)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if isPasswordCorrect := auth.ComparePassword(*payload.CurrentPassword, user.Password); !isPasswordCorrect {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidCredentials)
		return
	}

	hashedPassword, err := auth.HashPassword(*payload.NewPassword)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = h.db.UpdateUser(userId, types.UpdateUserPayload{
		Password: &hashedPassword,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) updateSettings(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateUserSettingsPayload
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

	err = h.db.UpdateUserSettings(userId, payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) banUser(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseIntURLParam("userId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.db.GetUserById(userId)
	if err != nil {
		if err == types.ErrUserNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	userRole, err := h.db.GetRoleById(user.RoleId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	if userRole.Name == types.DefaultRoleSuperAdmin.String() ||
		userRole.Name == types.DefaultRoleAdmin.String() {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotBanThisUser)
		return
	}

	err = h.db.UpdateUser(userId, types.UpdateUserPayload{
		IsBanned: utils.Ptr(true),
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) unbanUser(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseIntURLParam("userId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.UpdateUser(userId, types.UpdateUserPayload{
		IsBanned: utils.Ptr(false),
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) getUserSettings(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseIntURLParam("userId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	settings, err := h.db.GetUserSettings(userId)
	if err != nil {
		if err == types.ErrUserSettingsNotFound {
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

func (h *Handler) getMySettings(w http.ResponseWriter, r *http.Request) {
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

	settings, err := h.db.GetUserSettings(userId.(int))
	if err != nil {
		if err == types.ErrUserSettingsNotFound {
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

func (h *Handler) forgotPasswordRequest(w http.ResponseWriter, r *http.Request) {
	var payload types.ForgotPasswordRequestPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.db.GetUserByEmail(payload.Email)
	if err != nil {
		if err == types.ErrUserNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, types.ErrInvalidCredentials)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}
	if !user.EmailVerified {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrEmailNotVerified)
		return
	}

	expirationInMin := config.Env.ForgotPasswordTokenExpirationInMin

	resetToken, _, err := h.authHandler.GenerateToken(user.Id, expirationInMin)

	resetLink := fmt.Sprintf("%s?token=%s", config.Env.ResetPasswordWebsitePageUrl, resetToken)

	err = h.smtpServer.SendPasswordResetRequestMail(
		user.FullName.String,
		user.Email,
		resetLink,
		config.Env.WebsiteName,
		config.Env.WebsiteUrl,
		int(expirationInMin),
	)
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusInternalServerError,
			types.ErrOnSendingMail,
		)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) resetPassword(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			types.ErrTokenIsMissing,
		)
		return
	}

	var payload types.ResetPasswordPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	claims := types.UserJWTClaims{}

	_, err = h.authHandler.ValidateToken(token, &claims)
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			types.ErrInvalidTokenReceived,
		)
		return
	}

	user, err := h.db.GetUserById(claims.UserId)
	if err != nil {
		if err == types.ErrUserNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, types.ErrInvalidCredentials)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	hashedPassword, err := auth.HashPassword(payload.NewPassword)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = h.db.UpdateUser(user.Id, types.UpdateUserPayload{
		Password: &hashedPassword,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) verifyEmailRequest(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.db.GetUserById(userId)
	if err != nil {
		if err == types.ErrUserNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, types.ErrInvalidCredentials)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}
	if user.EmailVerified {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrEmailAlreadyVerified)
		return
	}

	expirationInMin := config.Env.EmailVerificationTokenExpirationInMin

	verificationToken, _, err := h.authHandler.GenerateToken(user.Id, expirationInMin)

	verificationLink := fmt.Sprintf(
		"%s?token=%s",
		config.Env.EmailVerificationWebsitePageUrl,
		verificationToken,
	)

	err = h.smtpServer.SendEmailVerificationRequestMail(
		user.FullName.String,
		user.Email,
		verificationLink,
		config.Env.WebsiteName,
		config.Env.WebsiteUrl,
		int(expirationInMin),
	)
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusInternalServerError,
			types.ErrOnSendingMail,
		)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

func (h *Handler) verifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			types.ErrTokenIsMissing,
		)
		return
	}

	claims := types.UserJWTClaims{}

	_, err := h.authHandler.ValidateToken(token, &claims)
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			types.ErrInvalidTokenReceived,
		)
		return
	}

	user, err := h.db.GetUserById(claims.UserId)
	if err != nil {
		if err == types.ErrUserNotFound {
			utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidCredentials)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	err = h.db.UpdateUser(user.Id, types.UpdateUserPayload{
		EmailVerified: utils.Ptr(true),
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}
