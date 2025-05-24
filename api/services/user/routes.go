package user

import (
	"net/http"
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
	router.HandleFunc("/register/vendor", h.register("Vendor")).Methods("POST")
	router.HandleFunc("/register/customer", h.register("Customer")).Methods("POST")

	router.HandleFunc("/login", h.login).Methods("POST")
	router.HandleFunc("/refresh", h.refresh).Methods("POST")
	router.HandleFunc("/logout",
		h.authHandler.WithJWTAuth(
			h.authHandler.WithCSRFToken(h.logout), h.db,
		),
	).Methods("POST")

	router.HandleFunc("/me", h.authHandler.WithJWTAuth(h.getMe, h.db))
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
				types.ErrInvalidPayload(errors),
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
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidPayload(errors))
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
