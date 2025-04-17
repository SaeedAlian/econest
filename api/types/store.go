package types

import "time"

type Store struct {
	Id          int       `json:"id"          exposure:"public"`
	Name        string    `json:"name"        exposure:"public"`
	Description string    `json:"description" exposure:"public"`
	Verified    bool      `json:"verified"    exposure:"private,needPermission"`
	CreatedAt   time.Time `json:"createdAt"   exposure:"public"`
	UpdatedAt   time.Time `json:"updatedAt"   exposure:"public"`
	OwnerId     int       `json:"ownerId"     exposure:"publicOwner,needPermission"`
}

type StoreSettings struct {
	Id          int       `json:"id"          exposure:"public"`
	PublicOwner bool      `json:"publicOwner" exposure:"public"`
	UpdatedAt   time.Time `json:"updatedAt"   exposure:"public"`
	StoreId     int       `json:"storeId"     exposure:"public"`
}

type StoreWithSettings struct {
	Store
	SettingsId        int       `json:"settingsId"        exposure:"public"`
	PublicOwner       bool      `json:"publicOwner"       exposure:"public"`
	SettingsUpdatedAt time.Time `json:"settingsUpdatedAt" exposure:"public"`
}

type StorePhoneNumber struct {
	Id          int       `json:"id"          exposure:"public"`
	CountryCode string    `json:"countryCode" exposure:"isPublic,needPermission"`
	Number      string    `json:"number"      exposure:"isPublic,needPermission"`
	IsPublic    bool      `json:"isPublic"    exposure:"public"`
	Verified    bool      `json:"verified"    exposure:"isPublic,needPermission"`
	CreatedAt   time.Time `json:"createdAt"   exposure:"public"`
	UpdatedAt   time.Time `json:"updatedAt"   exposure:"public"`
	StoreId     int       `json:"storeId"     exposure:"public"`
}

type StoreAddress struct {
	Id        int       `json:"id"        exposure:"public"`
	State     string    `json:"state"     exposure:"isPublic,needPermission"`
	City      string    `json:"city"      exposure:"isPublic,needPermission"`
	Street    string    `json:"street"    exposure:"isPublic,needPermission"`
	Zipcode   string    `json:"zipcode"   exposure:"isPublic,needPermission"`
	Details   string    `json:"details"   exposure:"isPublic,needPermission"`
	IsPublic  bool      `json:"isPublic"  exposure:"public"`
	CreatedAt time.Time `json:"createdAt" exposure:"public"`
	UpdatedAt time.Time `json:"updatedAt" exposure:"public"`
	StoreId   int       `json:"storeId"   exposure:"public"`
}

type StoreOwnedProduct struct {
	StoreId   int `json:"storeId"   exposure:"public"`
	ProductId int `json:"productId" exposure:"public"`
}

type CreateStorePayload struct {
	Name        string `json:"name"        validate:"required"`
	Description string `json:"description"`
	OwnerId     int    `json:"ownerId"     validate:"required"`
}

type CreateStorePhoneNumberPayload struct {
	CountryCode string `json:"countryCode" validate:"required,min=1,max=4"`
	Number      string `json:"number"      validate:"required,min=5,max=20"`
	StoreId     int    `json:"storeId"     validate:"required"`
}

type CreateStoreAddressPayload struct {
	State   string `json:"state"   validate:"required"`
	City    string `json:"city"    validate:"required"`
	Street  string `json:"street"  validate:"required"`
	Zipcode string `json:"zipcode" validate:"required"`
	Details string `json:"details"`
	StoreId int    `json:"storeId" validate:"required"`
}

type StoreSearchQuery struct {
	Name   *string `json:"name"`
	Limit  *int    `json:"limit"`
	Offset *int    `json:"offset"`
}

type UpdateStorePhoneNumberPayload struct {
	CountryCode *string `json:"countryCode" validate:"min=1,max=4"`
	Number      *string `json:"number"      validate:"min=5,max=20"`
	IsPublic    *bool   `json:"isPublic"`
	Verified    *bool   `json:"verified"`
}

type UpdateStoreAddressPayload struct {
	State    *string `json:"state"`
	City     *string `json:"city"`
	Street   *string `json:"street"`
	Zipcode  *string `json:"zipcode"`
	Details  *string `json:"details"`
	IsPublic *bool   `json:"isPublic"`
}

type StorePhoneNumberSearchQuery struct {
	VisibilityStatus   *SettingVisibilityStatus      `json:"visibilityStatus"`
	VerificationStatus *CredentialVerificationStatus `json:"verificationStatus"`
}

type StoreAddressSearchQuery struct {
	VisibilityStatus *SettingVisibilityStatus `json:"visibilityStatus"`
}

type UpdateStorePayload struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Verified    *bool   `json:"verified"`
}

type UpdateStoreSettingsPayload struct {
	PublicOwner *bool `json:"publicOwner"`
}

type AssignProductToStorePayload struct {
	StoreId   int `json:"storeId"   validate:"required"`
	ProductId int `json:"productId" validate:"required"`
}
