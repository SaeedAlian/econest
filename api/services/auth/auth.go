package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	db_manager "github.com/SaeedAlian/econest/api/db/manager"
	"github.com/SaeedAlian/econest/api/types"
	"github.com/SaeedAlian/econest/api/utils"
)

type JWTClaims interface {
	PopulateFromToken(claims jwt.MapClaims) error
}

type AuthHandler struct {
	cache     *redis.Client
	keyServer *KeyServer
}

func NewAuthHandler(cache *redis.Client, keyServer *KeyServer) *AuthHandler {
	return &AuthHandler{cache: cache, keyServer: keyServer}
}

func (h *AuthHandler) WithJWTAuth(
	manager *db_manager.Manager,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteErrorInResponse(
					w,
					http.StatusUnauthorized,
					types.ErrValidateTokenFailure,
				)
				return
			}

			tokenStr := strings.Split(authHeader, "Bearer ")[1]

			claims := types.UserJWTClaims{}
			token, err := h.ValidateToken(tokenStr, &claims)
			if err != nil {
				utils.WriteErrorInResponse(
					w,
					http.StatusUnauthorized,
					types.ErrValidateTokenFailure,
				)
				return
			}

			if !token.Valid {
				utils.WriteErrorInResponse(
					w,
					http.StatusUnauthorized,
					types.ErrInvalidTokenReceived,
				)
				return
			}

			userId := claims.UserId

			u, err := manager.GetUserById(userId)
			if u == nil || err != nil {
				utils.WriteErrorInResponse(
					w,
					http.StatusUnauthorized,
					types.ErrInvalidTokenReceived,
				)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "userId", u.Id)
			ctx = context.WithValue(ctx, "userRoleId", u.RoleId)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func (h *AuthHandler) WithJWTAuthOptional(
	manager *db_manager.Manager,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := strings.Split(authHeader, "Bearer ")[1]

			claims := types.UserJWTClaims{}
			token, err := h.ValidateToken(tokenStr, &claims)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if !token.Valid {
				next.ServeHTTP(w, r)
				return
			}

			userId := claims.UserId

			u, err := manager.GetUserById(userId)
			if u == nil || err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "userId", u.Id)
			ctx = context.WithValue(ctx, "userRoleId", u.RoleId)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func (h *AuthHandler) WithCSRFToken() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" || r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			csrfHeader := r.Header.Get("X-CSRF-Token")
			if csrfHeader == "" {
				utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCSRFMissing)
				return
			}

			ctx := r.Context()
			userId := ctx.Value("userId")
			if userId == nil {
				http.Error(
					w,
					types.ErrAuthenticationCredentialsNotFound.Error(),
					http.StatusUnauthorized,
				)
				return
			}

			savedToken, isValid, err := h.GetCSRFToken(userId.(int))
			if err != nil {
				utils.WriteErrorInResponse(
					w,
					http.StatusInternalServerError,
					types.ErrInternalServer,
				)
				return
			}
			if !isValid || savedToken != csrfHeader {
				utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrInvalidCSRFToken)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (h *AuthHandler) WithVerifiedEmail(
	manager *db_manager.Manager,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" || r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()
			userId := ctx.Value("userId")
			if userId == nil {
				http.Error(
					w,
					types.ErrAuthenticationCredentialsNotFound.Error(),
					http.StatusUnauthorized,
				)
				return
			}

			user, err := manager.GetUserById(userId.(int))
			if err != nil {
				utils.WriteErrorInResponse(
					w,
					http.StatusInternalServerError,
					types.ErrInternalServer,
				)
				return
			}
			if !user.EmailVerified {
				utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrEmailNotVerified)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (h *AuthHandler) WithUnbannedProfile(
	manager *db_manager.Manager,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" || r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()
			userId := ctx.Value("userId")
			if userId == nil {
				http.Error(
					w,
					types.ErrAuthenticationCredentialsNotFound.Error(),
					http.StatusUnauthorized,
				)
				return
			}

			user, err := manager.GetUserById(userId.(int))
			if err != nil {
				utils.WriteErrorInResponse(
					w,
					http.StatusInternalServerError,
					types.ErrInternalServer,
				)
				return
			}
			if user.IsBanned {
				refreshToken, err := r.Cookie("refresh_token")

				if err == nil && refreshToken != nil {
					claims := types.UserJWTClaims{}
					_, err := h.ValidateToken(refreshToken.Value, &claims)
					if err == nil {
						h.RevokeRefreshToken(claims.JTI)
						h.RevokeCSRFToken(claims.UserId)
					}
				}

				utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrUserIsBanned)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (h *AuthHandler) WithActionPermissionAuth(
	handler http.HandlerFunc,
	manager *db_manager.Manager,
	actionPermissions []types.Action,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId := ctx.Value("userId")
		userRoleId := ctx.Value("userRoleId")

		if userId == nil || userRoleId == nil {
			utils.WriteErrorInResponse(
				w,
				http.StatusUnauthorized,
				types.ErrAuthenticationCredentialsNotFound,
			)
			return
		}

		acceptedRoles, err := manager.GetRolesBasedOnActionPermission(actionPermissions)
		if err != nil {
			utils.WriteErrorInResponse(
				w,
				http.StatusInternalServerError,
				types.ErrInternalServer,
			)
			return
		}

		found := false

		for _, r := range acceptedRoles {
			if r.Id == userRoleId.(int) {
				found = true
				break
			}
		}

		if !found {
			utils.WriteErrorInResponse(
				w,
				http.StatusForbidden,
				types.ErrAccessDenied,
			)
			return
		}

		handler(w, r)
	}
}

func (h *AuthHandler) WithResourcePermissionAuth(
	handler http.HandlerFunc,
	manager *db_manager.Manager,
	resourcePermissions []types.Resource,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId := ctx.Value("userId")
		userRoleId := ctx.Value("userRoleId")

		if userId == nil || userRoleId == nil {
			utils.WriteErrorInResponse(
				w,
				http.StatusUnauthorized,
				types.ErrAuthenticationCredentialsNotFound,
			)
			return
		}

		acceptedRoles, err := manager.GetRolesBasedOnResourcePermission(resourcePermissions)
		if err != nil {
			utils.WriteErrorInResponse(
				w,
				http.StatusInternalServerError,
				types.ErrInternalServer,
			)
			return
		}

		found := false

		for _, r := range acceptedRoles {
			if r.Id == userRoleId.(int) {
				found = true
				break
			}
		}

		if !found {
			utils.WriteErrorInResponse(
				w,
				http.StatusForbidden,
				types.ErrAccessDenied,
			)
			return
		}

		handler(w, r)
	}
}

func (h *AuthHandler) GenerateToken(
	userId int,
	expiresAtInMinutes float64,
) (token string, jti string, error error) {
	now := time.Now().UTC()
	expiration := time.Minute * time.Duration(expiresAtInMinutes)

	new_jti := uuid.NewString()

	tokenClaims := jwt.MapClaims{
		"sub": userId,
		"jti": new_jti,
		"iat": now.Unix(),
		"exp": now.Add(expiration).Unix(),
	}

	new_token := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaims)

	currentKID, err := h.keyServer.GetCurrentKID()
	if err != nil {
		return "", "", err
	}

	currentPK, err := h.keyServer.GetPrivateKey(currentKID)
	if err != nil {
		return "", "", err
	}

	new_token.Header["kid"] = currentKID

	tokenStr, err := new_token.SignedString(currentPK)
	if err != nil {
		return "", "", err
	}

	return tokenStr, new_jti, nil
}

func (h *AuthHandler) GenerateCSRFToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", nil
	}
	token := base64.RawURLEncoding.EncodeToString(bytes)
	return token, nil
}

func (h *AuthHandler) ValidateToken(token string, claims *types.UserJWTClaims) (*jwt.Token, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, types.ErrKIDHeaderMissing
		}

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, types.ErrUnexpectedSigningMethod(token.Header["alg"])
		}

		pubKey, err := h.keyServer.GetPublicKey(kid)
		if err != nil {
			return nil, err
		}
		if pubKey == nil {
			return nil, types.ErrPubKeyIdNotFound
		}

		return pubKey, nil
	})
	if err != nil {
		return nil, err
	}

	mapClaims := parsed.Claims.(jwt.MapClaims)

	if err := claims.PopulateFromToken(mapClaims); err != nil {
		return nil, err
	}

	return parsed, nil
}

func (h *AuthHandler) SaveRefreshToken(jti string, userId int, expiresAtInMinutes float64) error {
	ttl := time.Duration(expiresAtInMinutes) * time.Minute
	err := h.cache.Set(ctx, "refresh:"+jti, userId, ttl).Err()
	return err
}

func (h *AuthHandler) SaveCSRFToken(userId int, token string, expiresAtInMinutes float64) error {
	ttl := time.Duration(expiresAtInMinutes) * time.Minute
	err := h.cache.Set(ctx, fmt.Sprintf("csrf:%d", userId), token, ttl).Err()
	return err
}

func (h *AuthHandler) GetCSRFToken(userId int) (token string, isValid bool, err error) {
	t, e := h.cache.Get(ctx, fmt.Sprintf("csrf:%d", userId)).Result()
	if e == redis.Nil {
		return "", false, nil
	} else if e != nil {
		return "", false, e
	} else {
		return t, true, nil
	}
}

func (h *AuthHandler) IsRefreshTokenValid(jti string) (isValid bool, err error) {
	_, e := h.cache.Get(ctx, "refresh:"+jti).Result()
	if e == redis.Nil {
		return false, nil
	} else if e != nil {
		return false, e
	} else {
		return true, nil
	}
}

func (h *AuthHandler) RevokeRefreshToken(jti string) error {
	err := h.cache.Del(ctx, "refresh:"+jti).Err()
	return err
}

func (h *AuthHandler) RevokeCSRFToken(userId int) error {
	err := h.cache.Del(ctx, fmt.Sprintf("csrf:%d", userId)).Err()
	return err
}

func (h *AuthHandler) RotateRefreshToken(
	oldJTI string,
	newJTI string,
	userId int,
	expiresAtInMinutes float64,
) error {
	if err := h.RevokeRefreshToken(oldJTI); err != nil {
		return err
	}
	return h.SaveRefreshToken(newJTI, userId, expiresAtInMinutes)
}

func (h *AuthHandler) RotateCSRFToken(
	userId int,
	newToken string,
	expiresAtInMinutes float64,
) error {
	if err := h.RevokeCSRFToken(userId); err != nil {
		return err
	}
	return h.SaveCSRFToken(userId, newToken, expiresAtInMinutes)
}
