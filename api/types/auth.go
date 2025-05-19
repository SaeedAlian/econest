package types

import "github.com/golang-jwt/jwt/v5"

type LoginResponsePayload struct {
	AccessToken string `json:"accessToken"`
}

type LoginUserPayload struct {
	Username string `json:"username" validate:"required,min=5"`
	Password string `json:"password" validate:"required,min=6,max=130"`
}

type UserJWTClaims struct {
	UserId    int    `json:"userId"`
	ExpiresAt int64  `json:"expiresAt"`
	IssuedAt  int64  `json:"issuedAt"`
	JTI       string `json:"jti"`
}

func (c *UserJWTClaims) PopulateFromToken(claims jwt.MapClaims) error {
	c.UserId = claims["userId"].(int)
	c.ExpiresAt = int64(claims["exp"].(float64))
	c.IssuedAt = int64(claims["iat"].(float64))
	c.JTI = claims["jti"].(string)
	return nil
}
