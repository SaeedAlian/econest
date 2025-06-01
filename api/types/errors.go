package types

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound                   = errors.New("user not found")
	ErrUserAddressNotFound            = errors.New("user address not found")
	ErrStoreAddressNotFound           = errors.New("store address not found")
	ErrShipmentAddressNotFound        = errors.New("shipment address not found")
	ErrUserPhoneNumberNotFound        = errors.New("user phone number not found")
	ErrStorePhoneNumberNotFound       = errors.New("store phone number not found")
	ErrStoreNotFound                  = errors.New("store not found")
	ErrRoleNotFound                   = errors.New("role not found")
	ErrWalletNotFound                 = errors.New("wallet not found")
	ErrWalletTransactionNotFound      = errors.New("wallet transaction not found")
	ErrProductNotFound                = errors.New("product not found")
	ErrSubcategoryNotFound            = errors.New("subcategory not found")
	ErrProductCategoryNotFound        = errors.New("product category not found")
	ErrProductTagNotFound             = errors.New("product tag not found")
	ErrProductOfferNotFound           = errors.New("product offer not found")
	ErrProductImageNotFound           = errors.New("product image not found")
	ErrProductAttributeNotFound       = errors.New("product attribute not found")
	ErrProductAttributeOptionNotFound = errors.New("product attribute option not found")
	ErrProductSpecNotFound            = errors.New("product spec not found")
	ErrProductVariantNotFound         = errors.New("product variant not found")
	ErrProductCommentNotFound         = errors.New("product comment not found")
	ErrPermissionGroupNotFound        = errors.New("permission group not found")
	ErrUserSettingsNotFound           = errors.New("user settings not found")
	ErrStoreSettingsNotFound          = errors.New("store settings not found")
	ErrStoreOwnerNotFound             = errors.New("store owner not found")
	ErrOrderNotFound                  = errors.New("order not found")
	ErrForeignKeyViolationForColumn   = errors.New(
		"invalid reference: a related record does not exist",
	)

	ErrReqBodyNotFound                   = errors.New("request body is not found")
	ErrPubKeyIdNotFound                  = errors.New("public key not found")
	ErrAuthenticationCredentialsNotFound = errors.New("authentication credentials not found")
	ErrRefreshTokenNotFound              = errors.New("refresh token not found")
	ErrValidateTokenFailure              = errors.New("failed to validate token")
	ErrInvalidTokenReceived              = errors.New("invalid token received")
	ErrInvalidPEMBlockForPrivateKey      = errors.New("invalid PEM block for private key")
	ErrInvalidPEMBlockForPublicKey       = errors.New("invalid PEM block for public key")
	ErrKIDHeaderMissing                  = errors.New("missing kid in token header")
	ErrUnexpectedSigningMethod           = func(alg any) error {
		return errors.New(fmt.Sprintf("unexpected signing method: %v", alg))
	}
	ErrAccessDenied        = errors.New("access denied")
	ErrCSRFMissing         = errors.New("csrf token is missing")
	ErrInvalidCSRFToken    = errors.New("csrf token is invalid")
	ErrInvalidRefreshToken = errors.New("refresh token is invalid")
	ErrTokenIsMissing      = errors.New("token is missing")

	ErrInvalidOptionId  = errors.New("invalid option id")
	ErrInvalidPageQuery = errors.New("invalid page")

	ErrProductQuantityIsNotEnough = errors.New("product quantity is not enough")
	ErrProductVariantsAreEmpty    = errors.New("product variants are empty")
	ErrBalanceInsufficient        = errors.New("insufficient wallet balance")

	ErrInvalidCredentials  = errors.New("invalid credentials received")
	ErrInvalidPayload      = errors.New("invalid payload received")
	ErrInvalidPayloadField = func(err error) error {
		return errors.New(fmt.Sprintf("invalid payload: %v", err))
	}

	ErrDuplicateRoleName            = errors.New("another role with this name already exists")
	ErrDuplicatePermissionGroupName = errors.New(
		"another permission group with this name already exists",
	)
	ErrDuplicateStoreName = errors.New(
		"another store with this name already exists",
	)
	ErrDuplicatePhoneNumber = errors.New(
		"another phone number with this number already exists",
	)
	ErrDuplicateProductCategoryImageName = errors.New(
		"another product category image with this name already exists",
	)
	ErrDuplicateProductImageName = errors.New(
		"another product image with this name already exists",
	)
	ErrDuplicateUserEmail = errors.New(
		"another user with this email already exists",
	)
	ErrDuplicateUsername = errors.New(
		"another user with this username already exists",
	)
	ErrDuplicateProductSlug = errors.New(
		"another product with this slug already exists",
	)
	ErrUniqueConstraintViolation          = errors.New("a unique constraint has been violated")
	ErrUniqueConstraintViolationForColumn = func(col string) error {
		return errors.New(fmt.Sprintf("the value for '%s' must be unique.", col))
	}

	ErrInvalidActionEnum               = errors.New("invalid action specified")
	ErrInvalidResourceEnum             = errors.New("invalid resource specified")
	ErrInvalidTransactionTypeEnum      = errors.New("invalid transaction type specified")
	ErrInvalidTransactionStatusEnum    = errors.New("invalid transaction status specified")
	ErrInvalidOrderPaymentStatusEnum   = errors.New("invalid order payment status specified")
	ErrInvalidOrderShipmentStatusEnum  = errors.New("invalid order shipment status specified")
	ErrInvalidVisibilityStatusOption   = errors.New("invalid visibility status option")
	ErrInvalidVerificationStatusOption = errors.New("invalid verification status option")
	ErrInvalidInputFormat              = errors.New("invalid input format")

	ErrNoFieldsReceivedToUpdate = errors.New("no fields received to update")

	ErrCannotAccessAddress           = errors.New("you cannot access this address")
	ErrCannotAccessPhoneNumber       = errors.New("you cannot access this phone number")
	ErrCannotAccessWalletTransaction = errors.New("you cannot access this wallet transaction")
	ErrCannotBanThisUser             = errors.New("you cannot ban this user")
	ErrCannotAccessStore             = errors.New("you cannot access this store")
	ErrCannotAccessComment           = errors.New("you cannot access this comment")
	ErrTransactionIsNotForWallet     = errors.New(
		"this transaction is not for the provided user wallet",
	)

	ErrUserIsBanned         = errors.New("this user is banned")
	ErrEmailNotVerified     = errors.New("email is not verified yet")
	ErrEmailAlreadyVerified = errors.New("email is already verified")

	ErrOnSendingMail = errors.New("error on sending mail")

	ErrImageNotExist                = errors.New("image file does not exist")
	ErrCouldNotOpenFile             = errors.New("could not open file")
	ErrCouldNotGetFileMimeType      = errors.New("could not get the file type")
	ErrCouldNotResetFileReader      = errors.New("could not read the file")
	ErrCouldNotGetFileStats         = errors.New("could not get the file information")
	ErrCouldNotCopyFileIntoResponse = errors.New("could not send the file")

	ErrCannotLoginWithThisUser = errors.New("cannot login with this user")
	ErrCannotRegisterThisUser  = errors.New("cannot login with this user")

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
