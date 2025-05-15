package auth

import (
	"context"
	"log"
	"net/http"
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

type JWTHandler struct {
	cache     *redis.Client
	keyServer *KeyServer
}

func NewJWTHandler(cache *redis.Client, keyServer *KeyServer) *JWTHandler {
	return &JWTHandler{cache: cache, keyServer: keyServer}
}

func (h *JWTHandler) WithJWTAuth(
	handler http.HandlerFunc,
	manager *db_manager.Manager,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")

		claims := types.UserJWTClaims{}
		token, err := h.ValidateToken(tokenStr, &claims)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			utils.WriteErrorInResponse(w, http.StatusUnauthorized, types.ErrValidateTokenFailure)
			return
		}

		if !token.Valid {
			log.Printf("invalid token received")
			utils.WriteErrorInResponse(w, http.StatusUnauthorized, types.ErrInvalidTokenReceived)
			return
		}

		userId := claims.UserId

		u, err := manager.GetUserById(userId)
		if u == nil || err != nil {
			log.Printf("invalid token received")
			utils.WriteErrorInResponse(w, http.StatusUnauthorized, types.ErrInvalidTokenReceived)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userId", u.Id)
		ctx = context.WithValue(ctx, "userRoleId", u.RoleId)
		r = r.WithContext(ctx)

		handler(w, r)
	}
}

func (h *JWTHandler) GenerateToken(userId int, expiresAtInMinutes float64) (string, error) {
	now := time.Now().UTC()
	expiration := time.Minute * time.Duration(expiresAtInMinutes)

	jti := uuid.NewString()

	tokenClaims := jwt.MapClaims{
		"sub": userId,
		"jti": jti,
		"iat": now.Unix(),
		"exp": now.Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaims)

	currentKID, err := h.keyServer.GetCurrentKID()
	if err != nil {
		return "", err
	}

	currentPK, err := h.keyServer.GetPrivateKey(currentKID)
	if err != nil {
		return "", err
	}

	token.Header["kid"] = currentKID

	tokenStr, err := token.SignedString(currentPK)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (h *JWTHandler) ValidateToken(token string, claims *types.UserJWTClaims) (*jwt.Token, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, types.ErrKIDHeaderMissing
		}

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
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

func (h *JWTHandler) SaveRefreshToken(jti string, userId int, exp int64) error {
	ttl := time.Duration(exp - time.Now().Unix())
	err := h.cache.Set(ctx, "refresh:"+jti, userId, ttl*time.Second).Err()
	return err
}

func (h *JWTHandler) IsRefreshTokenValid(jti string) (isValid bool, err error) {
	_, e := h.cache.Get(ctx, "refresh:"+jti).Result()
	if e == redis.Nil {
		return false, nil
	} else if e != nil {
		return false, e
	} else {
		return true, nil
	}
}

func (h *JWTHandler) RevokeRefreshToken(jti string) error {
	err := h.cache.Del(ctx, "refresh:"+jti).Err()
	return err
}

func (h *JWTHandler) RotateRefreshToken(oldJTI string, newJTI string, userId int, exp int64) error {
	if err := h.RevokeRefreshToken(oldJTI); err != nil {
		return err
	}
	return h.SaveRefreshToken(newJTI, userId, exp)
}
