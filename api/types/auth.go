package types

import "github.com/golang-jwt/jwt/v5"

// LoginResponsePayload contains the response data after successful login
// @model LoginResponsePayload
type LoginResponsePayload struct {
	// JWT access token for authenticated requests
	AccessToken string `json:"accessToken"`
}

// LogoutResponsePayload contains the response data after logout attempt
// @model LogoutResponsePayload
type LogoutResponsePayload struct {
	// Indicates whether the logout was successful
	Success bool `json:"success"`
}

// LoginUserPayload contains credentials for user authentication
// @model LoginUserPayload
type LoginUserPayload struct {
	// Username for login (5+ characters, required)
	Username string `json:"username" validate:"required,min=5"`
	// Password for login (6-130 characters, required)
	Password string `json:"password" validate:"required,min=6,max=130"`
}

// ForgotPasswordRequestPayload contains data for password reset request
// @model ForgotPasswordRequestPayload
type ForgotPasswordRequestPayload struct {
	// Email address to send reset instructions to (valid email, required)
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordPayload contains data for setting a new password
// @model ResetPasswordPayload
type ResetPasswordPayload struct {
	// New password to set (6-130 characters, required)
	NewPassword string `json:"newPassword" validate:"required,min=6,max=130"`
}

// UserJWTClaims represents the claims contained in a JWT token
// @model UserJWTClaims
type UserJWTClaims struct {
	// ID of the authenticated user
	UserId int `json:"userId"`
	// Unix timestamp when the token expires
	ExpiresAt int64 `json:"expiresAt"`
	// Unix timestamp when the token was issued
	IssuedAt int64 `json:"issuedAt"`
	// Unique identifier for the token (JWT ID)
	JTI string `json:"jti"`
}

// PopulateFromToken populates the claims from a JWT token
// This method is used to extract claims from a parsed JWT token
func (c *UserJWTClaims) PopulateFromToken(claims jwt.MapClaims) error {
	c.UserId = int(claims["sub"].(float64))
	c.ExpiresAt = int64(claims["exp"].(float64))
	c.IssuedAt = int64(claims["iat"].(float64))
	c.JTI = claims["jti"].(string)
	return nil
}
