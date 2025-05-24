package auth

import (
	"database/sql"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/SaeedAlian/econest/api/types"
	json_types "github.com/SaeedAlian/econest/api/types/json"
)

func TestAuthHandler(t *testing.T) {
	ksCache := redis.NewClient(&redis.Options{
		Addr: testRedisAddr,
	})

	authHandlerCache := redis.NewClient(&redis.Options{
		Addr: testRedisAddr,
	})

	ks := NewKeyServer(ksCache)
	h := NewAuthHandler(authHandlerCache, ks)

	kid1 := "kid1"
	testUser := &types.User{
		Id:            1,
		Username:      "test_user",
		Email:         "test@email.com",
		EmailVerified: true,
		Password:      "passwd",
		FullName: json_types.JSONNullString{
			NullString: sql.NullString{
				String: "Test User",
			},
		},
		BirthDate: json_types.JSONNullTime{
			NullTime: sql.NullTime{
				Time: time.Now(),
			},
		},
		IsBanned:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("should initialize key server successfully", func(t *testing.T) {
		err := ks.RotateKeys(kid1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should generate token successfully", func(t *testing.T) {
		_, _, err := h.GenerateToken(testUser.Id, 1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should generate csrf token successfully", func(t *testing.T) {
		token, err := h.GenerateCSRFToken()
		if err != nil {
			t.Fatal(err)
		}

		if token == "" {
			t.Fatal("expected token, but got empty string")
		}
	})

	t.Run("should generate & verify token successfully", func(t *testing.T) {
		token, _, err := h.GenerateToken(testUser.Id, 1)
		if err != nil {
			t.Fatal(err)
		}

		claims := types.UserJWTClaims{}

		res, err := h.ValidateToken(token, &claims)
		if err != nil {
			t.Fatal(err)
		}

		if res == nil {
			t.Fatal("expected token result, but got nil")
		}

		if claims.UserId != testUser.Id {
			t.Fatal("expected test user id, but got another unrelated value")
		}
	})

	t.Run("should save refresh token", func(t *testing.T) {
		_, jti, err := h.GenerateToken(testUser.Id, 1)
		if err != nil {
			t.Fatal(err)
		}

		err = h.SaveRefreshToken(jti, testUser.Id, 1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should save csrf token", func(t *testing.T) {
		token, err := h.GenerateCSRFToken()
		if err != nil {
			t.Fatal(err)
		}

		err = h.SaveCSRFToken(testUser.Id, token, 1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should save & validate refresh token", func(t *testing.T) {
		_, jti, err := h.GenerateToken(testUser.Id, 1)
		if err != nil {
			t.Fatal(err)
		}

		err = h.SaveRefreshToken(jti, testUser.Id, 1)
		if err != nil {
			t.Fatal(err)
		}

		isValid, err := h.IsRefreshTokenValid(jti)
		if err != nil {
			t.Fatal(err)
		}

		if !isValid {
			t.Fatal("expected valid refresh token, but got invalid")
		}
	})

	t.Run("should save and get csrf token", func(t *testing.T) {
		token, err := h.GenerateCSRFToken()
		if err != nil {
			t.Fatal(err)
		}

		err = h.SaveCSRFToken(testUser.Id, token, 1)
		if err != nil {
			t.Fatal(err)
		}

		savedToken, isValid, err := h.GetCSRFToken(testUser.Id)
		if err != nil {
			t.Fatal(err)
		}
		if !isValid {
			t.Fatal("expected valid csrf token, but got invalid")
		}
		if savedToken != token {
			t.Fatal("expected saved token to be equal to the generated token, but it is not")
		}
	})

	t.Run("should save & revoke refresh token", func(t *testing.T) {
		_, jti, err := h.GenerateToken(testUser.Id, 1)
		if err != nil {
			t.Fatal(err)
		}

		err = h.SaveRefreshToken(jti, testUser.Id, 1)
		if err != nil {
			t.Fatal(err)
		}

		isValid, err := h.IsRefreshTokenValid(jti)
		if err != nil {
			t.Fatal(err)
		}

		if !isValid {
			t.Fatal("expected valid refresh token, but got invalid")
		}

		err = h.RevokeRefreshToken(jti)
		if err != nil {
			t.Fatal(err)
		}

		isValid, err = h.IsRefreshTokenValid(jti)
		if err != nil {
			t.Fatal(err)
		}

		if isValid {
			t.Fatal("expected not valid refresh token, but got valid")
		}
	})

	t.Run("should save and revoke csrf token", func(t *testing.T) {
		token, err := h.GenerateCSRFToken()
		if err != nil {
			t.Fatal(err)
		}

		err = h.SaveCSRFToken(testUser.Id, token, 1)
		if err != nil {
			t.Fatal(err)
		}

		savedToken, isValid, err := h.GetCSRFToken(testUser.Id)
		if err != nil {
			t.Fatal(err)
		}
		if !isValid {
			t.Fatal("expected valid csrf token, but got invalid")
		}
		if savedToken != token {
			t.Fatal("expected saved token to be equal to the generated token, but it is not")
		}

		err = h.RevokeCSRFToken(testUser.Id)
		if err != nil {
			t.Fatal(err)
		}

		_, isValid, err = h.GetCSRFToken(testUser.Id)
		if err != nil {
			t.Fatal(err)
		}
		if isValid {
			t.Fatal("expected not valid csrf token, but got valid")
		}
	})
}
