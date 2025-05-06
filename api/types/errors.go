package types

import "errors"

var (
	ErrUserNotFound                    = errors.New("user not found")
	ErrStoreNotFound                   = errors.New("store not found")
	ErrRoleNotFound                    = errors.New("role not found")
	ErrWalletNotFound                  = errors.New("wallet not found")
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
	ErrInconsistentAttributePresence   = errors.New(
		"inconsistent attribute presence: some but not all pvos have attribute_id",
	)
)
