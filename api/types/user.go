package types

import "time"

type User struct {
	Id            int       `json:"id"            exposure:"public"`
	Username      string    `json:"username"      exposure:"private,needPermission"`
	Email         string    `json:"email"         exposure:"publicEmail,needPermission"`
	EmailVerified bool      `json:"emailVerified" exposure:"publicEmail,needPermission"`
	Password      string    `json:"-"`
	FullName      string    `json:"fullName"      exposure:"public"`
	BirthDate     time.Time `json:"birthDate"     exposure:"publicBirthDate,needPermission"`
	IsBanned      bool      `json:"isBanned"      exposure:"public"`
	CreatedAt     time.Time `json:"createdAt"     exposure:"public"`
	UpdatedAt     time.Time `json:"updatedAt"     exposure:"public"`
	RoleId        int       `json:"roleId"        exposure:"private,needPermission"`
}

type UserSettings struct {
	Id               int       `json:"id"               exposure:"public"`
	PublicEmail      bool      `json:"publicEmail"      exposure:"public"`
	PublicBirthDate  bool      `json:"publicBirthDate"  exposure:"public"`
	IsUsingDarkTheme bool      `json:"isUsingDarkTheme" exposure:"public"`
	Language         string    `json:"language"         exposure:"public"`
	UpdatedAt        time.Time `json:"updatedAt"        exposure:"public"`
	UserId           int       `json:"userId"           exposure:"public"`
}

type UserWithSettings struct {
	User
	SettingsId        int       `json:"settingsId"        exposure:"public"`
	PublicEmail       bool      `json:"publicEmail"       exposure:"public"`
	PublicBirthDate   bool      `json:"publicBirthDate"   exposure:"public"`
	IsUsingDarkTheme  bool      `json:"isUsingDarkTheme"  exposure:"public"`
	Language          string    `json:"language"          exposure:"public"`
	SettingsUpdatedAt time.Time `json:"settingsUpdatedAt" exposure:"public"`
}

type UserPhoneNumber struct {
	Id          int       `json:"id"          exposure:"public"`
	CountryCode string    `json:"countryCode" exposure:"isPublic,needPermission"`
	Number      string    `json:"number"      exposure:"isPublic,needPermission"`
	IsPublic    bool      `json:"isPublic"    exposure:"public"`
	Verified    bool      `json:"verified"    exposure:"isPublic,needPermission"`
	CreatedAt   time.Time `json:"createdAt"   exposure:"public"`
	UpdatedAt   time.Time `json:"updatedAt"   exposure:"public"`
	UserId      int       `json:"userId"      exposure:"public"`
}

type UserAddress struct {
	Id        int       `json:"id"        exposure:"public"`
	State     string    `json:"state"     exposure:"isPublic,needPermission"`
	City      string    `json:"city"      exposure:"isPublic,needPermission"`
	Street    string    `json:"street"    exposure:"isPublic,needPermission"`
	Zipcode   string    `json:"zipcode"   exposure:"isPublic,needPermission"`
	Details   string    `json:"details"   exposure:"isPublic,needPermission"`
	IsPublic  bool      `json:"isPublic"  exposure:"public"`
	CreatedAt time.Time `json:"createdAt" exposure:"public"`
	UpdatedAt time.Time `json:"updatedAt" exposure:"public"`
	UserId    int       `json:"userId"    exposure:"public"`
}

type CreateUserPayload struct {
	Username  string    `json:"username"  validate:"required,min=5"`
	Email     string    `json:"email"     validate:"required,email"`
	Password  string    `json:"password"  validate:"required,min=6,max=130"`
	FullName  string    `json:"fullName"`
	BirthDate time.Time `json:"birthDate" validate:"required"`
	RoleId    int       `json:"roleId"`
}

type CreateUserPhoneNumberPayload struct {
	CountryCode string `json:"countryCode" validate:"required,min=1,max=4"`
	Number      string `json:"number"      validate:"required,min=5,max=20"`
	UserId      int    `json:"userId"      validate:"required"`
}

type CreateUserAddressPayload struct {
	State   string `json:"state"   validate:"required"`
	City    string `json:"city"    validate:"required"`
	Street  string `json:"street"  validate:"required"`
	Zipcode string `json:"zipcode" validate:"required"`
	Details string `json:"details"`
	UserId  int    `json:"userId"  validate:"required"`
}

type UserSearchQuery struct {
	FullName *string `json:"fullName"`
	RoleId   *int    `json:"roleId"`
	Limit    *int    `json:"limit"`
	Offset   *int    `json:"offset"`
}

type UserPhoneNumberSearchQuery struct {
	VisibilityStatus   *SettingVisibilityStatus      `json:"visibilityStatus"`
	VerificationStatus *CredentialVerificationStatus `json:"verificationStatus"`
}

type UserAddressSearchQuery struct {
	VisibilityStatus *SettingVisibilityStatus `json:"visibilityStatus"`
}

type UpdateUserPayload struct {
	Username      *string    `json:"username"      validate:"min=5"`
	Email         *string    `json:"email"         validate:"email"`
	EmailVerified *bool      `json:"emailVerified"`
	Password      *string    `json:"password"      validate:"min=6,max=130"`
	FullName      *string    `json:"fullName"`
	BirthDate     *time.Time `json:"birthDate"`
	IsBanned      *bool      `json:"isBanned"`
}

type UpdateUserSettingsPayload struct {
	PublicEmail      *bool   `json:"publicEmail"`
	PublicBirthDate  *bool   `json:"publicBirthDate"`
	IsUsingDarkTheme *bool   `json:"isUsingDarkTheme"`
	Language         *string `json:"language"`
}

type UpdateUserPhoneNumberPayload struct {
	CountryCode *string `json:"countryCode" validate:"min=1,max=4"`
	Number      *string `json:"number"      validate:"min=5,max=20"`
	IsPublic    *bool   `json:"isPublic"`
	Verified    *bool   `json:"verified"`
}

type UpdateUserAddressPayload struct {
	State    *string `json:"state"`
	City     *string `json:"city"`
	Street   *string `json:"street"`
	Zipcode  *string `json:"zipcode"`
	Details  *string `json:"details"`
	IsPublic *bool   `json:"isPublic"`
}
