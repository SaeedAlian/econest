package types

import (
	"time"

	json_types "github.com/SaeedAlian/econest/api/types/json"
)

// ProductBase represents the core product information
// @model ProductBase
type ProductBase struct {
	// Unique product identifier (public)
	Id int `json:"id"             exposure:"public"`
	// Name of the product (public)
	Name string `json:"name"           exposure:"public"`
	// URL-friendly product identifier (public)
	Slug string `json:"slug"           exposure:"public"`
	// Current price of the product (public)
	Price float64 `json:"price"          exposure:"public"`
	// Factor used to calculate shipping costs (public)
	ShipmentFactor float64 `json:"shipmentFactor" exposure:"public"`
	// Detailed product description (public)
	Description string `json:"description"    exposure:"public"`
	// Whether the product is active/available (public)
	IsActive bool `json:"isActive"       exposure:"public"`
	// When the product was created (public)
	CreatedAt time.Time `json:"createdAt"      exposure:"public"`
	// When the product was last updated (public)
	UpdatedAt time.Time `json:"updatedAt"      exposure:"public"`
	// ID of the subcategory this product belongs to (public)
	SubcategoryId int `json:"subcategoryId"  exposure:"public"`
}

// ProductCategory represents a product category
// @model ProductCategory
type ProductCategory struct {
	// Unique category identifier (public)
	Id int `json:"id"               exposure:"public"`
	// Name of the category (public)
	Name string `json:"name"             exposure:"public"`
	// Filename of the category image (public)
	ImageName string `json:"imageName"        exposure:"public"`
	// When the category was created (public)
	CreatedAt time.Time `json:"createdAt"        exposure:"public"`
	// When the category was last updated (public)
	UpdatedAt time.Time `json:"updatedAt"        exposure:"public"`
	// ID of the parent category, if any (public)
	ParentCategoryId json_types.JSONNullInt32 `json:"parentCategoryId" exposure:"public" swaggertype:"primitive,number"`
}

// ProductCategoryWithParents represents a category with its parent hierarchy
// @model ProductCategoryWithParents
type ProductCategoryWithParents struct {
	ProductCategory
	// Parent category information (public, optional)
	ParentCategory *ProductCategoryWithParents `json:"parentCategory,omitempty" exposure:"public"`
}

// ProductOffer represents a special offer/discount for a product
// @model ProductOffer
type ProductOffer struct {
	// Unique offer identifier (public)
	Id int `json:"id"        exposure:"public"`
	// Discount percentage (public)
	Discount float64 `json:"discount"  exposure:"public"`
	// When the offer expires (public)
	ExpireAt time.Time `json:"expireAt"  exposure:"public"`
	// When the offer was created (public)
	CreatedAt time.Time `json:"createdAt" exposure:"public"`
	// When the offer was last updated (public)
	UpdatedAt time.Time `json:"updatedAt" exposure:"public"`
	// ID of the product this offer applies to (public)
	ProductId int `json:"productId" exposure:"public"`
}

// ProductImage represents an image associated with a product
// @model ProductImage
type ProductImage struct {
	// Unique image identifier (public)
	Id int `json:"id"        exposure:"public"`
	// Filename of the image (public)
	ImageName string `json:"imageName" exposure:"public"`
	// Whether this is the main product image (public)
	IsMain bool `json:"isMain"    exposure:"public"`
	// ID of the product this image belongs to (public)
	ProductId int `json:"productId" exposure:"public"`
}

// ProductSpec represents a product specification/feature
// @model ProductSpec
type ProductSpec struct {
	// Unique specification identifier (public)
	Id int `json:"id"        exposure:"public"`
	// Label/name of the specification (public)
	Label string `json:"label"     exposure:"public"`
	// Value of the specification (public)
	Value string `json:"value"     exposure:"public"`
	// ID of the product this spec belongs to (public)
	ProductId int `json:"productId" exposure:"public"`
}

// ProductTag represents a tag that can be assigned to products
// @model ProductTag
type ProductTag struct {
	// Unique tag identifier (public)
	Id int `json:"id"        exposure:"public"`
	// Name of the tag (public)
	Name string `json:"name"      exposure:"public"`
	// When the tag was created (public)
	CreatedAt time.Time `json:"createdAt" exposure:"public"`
	// When the tag was last updated (public)
	UpdatedAt time.Time `json:"updatedAt" exposure:"public"`
}

// ProductTagAssignment represents an assignment of a tag to a product
// @model ProductTagAssignment
type ProductTagAssignment struct {
	// ID of the product being tagged (public)
	ProductId int `json:"productId" exposure:"public"`
	// ID of the tag being assigned (public)
	TagId int `json:"tagId"     exposure:"public"`
}

// ProductAttribute represents a product attribute (like color, size)
// @model ProductAttribute
type ProductAttribute struct {
	// Unique attribute identifier (public)
	Id int `json:"id"    exposure:"public"`
	// Name/label of the attribute (public)
	Label string `json:"label" exposure:"public"`
}

// ProductAttributeOption represents an option for a product attribute
// @model ProductAttributeOption
type ProductAttributeOption struct {
	// Unique option identifier (public)
	Id int `json:"id"          exposure:"public"`
	// Value of the option (public)
	Value string `json:"value"       exposure:"public"`
	// ID of the attribute this option belongs to (public)
	AttributeId int `json:"attributeId" exposure:"public"`
}

// ProductAttributeWithOptions combines an attribute with its options
// @model ProductAttributeWithOptions
type ProductAttributeWithOptions struct {
	ProductAttribute
	// List of available options for this attribute (public)
	Options []ProductAttributeOption `json:"options" exposure:"public"`
}

// ProductVariant represents a specific variant of a product
// @model ProductVariant
type ProductVariant struct {
	// Unique variant identifier (public)
	Id int `json:"id"        exposure:"public"`
	// Current stock quantity (public)
	Quantity int `json:"quantity"  exposure:"public"`
	// ID of the product this variant belongs to (public)
	ProductId int `json:"productId" exposure:"public"`
}

// ProductVariantAttributeOption represents an attribute option assigned to a variant
// @model ProductVariantAttributeOption
type ProductVariantAttributeOption struct {
	// ID of the variant (public)
	VariantId int `json:"variantId"   exposure:"public"`
	// ID of the attribute (public)
	AttributeId int `json:"attributeId" exposure:"public"`
	// ID of the selected option (public)
	OptionId int `json:"optionId"    exposure:"public"`
}

// ProductVariantSelectedAttributeOption represents a selected attribute option for display
// @model ProductVariantSelectedAttributeOption
type ProductVariantSelectedAttributeOption struct {
	ProductAttribute
	// The selected option for this attribute (public)
	SelectedOption ProductAttributeOption `json:"selectedOption" exposure:"public"`
}

// ProductVariantWithAttributeSet represents a variant with its complete attribute set
// @model ProductVariantWithAttributeSet
type ProductVariantWithAttributeSet struct {
	ProductVariant
	// Complete set of attributes defining this variant (public)
	AttributeSet []ProductVariantSelectedAttributeOption `json:"attributeSet" exposure:"public"`
}

// ProductComment represents a user's comment/rating on a product
// @model ProductComment
type ProductComment struct {
	// Unique comment identifier (public)
	Id int `json:"id"        exposure:"public"`
	// Rating score (1-5) (public)
	Scoring int `json:"scoring"   exposure:"public"`
	// Text of the comment (public)
	Comment json_types.JSONNullString `json:"comment"   exposure:"public" swaggertype:"string"`
	// When the comment was created (public)
	CreatedAt time.Time `json:"createdAt" exposure:"public"`
	// When the comment was last updated (public)
	UpdatedAt time.Time `json:"updatedAt" exposure:"public"`
	// ID of the product being commented on (public)
	ProductId int `json:"productId" exposure:"public"`
	// ID of the user who made the comment (public)
	UserId int `json:"userId"    exposure:"public"`
}

// ProductCommentWithUser combines a product comment with user information
// @model ProductCommentWithUser
type ProductCommentWithUser struct {
	// Unique comment identifier (public)
	Id int `json:"id"        exposure:"public"`
	// Rating score (1-5) (public)
	Scoring int `json:"scoring"   exposure:"public"`
	// Text of the comment (public)
	Comment json_types.JSONNullString `json:"comment"   exposure:"public" swaggertype:"string"`
	// When the comment was created (public)
	CreatedAt time.Time `json:"createdAt" exposure:"public"`
	// When the comment was last updated (public)
	UpdatedAt time.Time `json:"updatedAt" exposure:"public"`
	// ID of the product being commented on (public)
	ProductId int `json:"productId" exposure:"public"`
	// User who made the comment (public)
	User CommentUser `json:"user"      exposure:"public"`
}

// Product combines basic product information with additional details
// @model Product
type Product struct {
	ProductBase
	// Product leaf subcategory info
	Subcategory ProductCategory `json:"subcategory" exposure:"public"`
	// Average score of the product based on the comments
	AverageScore float32 `json:"averageScore" exposure:"public"`
	// Total available quantity across all variants (public)
	TotalQuantity int `json:"totalQuantity"       exposure:"public"`
	// Current offer/discount, if any (public, optional)
	Offer *ProductOffer `json:"offer,omitempty"     exposure:"public"`
	// Main product image (public, optional)
	MainImage *ProductImage `json:"mainImage,omitempty" exposure:"public"`
	// Store information (public)
	Store StoreInfo `json:"store"               exposure:"public"`
}

// ProductExtended provides comprehensive product information
// @model ProductExtended
type ProductExtended struct {
	ProductBase
	// Subcategory information with parent hierarchy (public)
	Subcategory ProductCategoryWithParents `json:"subcategory"     exposure:"public"`
	// List of product specifications (public)
	Specs []ProductSpec `json:"specs"           exposure:"public"`
	// List of tags assigned to the product (public)
	Tags []ProductTag `json:"tags"            exposure:"public"`
	// List of available variants with attributes (public)
	Variants []ProductVariantWithAttributeSet `json:"variants"        exposure:"public"`
	// List of available attributes with their available options based on the product variants (public)
	Attributes []ProductAttributeWithOptions `json:"attributes"  exposure:"public"`
	// Current offer/discount, if any (public, optional)
	Offer *ProductOffer `json:"offer,omitempty" exposure:"public"`
	// List of product images (public)
	Images []ProductImage `json:"images"          exposure:"public"`
	// Store information (public)
	Store StoreInfo `json:"store"           exposure:"public"`
}

// CreateProductTagPayload contains data needed to create a new product tag
// @model CreateProductTagPayload
type CreateProductTagPayload struct {
	// Name of the tag (required)
	Name string `json:"name" validate:"required"`
}

// UpdateProductTagPayload contains data for updating a product tag
// @model UpdateProductTagPayload
type UpdateProductTagPayload struct {
	// New name for the tag
	Name *string `json:"name"`
}

// ProductTagSearchQuery contains parameters for searching product tags
// @model ProductTagSearchQuery
type ProductTagSearchQuery struct {
	// Filter by tag name
	Name *string `json:"name"`
	// Filter by product ID
	ProductId *int `json:"productId"`
	// Maximum number of results
	Limit *int `json:"limit"`
	// Number of results to skip
	Offset *int `json:"offset"`
}

// CreateProductCategoryPayload contains data needed to create a new product category
// @model CreateProductCategoryPayload
type CreateProductCategoryPayload struct {
	// Name of the category (required)
	Name string `json:"name"             validate:"required"`
	// Image filename (required)
	ImageName string `json:"imageName"        validate:"required"`
	// ID of the parent category, if any
	ParentCategoryId *int `json:"parentCategoryId"`
}

// UpdateProductCategoryPayload contains data for updating a product category
// @model UpdateProductCategoryPayload
type UpdateProductCategoryPayload struct {
	// New name for the category
	Name *string `json:"name"`
	// New image filename
	ImageName *string `json:"imageName"`
}

// ProductCategorySearchQuery contains parameters for searching product categories
// @model ProductCategorySearchQuery
type ProductCategorySearchQuery struct {
	// Filter by category name
	Name *string `json:"name"`
	// Filter by parent category ID
	ParentCategoryId *int `json:"parentCategoryId"`
	// Maximum number of results
	Limit *int `json:"limit"`
	// Number of results to skip
	Offset *int `json:"offset"`
}

// CreateProductBasePayload contains core data needed to create a new product
// @model CreateProductBasePayload
type CreateProductBasePayload struct {
	// Product name (required)
	Name string `json:"name"           validate:"required"`
	// URL-friendly slug
	Slug string `json:"slug"`
	// Product price (required)
	Price float64 `json:"price"          validate:"required"`
	// Shipping cost factor (required)
	ShipmentFactor float64 `json:"shipmentFactor" validate:"required"`
	// Product description
	Description string `json:"description"`
	// Subcategory ID (required)
	SubcategoryId int `json:"subcategoryId"  validate:"required"`
	// Store ID (required)
	StoreId int `json:"storeId"        validate:"required"`
}

// CreateProductImagePayload contains data needed to add a product image
// @model CreateProductImagePayload
type CreateProductImagePayload struct {
	// Image filename (required)
	ImageName string `json:"imageName" validate:"required"`
	// Whether this is the main product image (required)
	IsMain bool `json:"isMain"    validate:"required"`
}

// CreateProductSpecPayload contains data needed to add a product specification
// @model CreateProductSpecPayload
type CreateProductSpecPayload struct {
	// Specification label (required)
	Label string `json:"label" validate:"required"`
	// Specification value (required)
	Value string `json:"value" validate:"required"`
}

// ProductVariantAttributeSetPayload contains attribute data for a product variant
// @model ProductVariantAttributeSetPayload
type ProductVariantAttributeSetPayload struct {
	// Attribute ID (required)
	AttributeId int `json:"attributeId" validate:"required"`
	// Selected option ID (required)
	OptionId int `json:"optionId"    validate:"required"`
}

// CreateProductVariantPayload contains data needed to create a product variant
// @model CreateProductVariantPayload
type CreateProductVariantPayload struct {
	// Initial stock quantity (required)
	Quantity int `json:"quantity"      validate:"required"`
	// Set of attributes defining this variant (required)
	AttributeSets []ProductVariantAttributeSetPayload `json:"attributeSets" validate:"required"`
}

// CreateProductPayload contains complete data needed to create a new product
// @model CreateProductPayload
type CreateProductPayload struct {
	// Base product information (required)
	Base CreateProductBasePayload `json:"base"     validate:"required"`
	// List of tag IDs to assign (required)
	TagIds []int `json:"tagIds"   validate:"required"`
	// List of product images (required)
	Images []CreateProductImagePayload `json:"images"   validate:"required"`
	// List of product specifications (required)
	Specs []CreateProductSpecPayload `json:"specs"    validate:"required"`
	// List of product variants (required)
	Variants []CreateProductVariantPayload `json:"variants" validate:"required"`
}

// UpdateProductBasePayload contains data for updating core product information
// @model UpdateProductBasePayload
type UpdateProductBasePayload struct {
	// New product name
	Name *string `json:"name"`
	// New URL-friendly slug
	Slug *string `json:"slug"`
	// New product price
	Price *float64 `json:"price"`
	// New shipping cost factor
	ShipmentFactor *float64 `json:"shipmentFactor"`
	// New product description
	Description *string `json:"description"`
	// New subcategory ID
	SubcategoryId *int `json:"subcategoryId"`
	// New active status
	IsActive *bool `json:"isActive"`
}

// UpdateProductSpecPayload contains data for updating a product specification
// @model UpdateProductSpecPayload
type UpdateProductSpecPayload struct {
	// New specification label
	Label *string `json:"label"`
	// New specification value
	Value *string `json:"value"`
}

// UpdateProductVariantPayload contains data for updating a product variant
// @model UpdateProductVariantPayload
type UpdateProductVariantPayload struct {
	// New stock quantity
	Quantity *int `json:"quantity"`
	// New attribute sets to add
	NewAttributeSets []ProductVariantAttributeSetPayload `json:"newAttributeSets"`
	// Attribute IDs to remove
	DelAttributeIds []int `json:"delAttributeIds"`
}

// UpdatedProductSpecPayload combines spec ID with update data
// @model UpdatedProductSpecPayload
type UpdatedProductSpecPayload struct {
	// ID of the spec to update
	Id int
	UpdateProductSpecPayload
}

// UpdatedProductVariantPayload combines variant ID with update data
// @model UpdatedProductVariantPayload
type UpdatedProductVariantPayload struct {
	// ID of the variant to update
	Id int
	UpdateProductVariantPayload
}

// UpdateProductPayload contains complete data for updating a product
// @model UpdateProductPayload
type UpdateProductPayload struct {
	// Base product updates
	Base *UpdateProductBasePayload `json:"base"`
	// New tag IDs to add
	NewTagIds []int `json:"newTagIds"`
	// Tag IDs to remove
	DelTagIds []int `json:"delTagIds"`
	// New images to add
	NewImages []CreateProductImagePayload `json:"newImages"`
	// ID of the new main image
	NewMainImage *int `json:"newMainImage"`
	// Image IDs to remove
	DelImageIds []int `json:"delImageIds"`
	// New specs to add
	NewSpecs []CreateProductSpecPayload `json:"newSpecs"`
	// Existing specs to update
	UpdatedSpecs []UpdatedProductSpecPayload `json:"updatedSpecs"`
	// Spec IDs to remove
	DelSpecIds []int `json:"delSpecIds"`
	// New variants to add
	NewVariants []CreateProductVariantPayload `json:"newVariants"`
	// Existing variants to update
	UpdatedVariants []UpdatedProductVariantPayload `json:"updatedVariants"`
	// Variant IDs to remove
	DelVariantIds []int `json:"delVariantIds"`
}

// ProductSearchQuery contains parameters for searching products
// @model ProductSearchQuery
type ProductSearchQuery struct {
	// General search keyword
	Keyword *string `json:"keyword"`
	// Filter by exact product name
	Name *string `json:"name"`
	// Filter by exact slug
	Slug *string `json:"slug"`
	// Minimum available quantity
	MinQuantity *int `json:"minQuantity"`
	// Maximum available quantity
	MaxQuantity *int `json:"maxQuantity"`
	// Whether product has an active offer
	HasOffer *bool `json:"hasOffer"`
	// Filter by category ID
	CategoryId *int `json:"categoryId"`
	// Filter by tag name
	TagName *string `json:"tagName"`
	// Filter by tag IDs (separated by comma ',')
	TagIds *string `json:"tagIds"`
	// Maximum price
	PriceLessThan *float32 `json:"priceLessThan"`
	// Minimum price
	PriceMoreThan *float32 `json:"priceMoreThan"`
	// Filter by store ID
	StoreId *int `json:"storeId"`
	// Minimum average score
	AverageScore *float32 `json:"averageScore"`
	// Filter by active status
	IsActive *bool `json:"isActive"`
	// Maximum number of results
	Limit *int `json:"limit"`
	// Number of results to skip
	Offset *int `json:"offset"`
}

// CreateProductOfferPayload contains data needed to create a product offer
// @model CreateProductOfferPayload
type CreateProductOfferPayload struct {
	// Discount percentage (required)
	Discount float64 `json:"discount"  validate:"required"`
	// Expiration date (required)
	ExpireAt time.Time `json:"expireAt"  validate:"required"`
	// Product ID this offer applies to
	ProductId int `json:"productId"`
}

// UpdateProductOfferPayload contains data for updating a product offer
// @model UpdateProductOfferPayload
type UpdateProductOfferPayload struct {
	// New discount percentage
	Discount *float64 `json:"discount"`
	// New expiration date
	ExpireAt *time.Time `json:"expireAt"`
}

// ProductOfferSearchQuery contains parameters for searching product offers
// @model ProductOfferSearchQuery
type ProductOfferSearchQuery struct {
	// Maximum discount percentage
	DiscountLessThan *float64 `json:"discountLessThan"`
	// Minimum discount percentage
	DiscountMoreThan *float64 `json:"discountMoreThan"`
	// Offers expiring before this date
	ExpireAtLessThan *time.Time `json:"expireAtLessThan"`
	// Offers expiring after this date
	ExpireAtMoreThan *time.Time `json:"expireAtMoreThan"`
	// Maximum number of results
	Limit *int `json:"limit"`
	// Number of results to skip
	Offset *int `json:"offset"`
}

// CreateProductAttributePayload contains data needed to create a product attribute
// @model CreateProductAttributePayload
type CreateProductAttributePayload struct {
	// Attribute name/label (required)
	Label string `json:"label"   validate:"required"`
	// List of option values (required)
	Options []string `json:"options" validate:"required"`
}

// UpdateProductAttributeOptionPayload contains data for updating an attribute option
// @model UpdateProductAttributeOptionPayload
type UpdateProductAttributeOptionPayload struct {
	// New option value
	Value *string `json:"value"`
}

// UpdatedProductAttributeOptionPayload combines option ID with update data
// @model UpdatedProductAttributeOptionPayload
type UpdatedProductAttributeOptionPayload struct {
	// ID of the option to update
	Id int
	UpdateProductAttributeOptionPayload
}

// UpdateProductAttributePayload contains data for updating a product attribute
// @model UpdateProductAttributePayload
type UpdateProductAttributePayload struct {
	// New attribute label
	Label *string `json:"label"`
	// New options to add
	NewOptions []string `json:"newOptions"`
	// Existing options to update
	UpdatedOptions []UpdatedProductAttributeOptionPayload `json:"updatedOptions"`
	// Option IDs to remove
	DelOptionIds []int `json:"delOptionIds"`
}

// ProductAttributeSearchQuery contains parameters for searching product attributes
// @model ProductAttributeSearchQuery
type ProductAttributeSearchQuery struct {
	// Filter by attribute label
	Label *string `json:"label"`
	// Maximum number of results
	Limit *int `json:"limit"`
	// Number of results to skip
	Offset *int `json:"offset"`
}

// CreateProductCommentPayload contains data needed to create a product comment
// @model CreateProductCommentPayload
type CreateProductCommentPayload struct {
	// Rating score (1-5) (required)
	Scoring int `json:"scoring"   validate:"required"`
	// Comment text (required)
	Comment string `json:"comment"   validate:"required"`
	// Product ID being commented on
	ProductId int `json:"productId"`
	// User ID making the comment
	UserId int `json:"userId"`
}

// UpdateProductCommentPayload contains data for updating a product comment
// @model UpdateProductCommentPayload
type UpdateProductCommentPayload struct {
	// New rating score
	Scoring *int `json:"scoring"`
	// Updated comment text
	Comment *string `json:"comment"`
}

// ProductCommentSearchQuery contains parameters for searching product comments
// @model ProductCommentSearchQuery
type ProductCommentSearchQuery struct {
	// Maximum rating score
	ScoringLessThan *int `json:"scoringLessThan"`
	// Minimum rating score
	ScoringMoreThan *int `json:"scoringMoreThan"`
	// Maximum number of results
	Limit *int `json:"limit"`
	// Number of results to skip
	Offset *int `json:"offset"`
}
