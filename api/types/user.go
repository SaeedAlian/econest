package types

import (
	"time"

	json_types "github.com/SaeedAlian/econest/api/types/json"
)

// User represents a user account
// @model User
type User struct {
	Id int `json:"id"            exposure:"public"`
	// Username of the user (private, needs permission)
	Username string `json:"username"      exposure:"private,needPermission"`
	// Email address of the user (visibility depends on settings, needs permission)
	Email string `json:"email"         exposure:"publicEmail,needPermission"`
	// Indicates if email is verified (visibility depends on settings, needs permission)
	EmailVerified bool `json:"emailVerified" exposure:"publicEmail,needPermission"`
	// User password (which doesn't include in the json)
	Password string `json:"-"`
	// Full name of the user
	FullName json_types.JSONNullString `json:"fullName"      exposure:"public"                         swaggertype:"string"`
	// Birth date of the user (visibility depends on settings, needs permission)
	BirthDate json_types.JSONNullTime `json:"birthDate"     exposure:"publicBirthDate,needPermission" swaggertype:"string" format:"date-time"`
	// Indicates if user is banned
	IsBanned bool `json:"isBanned"      exposure:"public"`
	// When the user was created
	CreatedAt time.Time `json:"createdAt"     exposure:"public"`
	// When the user was last updated
	UpdatedAt time.Time `json:"updatedAt"     exposure:"public"`
	// Role ID of the user (private, needs permission)
	RoleId int `json:"roleId"        exposure:"private,needPermission"`
}

// UserSettings represents user preferences and visibility settings
// @model UserSettings
type UserSettings struct {
	Id int `json:"id"               exposure:"public"`
	// Whether email is public
	PublicEmail bool `json:"publicEmail"      exposure:"public"`
	// Whether birth date is public
	PublicBirthDate bool `json:"publicBirthDate"  exposure:"public"`
	// Whether dark theme is enabled
	IsUsingDarkTheme bool `json:"isUsingDarkTheme" exposure:"public"`
	// User's preferred language
	Language string `json:"language"         exposure:"public"`
	// When settings were last updated
	UpdatedAt time.Time `json:"updatedAt"        exposure:"public"`
	// ID of the user these settings belong to
	UserId int `json:"userId"           exposure:"public"`
}

// UserWithSettings combines User with their settings
// @model UserWithSettings
type UserWithSettings struct {
	User
	SettingsId        int       `json:"settingsId"        exposure:"public"`
	PublicEmail       bool      `json:"publicEmail"       exposure:"public"`
	PublicBirthDate   bool      `json:"publicBirthDate"   exposure:"public"`
	IsUsingDarkTheme  bool      `json:"isUsingDarkTheme"  exposure:"public"`
	Language          string    `json:"language"          exposure:"public"`
	SettingsUpdatedAt time.Time `json:"settingsUpdatedAt" exposure:"public"`
}

// UserPhoneNumber represents a user's phone number
// @model UserPhoneNumber
type UserPhoneNumber struct {
	Id int `json:"id"          exposure:"public"`
	// Country code (visibility depends on isPublic, needs permission)
	CountryCode string `json:"countryCode" exposure:"isPublic,needPermission"`
	// Phone number (visibility depends on isPublic, needs permission)
	Number string `json:"number"      exposure:"isPublic,needPermission"`
	// Whether the phone number is public
	IsPublic bool `json:"isPublic"    exposure:"public"`
	// Whether the phone number is verified (visibility depends on isPublic, needs permission)
	Verified bool `json:"verified"    exposure:"isPublic,needPermission"`
	// When the phone number was added
	CreatedAt time.Time `json:"createdAt"   exposure:"public"`
	// When the phone number was last updated
	UpdatedAt time.Time `json:"updatedAt"   exposure:"public"`
	// ID of the user this phone number belongs to
	UserId int `json:"userId"      exposure:"public"`
}

// UserAddress represents a user's address
// @model UserAddress
type UserAddress struct {
	Id int `json:"id"        exposure:"public"`
	// State/Province (visibility depends on isPublic, needs permission)
	State string `json:"state"     exposure:"isPublic,needPermission"`
	// City (visibility depends on isPublic, needs permission)
	City string `json:"city"      exposure:"isPublic,needPermission"`
	// Street address (visibility depends on isPublic, needs permission)
	Street string `json:"street"    exposure:"isPublic,needPermission"`
	// Zip/Postal code (visibility depends on isPublic, needs permission)
	Zipcode string `json:"zipcode"   exposure:"isPublic,needPermission"`
	// Additional address details (visibility depends on isPublic, needs permission)
	Details json_types.JSONNullString `json:"details"   exposure:"isPublic,needPermission" swaggertype:"string"`
	// Whether the address is public
	IsPublic bool `json:"isPublic"  exposure:"public"`
	// When the address was added
	CreatedAt time.Time `json:"createdAt" exposure:"public"`
	// When the address was last updated
	UpdatedAt time.Time `json:"updatedAt" exposure:"public"`
	// ID of the user this address belongs to
	UserId int `json:"userId"    exposure:"public"`
}

// CommentUser represents minimal user info for comments
// @model CommentUser
type CommentUser struct {
	Id        int                       `json:"id"        exposure:"public"`
	FullName  json_types.JSONNullString `json:"fullName"  exposure:"public" swaggertype:"string"`
	CreatedAt time.Time                 `json:"createdAt" exposure:"public"`
	UpdatedAt time.Time                 `json:"updatedAt" exposure:"public"`
}

// CreateUserPayload contains data needed to create a new user
// @model CreateUserPayload
type CreateUserPayload struct {
	// Username (5+ characters)
	Username string `json:"username"  validate:"required,min=5"`
	// Email address
	Email string `json:"email"     validate:"required,email"`
	// Password (6-130 characters)
	Password string `json:"password"  validate:"required,min=6,max=130"`
	// Full name
	FullName string `json:"fullName"`
	// Birth date
	BirthDate time.Time `json:"birthDate" validate:"required"`
	// Role ID
	RoleId int `json:"roleId"`
}

// CreateUserPhoneNumberPayload contains data needed to add a phone number
// @model CreateUserPhoneNumberPayload
type CreateUserPhoneNumberPayload struct {
	// Country code (1-4 characters)
	CountryCode string `json:"countryCode" validate:"required,min=1,max=4"`
	// Phone number (5-20 characters)
	Number string `json:"number"      validate:"required,min=5,max=20"`
	// User ID this phone number belongs to
	UserId int `json:"userId"`
}

// CreateUserAddressPayload contains data needed to add an address
// @model CreateUserAddressPayload
type CreateUserAddressPayload struct {
	// State/Province
	State string `json:"state"   validate:"required"`
	// City
	City string `json:"city"    validate:"required"`
	// Street address
	Street string `json:"street"  validate:"required"`
	// Zip/Postal code
	Zipcode string `json:"zipcode" validate:"required"`
	// Additional address details
	Details string `json:"details"`
	// User ID this address belongs to
	UserId int `json:"userId"`
}

// UserSearchQuery contains parameters for searching users
// @model UserSearchQuery
type UserSearchQuery struct {
	// Filter by full name
	FullName *string `json:"fullName"`
	// Filter by role ID
	RoleId *int `json:"roleId"`
	// Maximum number of results
	Limit *int `json:"limit"`
	// Number of results to skip
	Offset *int `json:"offset"`
}

// UserPhoneNumberSearchQuery contains parameters for searching phone numbers
// @model UserPhoneNumberSearchQuery
type UserPhoneNumberSearchQuery struct {
	// Filter by visibility status
	VisibilityStatus *SettingVisibilityStatus `json:"visibilityStatus"`
	// Filter by verification status
	VerificationStatus *CredentialVerificationStatus `json:"verificationStatus"`
}

// UserAddressSearchQuery contains parameters for searching addresses
// @model UserAddressSearchQuery
type UserAddressSearchQuery struct {
	// Filter by visibility status
	VisibilityStatus *SettingVisibilityStatus `json:"visibilityStatus"`
}

// UpdateUserPayload contains data for updating a user
// @model UpdateUserPayload
type UpdateUserPayload struct {
	// New username (5+ characters)
	Username *string `json:"username"      validate:"omitempty,min=5"`
	// New email address
	Email *string `json:"email"         validate:"omitempty,email"`
	// Email verification status
	EmailVerified *bool `json:"emailVerified"`
	// New password (6-130 characters)
	Password *string `json:"password"      validate:"omitempty,min=6,max=130"`
	// New full name
	FullName *string `json:"fullName"`
	// New birth date
	BirthDate *time.Time `json:"birthDate"`
	// Ban status
	IsBanned *bool `json:"isBanned"`
}

// UpdateUserPasswordPayload contains data for changing a password
// @model UpdateUserPasswordPayload
type UpdateUserPasswordPayload struct {
	// Current password (6-130 characters)
	CurrentPassword *string `json:"currentPassword" validate:"min=6,max=130"`
	// New password (6-130 characters)
	NewPassword *string `json:"newPassword"     validate:"min=6,max=130"`
}

// UpdateUserSettingsPayload contains data for updating user settings
// @model UpdateUserSettingsPayload
type UpdateUserSettingsPayload struct {
	// Whether email should be public
	PublicEmail *bool `json:"publicEmail"`
	// Whether birth date should be public
	PublicBirthDate *bool `json:"publicBirthDate"`
	// Whether to use dark theme
	IsUsingDarkTheme *bool `json:"isUsingDarkTheme"`
	// Preferred language
	Language *string `json:"language"`
}

// UpdateUserPhoneNumberPayload contains data for updating a phone number
// @model UpdateUserPhoneNumberPayload
type UpdateUserPhoneNumberPayload struct {
	// New country code (1-4 characters)
	CountryCode *string `json:"countryCode" validate:"omitempty,min=1,max=4"`
	// New phone number (5-20 characters)
	Number *string `json:"number"      validate:"omitempty,min=5,max=20"`
	// Whether the number should be public
	IsPublic *bool `json:"isPublic"`
	// Verification status
	Verified *bool `json:"verified"`
}

// UpdateUserAddressPayload contains data for updating an address
// @model UpdateUserAddressPayload
type UpdateUserAddressPayload struct {
	// New state/province
	State *string `json:"state"`
	// New city
	City *string `json:"city"`
	// New street address
	Street *string `json:"street"`
	// New zip/postal code
	Zipcode *string `json:"zipcode"`
	// New additional details
	Details *string `json:"details"`
	// Whether the address should be public
	IsPublic *bool `json:"isPublic"`
}
