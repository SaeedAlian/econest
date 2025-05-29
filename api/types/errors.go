package types

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound                    = errors.New("user not found")
	ErrUserAddressNotFound             = errors.New("user address not found")
	ErrStoreAddressNotFound            = errors.New("store address not found")
	ErrUserPhoneNumberNotFound         = errors.New("user phone number not found")
	ErrStorePhoneNumberNotFound        = errors.New("store phone number not found")
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
	ErrProductAttributeNotFound        = errors.New("product attribute not found")
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
	ErrInconsistentAttributePresence   = errors.New(
		"inconsistent attribute presence: some but not all pvos have attribute_id",
	)
	ErrDuplicateUsernameOrEmail = errors.New("another user with this username/email already exists")
	ErrDuplicateUsername        = errors.New("another user with this username already exists")
	ErrDuplicateEmail           = errors.New("another user with this email already exists")
	ErrInternalServer           = errors.New("internal server error")
	ErrValidateTokenFailure     = errors.New("failed to validate token")
	ErrInvalidTokenReceived     = errors.New("invalid token received")
	ErrInvalidCredentials       = errors.New("invalid credentials received")
	ErrInvalidPayload           = errors.New("invalid payload received")
	ErrInvalidPayloadField      = func(err error) error {
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
	ErrCSRFMissing                       = errors.New("csrf token is missing")
	ErrInvalidCSRFToken                  = errors.New("csrf token is invalid")
	ErrRefreshTokenNotFound              = errors.New("refresh token not found")
	ErrInvalidRefreshToken               = errors.New("refresh token is invalid")
	ErrCannotAccessAddress               = errors.New("you cannot access this address")
	ErrCannotAccessPhoneNumber           = errors.New("you cannot access this phone number")
	ErrUpdateAddress                     = errors.New("there was an error on updating address")
	ErrUpdatePhoneNumber                 = errors.New(
		"there was an error on updating phone number",
	)
	ErrDeleteAddress     = errors.New("there was an error on deleting address")
	ErrDeletePhoneNumber = errors.New(
		"there was an error on deleting phone number",
	)
	ErrInvalidPageQuery     = errors.New("invalid page")
	ErrCannotBanThisUser    = errors.New("you cannot ban this user")
	ErrEmailNotVerified     = errors.New("email is not verified yet")
	ErrOnSendingMail        = errors.New("error on sending mail")
	ErrTokenIsMissing       = errors.New("token is missing")
	ErrEmailAlreadyVerified = errors.New("email is already verified")
	ErrUserIsBanned         = errors.New("this user is banned")
	ErrCannotAccessStore    = errors.New("you cannot access this store")
	ErrCannotAccessComment  = errors.New("you cannot access this comment")
	ErrImageNotExist        = errors.New("image file does not exist")
	ErrDuplicateStoreName   = errors.New(
		"another store with this name already exists",
	)
	ErrInvalidOptionId  = errors.New("invalid option id")
	ErrUploadSizeTooBig = func(maxSize int) error {
		return errors.New(
			fmt.Sprintf(
				"uploaded file is too big, choose a file that's less than %d MB in size",
				maxSize,
			),
		)
	}
	ErrCannotRetrieveFile = func(err error) error {
		return errors.New(
			fmt.Sprintf(
				"cannot retrieve the file: %v",
				err,
			),
		)
	}
	ErrFileUpload = func(err error) error {
		return errors.New(
			fmt.Sprintf(
				"error in file uploading: %v",
				err,
			),
		)
	}
	ErrNotAllowedFileType = func(allowedFileTypes string) error {
		return errors.New(
			fmt.Sprintf(
				"cannot upload this file, please upload only %s files",
				allowedFileTypes,
			),
		)
	}
	ErrQueryMappingNilValueReceived = func(key string) error {
		return errors.New(
			fmt.Sprintf(
				"value for key %s is not a non-nil pointer",
				key,
			),
		)
	}
	ErrInvalidQueryValue = func(key string) error {
		return errors.New(
			fmt.Sprintf(
				"invalid value for %s query",
				key,
			),
		)
	}
	ErrInvalidParamValue = func(key string) error {
		return errors.New(
			fmt.Sprintf(
				"invalid value for %s parameter",
				key,
			),
		)
	}
)
