package types

import (
	"time"

	json_types "github.com/SaeedAlian/econest/api/types/json"
)

// Store represents a merchant store
// @model Store
type Store struct {
	// Unique identifier for the store
	Id int `json:"id"          exposure:"public"`
	// Name of the store
	Name string `json:"name"        exposure:"public"`
	// Description of the store
	Description string `json:"description" exposure:"public"`
	// Whether the store is verified (private, needs permission)
	Verified bool `json:"verified"    exposure:"private,needPermission"`
	// When the store was created
	CreatedAt time.Time `json:"createdAt"   exposure:"public"`
	// When the store was last updated
	UpdatedAt time.Time `json:"updatedAt"   exposure:"public"`
	// ID of the store owner (visibility depends on settings, needs permission)
	OwnerId int `json:"ownerId"     exposure:"publicOwner,needPermission"`
}

// StoreInfo represents basic store information
// @model StoreInfo
type StoreInfo struct {
	// Store ID
	Id int `json:"id"   exposure:"public"`
	// Store name
	Name string `json:"name" exposure:"public"`
}

// StoreSettings contains store configuration
// @model StoreSettings
type StoreSettings struct {
	// Settings ID
	Id int `json:"id"          exposure:"public"`
	// Whether owner information is public
	PublicOwner bool `json:"publicOwner" exposure:"public"`
	// When settings were last updated
	UpdatedAt time.Time `json:"updatedAt"   exposure:"public"`
	// ID of the store these settings belong to
	StoreId int `json:"storeId"     exposure:"public"`
}

// StoreWithSettings combines Store with its settings
// @model StoreWithSettings
type StoreWithSettings struct {
	Store
	// Settings ID
	SettingsId int `json:"settingsId"        exposure:"public"`
	// Whether owner information is public
	PublicOwner bool `json:"publicOwner"       exposure:"public"`
	// When settings were last updated
	SettingsUpdatedAt time.Time `json:"settingsUpdatedAt" exposure:"public"`
}

// StorePhoneNumber represents a store's contact number
// @model StorePhoneNumber
type StorePhoneNumber struct {
	// Phone number ID
	Id int `json:"id"          exposure:"public"`
	// Country code (visibility depends on isPublic, needs permission)
	CountryCode string `json:"countryCode" exposure:"isPublic,needPermission"`
	// Phone number (visibility depends on isPublic, needs permission)
	Number string `json:"number"      exposure:"isPublic,needPermission"`
	// Whether the number is publicly visible
	IsPublic bool `json:"isPublic"    exposure:"public"`
	// Whether the number is verified (visibility depends on isPublic, needs permission)
	Verified bool `json:"verified"    exposure:"isPublic,needPermission"`
	// When the number was added
	CreatedAt time.Time `json:"createdAt"   exposure:"public"`
	// When the number was last updated
	UpdatedAt time.Time `json:"updatedAt"   exposure:"public"`
	// ID of the store this number belongs to
	StoreId int `json:"storeId"     exposure:"public"`
}

// StoreAddress represents a store's physical location
// @model StoreAddress
type StoreAddress struct {
	// Address ID
	Id int `json:"id"        exposure:"public"`
	// State/Province (visibility depends on isPublic, needs permission)
	State string `json:"state"     exposure:"isPublic,needPermission"`
	// City (visibility depends on isPublic, needs permission)
	City string `json:"city"      exposure:"isPublic,needPermission"`
	// Street address (visibility depends on isPublic, needs permission)
	Street string `json:"street"    exposure:"isPublic,needPermission"`
	// Zip/Postal code (visibility depends on isPublic, needs permission)
	Zipcode string `json:"zipcode"   exposure:"isPublic,needPermission"`
	// Additional details (visibility depends on isPublic, needs permission)
	Details json_types.JSONNullString `json:"details"   exposure:"isPublic,needPermission" swaggertype:"string"`
	// Whether the address is publicly visible
	IsPublic bool `json:"isPublic"  exposure:"public"`
	// When the address was added
	CreatedAt time.Time `json:"createdAt" exposure:"public"`
	// When the address was last updated
	UpdatedAt time.Time `json:"updatedAt" exposure:"public"`
	// ID of the store this address belongs to
	StoreId int `json:"storeId"   exposure:"public"`
}

// StoreOwnedProduct links products to their store
// @model StoreOwnedProduct
type StoreOwnedProduct struct {
	// ID of the store that owns the product
	StoreId int `json:"storeId"   exposure:"public"`
	// ID of the product
	ProductId int `json:"productId" exposure:"public"`
}

// CreateStorePayload contains data needed to create a new store
// @model CreateStorePayload
type CreateStorePayload struct {
	// Store name (required)
	Name string `json:"name"        validate:"required"`
	// Store description
	Description string `json:"description"`
	// ID of the store owner
	OwnerId int `json:"ownerId"`
}

// CreateStorePhoneNumberPayload contains data needed to add a store phone number
// @model CreateStorePhoneNumberPayload
type CreateStorePhoneNumberPayload struct {
	// Country code (1-4 characters, required)
	CountryCode string `json:"countryCode" validate:"required,min=1,max=4"`
	// Phone number (5-20 characters, required)
	Number string `json:"number"      validate:"required,min=5,max=20"`
	// ID of the store this number belongs to
	StoreId int `json:"storeId"`
}

// CreateStoreAddressPayload contains data needed to add a store address
// @model CreateStoreAddressPayload
type CreateStoreAddressPayload struct {
	// State/Province (required)
	State string `json:"state"   validate:"required"`
	// City (required)
	City string `json:"city"    validate:"required"`
	// Street address (required)
	Street string `json:"street"  validate:"required"`
	// Zip/Postal code (required)
	Zipcode string `json:"zipcode" validate:"required"`
	// Additional address details
	Details string `json:"details"`
	// ID of the store this address belongs to
	StoreId int `json:"storeId"`
}

// StoreSearchQuery contains parameters for searching stores
// @model StoreSearchQuery
type StoreSearchQuery struct {
	// Filter by store name
	Name *string `json:"name"`
	// Filter by owner ID
	OwnerId *int `json:"ownerId"`
	// Maximum number of results
	Limit *int `json:"limit"`
	// Number of results to skip
	Offset *int `json:"offset"`
}

// UpdateStorePhoneNumberPayload contains data for updating a store phone number
// @model UpdateStorePhoneNumberPayload
type UpdateStorePhoneNumberPayload struct {
	// New country code (1-4 characters)
	CountryCode *string `json:"countryCode" validate:"omitempty,min=1,max=4"`
	// New phone number (5-20 characters)
	Number *string `json:"number"      validate:"omitempty,min=5,max=20"`
	// New visibility status
	IsPublic *bool `json:"isPublic"`
	// New verification status
	Verified *bool `json:"verified"`
}

// UpdateStoreAddressPayload contains data for updating a store address
// @model UpdateStoreAddressPayload
type UpdateStoreAddressPayload struct {
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
	// New visibility status
	IsPublic *bool `json:"isPublic"`
}

// StorePhoneNumberSearchQuery contains parameters for searching store phone numbers
// @model StorePhoneNumberSearchQuery
type StorePhoneNumberSearchQuery struct {
	// Filter by visibility status
	VisibilityStatus *SettingVisibilityStatus `json:"visibilityStatus"`
	// Filter by verification status
	VerificationStatus *CredentialVerificationStatus `json:"verificationStatus"`
}

// StoreAddressSearchQuery contains parameters for searching store addresses
// @model StoreAddressSearchQuery
type StoreAddressSearchQuery struct {
	// Filter by visibility status
	VisibilityStatus *SettingVisibilityStatus `json:"visibilityStatus"`
}

// UpdateStorePayload contains data for updating store information
// @model UpdateStorePayload
type UpdateStorePayload struct {
	// New store name
	Name *string `json:"name"`
	// New store description
	Description *string `json:"description"`
	// New verification status
	Verified *bool `json:"verified"`
}

// UpdateStoreSettingsPayload contains data for updating store settings
// @model UpdateStoreSettingsPayload
type UpdateStoreSettingsPayload struct {
	// New owner visibility setting
	PublicOwner *bool `json:"publicOwner"`
}
