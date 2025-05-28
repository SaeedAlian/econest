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
	ErrInvalidUserPayload              = errors.New("invalid user payload")
	ErrInvalidLoginPayload             = errors.New("invalid login payload")
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
	ErrCSRFMissing                       = errors.New("csrf token is missing")
	ErrInvalidCSRFToken                  = errors.New("csrf token is invalid")
	ErrRefreshTokenNotFound              = errors.New("refresh token not found")
	ErrInvalidRefreshToken               = errors.New("refresh token is invalid")
	ErrCreateAddress                     = errors.New("there was an error on creating address")
	ErrCreatePhoneNumber                 = errors.New(
		"there was an error on creating phone number",
	)
	ErrCannotAccessAddress     = errors.New("you cannot access this address")
	ErrCannotAccessPhoneNumber = errors.New("you cannot access this phone number")
	ErrInvalidAddressId        = errors.New("invalid address id")
	ErrInvalidPhoneNumberId    = errors.New("invalid phone number id")
	ErrUpdateAddress           = errors.New("there was an error on updating address")
	ErrUpdatePhoneNumber       = errors.New(
		"there was an error on updating phone number",
	)
	ErrInvalidAddressPayload     = errors.New("invalid address payload")
	ErrInvalidPhoneNumberPayload = errors.New("invalid phone number payload")
	ErrDeleteAddress             = errors.New("there was an error on deleting address")
	ErrDeletePhoneNumber         = errors.New(
		"there was an error on deleting phone number",
	)
	ErrInvalidUserId                       = errors.New("invalid user id")
	ErrInvalidRoleIdQuery                  = errors.New("invalid role id")
	ErrInvalidPageQuery                    = errors.New("invalid page")
	ErrInvalidProfilePayload               = errors.New("invalid profile payload")
	ErrInvalidUserSettingsPayload          = errors.New("invalid settings payload")
	ErrCannotBanThisUser                   = errors.New("you cannot ban this user")
	ErrEmailNotVerified                    = errors.New("email is not verified yet")
	ErrOnSendingMail                       = errors.New("error on sending mail")
	ErrInvalidForgotPasswordRequestPayload = errors.New("invalid password payload")
	ErrTokenIsMissing                      = errors.New("token is missing")
	ErrInvalidResetPasswordPayload         = errors.New("invalid password payload")
	ErrEmailAlreadyVerified                = errors.New("email is already verified")
	ErrInvalidPasswordPayload              = errors.New("invalid password payload")
	ErrUserIsBanned                        = errors.New("this user is banned")
	ErrInvalidStorePayload                 = errors.New("invalid store payload")
	ErrInvalidStoreId                      = errors.New("invalid store id")
	ErrInvalidOwnerIdQuery                 = errors.New("invalid owner id")
	ErrCannotAccessStore                   = errors.New("you cannot access this store")
	ErrDuplicateStoreName                  = errors.New(
		"another store with this name already exists",
	)
	ErrInvalidStoreSettingsPayload = errors.New("invalid store settings payload")
	ErrInvalidOptionId             = errors.New("invalid option id")
	ErrUploadSizeTooBig            = func(maxSize int) error {
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
)
