package types

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound                    = errors.New("user not found")
	ErrStoreNotFound                   = errors.New("store not found")
	ErrRoleNotFound                    = errors.New("role not found")
	ErrWalletNotFound                  = errors.New("wallet not found")
	ErrWalletTransactionNotFound       = errors.New("wallet transaction not found")
	ErrProductNotFound                 = errors.New("product not found")
	ErrSubcategoryNotFound             = errors.New("subcategory not found")
	ErrProductCategoryNotFound         = errors.New("product category not found")
	ErrProductTagNotFound              = errors.New("product tag not found")
	ErrProductOfferNotFound            = errors.New("product offer not found")
	ErrProductImageNotFound            = errors.New("product image not found")
	ErrProductSpecNotFound             = errors.New("product spec not found")
	ErrProductVariantNotFound          = errors.New("product variant not found")
	ErrProductCommentNotFound          = errors.New("product comment not found")
	ErrPermissionGroupNotFound         = errors.New("permission group not found")
	ErrUserSettingsNotFound            = errors.New("user settings not found")
	ErrStoreSettingsNotFound           = errors.New("store settings not found")
	ErrNoFieldsReceivedToUpdate        = errors.New("no fields received to update")
	ErrInvalidVisibilityStatusOption   = errors.New("invalid visibility status option")
	ErrInvalidVerificationStatusOption = errors.New("invalid verification status option")
	ErrReqBodyNotFound                 = errors.New("request body is not found")
	ErrInvalidUserPayload              = errors.New("invalid user payload")
	ErrInvalidLoginPayload             = errors.New("invalid login payload")
	ErrInconsistentAttributePresence   = errors.New(
		"inconsistent attribute presence: some but not all pvos have attribute_id",
	)
	ErrDuplicateUsernameOrEmail = errors.New("another user with this username/email already exists")
	ErrInternalServer           = errors.New("internal server error")
	ErrValidateTokenFailure     = errors.New("failed to validate token")
	ErrInvalidTokenReceived     = errors.New("invalid token received")
	ErrInvalidCredentials       = errors.New("invalid credentials received")
	ErrInvalidPayload           = func(err error) error {
		return errors.New(fmt.Sprintf("invalid payload: %v", err))
	}
	ErrInvalidPEMBlockForPrivateKey = errors.New("invalid PEM block for private key")
	ErrInvalidPEMBlockForPublicKey  = errors.New("invalid PEM block for public key")
	ErrKIDHeaderMissing             = errors.New("missing kid in token header")
	ErrUnexpectedSigningMethod      = func(alg any) error {
		return errors.New(fmt.Sprintf("unexpected signing method: %v", alg))
	}
	ErrPubKeyIdNotFound                  = errors.New("public key not found")
	ErrAuthenticationCredentialsNotFound = errors.New("authentication credentials not found")
	ErrAccessDenied                      = errors.New("access denied")
)
