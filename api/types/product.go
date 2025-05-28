package types

import (
	"time"

	json_types "github.com/SaeedAlian/econest/api/types/json"
)

type ProductBase struct {
	Id            int       `json:"id"            exposure:"public"`
	Name          string    `json:"name"          exposure:"public"`
	Slug          string    `json:"slug"          exposure:"public"`
	Price         float64   `json:"price"         exposure:"public"`
	Description   string    `json:"description"   exposure:"public"`
	IsActive      bool      `json:"isActive"      exposure:"public"`
	CreatedAt     time.Time `json:"createdAt"     exposure:"public"`
	UpdatedAt     time.Time `json:"updatedAt"     exposure:"public"`
	SubcategoryId int       `json:"subcategoryId" exposure:"public"`
}

type ProductCategory struct {
	Id               int                      `json:"id"               exposure:"public"`
	Name             string                   `json:"name"             exposure:"public"`
	ImageName        string                   `json:"imageName"        exposure:"public"`
	CreatedAt        time.Time                `json:"createdAt"        exposure:"public"`
	UpdatedAt        time.Time                `json:"updatedAt"        exposure:"public"`
	ParentCategoryId json_types.JSONNullInt32 `json:"parentCategoryId" exposure:"public"`
}

type ProductCategoryWithParents struct {
	ProductCategory
	ParentCategory *ProductCategoryWithParents `json:"parentCategory,omitempty" exposure:"public"`
}

type ProductOffer struct {
	Id        int       `json:"id"        exposure:"public"`
	Discount  float64   `json:"discount"  exposure:"public"`
	ExpireAt  time.Time `json:"expireAt"  exposure:"public"`
	CreatedAt time.Time `json:"createdAt" exposure:"public"`
	UpdatedAt time.Time `json:"updatedAt" exposure:"public"`
	ProductId int       `json:"productId" exposure:"public"`
}

type ProductImage struct {
	Id        int    `json:"id"        exposure:"public"`
	ImageName string `json:"imageName" exposure:"public"`
	IsMain    bool   `json:"isMain"    exposure:"public"`
	ProductId int    `json:"productId" exposure:"public"`
}

type ProductSpec struct {
	Id        int    `json:"id"        exposure:"public"`
	Label     string `json:"label"     exposure:"public"`
	Value     string `json:"value"     exposure:"public"`
	ProductId int    `json:"productId" exposure:"public"`
}

type ProductTag struct {
	Id        int       `json:"id"        exposure:"public"`
	Name      string    `json:"name"      exposure:"public"`
	CreatedAt time.Time `json:"createdAt" exposure:"public"`
	UpdatedAt time.Time `json:"updatedAt" exposure:"public"`
}

type ProductTagAssignment struct {
	ProductId int `json:"productId" exposure:"public"`
	TagId     int `json:"tagId"     exposure:"public"`
}

type ProductAttribute struct {
	Id    int    `json:"id"    exposure:"public"`
	Label string `json:"label" exposure:"public"`
}

type ProductAttributeOption struct {
	Id          int    `json:"id"          exposure:"public"`
	Value       string `json:"value"       exposure:"public"`
	AttributeId int    `json:"attributeId" exposure:"public"`
}

type ProductAttributeWithOptions struct {
	ProductAttribute
	Options []ProductAttributeOption `json:"options" exposure:"public"`
}

type ProductVariant struct {
	Id        int `json:"id"        exposure:"public"`
	Quantity  int `json:"quantity"  exposure:"public"`
	ProductId int `json:"productId" exposure:"public"`
}

type ProductVariantAttributeOption struct {
	VariantId   int `json:"variantId"   exposure:"public"`
	AttributeId int `json:"attributeId" exposure:"public"`
	OptionId    int `json:"optionId"    exposure:"public"`
}

type ProductVariantSelectedAttributeOption struct {
	ProductAttribute
	SelectedOption ProductAttributeOption `json:"selectedOption" exposure:"public"`
}

type ProductVariantWithAttributeSet struct {
	ProductVariant
	AttributeSet []ProductVariantSelectedAttributeOption `json:"attributeSet" exposure:"public"`
}

type ProductComment struct {
	Id        int                       `json:"id"        exposure:"public"`
	Scoring   int                       `json:"scoring"   exposure:"public"`
	Comment   json_types.JSONNullString `json:"comment"   exposure:"public"`
	CreatedAt time.Time                 `json:"createdAt" exposure:"public"`
	UpdatedAt time.Time                 `json:"updatedAt" exposure:"public"`
	ProductId int                       `json:"productId" exposure:"public"`
	UserId    int                       `json:"userId"    exposure:"public"`
}

type Product struct {
	ProductBase
	TotalQuantity int           `json:"totalQuantity"       exposure:"public"`
	Offer         *ProductOffer `json:"offer,omitempty"     exposure:"public"`
	MainImage     *ProductImage `json:"mainImage,omitempty" exposure:"public"`
	Store         StoreInfo     `json:"store"               exposure:"public"`
}

type ProductExtended struct {
	ProductBase
	Subcategory ProductCategoryWithParents       `json:"subcategory"     exposure:"public"`
	Specs       []ProductSpec                    `json:"specs"           exposure:"public"`
	Tags        []ProductTag                     `json:"tags"            exposure:"public"`
	Variants    []ProductVariantWithAttributeSet `json:"variants"        exposure:"public"`
	Offer       *ProductOffer                    `json:"offer,omitempty" exposure:"public"`
	Images      []ProductImage                   `json:"images"          exposure:"public"`
	Store       StoreInfo                        `json:"store"           exposure:"public"`
}

type CreateProductTagPayload struct {
	Name string `json:"name" validate:"required"`
}

type UpdateProductTagPayload struct {
	Name *string `json:"name"`
}

type ProductTagSearchQuery struct {
	Name      *string `json:"name"`
	ProductId *int    `json:"productId"`
	Limit     *int    `json:"limit"`
	Offset    *int    `json:"offset"`
}

type CreateProductCategoryPayload struct {
	Name             string `json:"name"             validate:"required"`
	ImageName        string `json:"imageName"        validate:"required"`
	ParentCategoryId *int   `json:"parentCategoryId"`
}

type UpdateProductCategoryPayload struct {
	Name *string `json:"name"`
}

type ProductCategorySearchQuery struct {
	Name             *string `json:"name"`
	ParentCategoryId *int    `json:"parentCategoryId"`
	Limit            *int    `json:"limit"`
	Offset           *int    `json:"offset"`
}

type CreateProductBasePayload struct {
	Name          string  `json:"name"          validate:"required"`
	Slug          string  `json:"slug"`
	Price         float64 `json:"price"         validate:"required"`
	Description   string  `json:"description"`
	SubcategoryId int     `json:"subcategoryId" validate:"required"`
	StoreId       int     `json:"storeId"       validate:"required"`
}

type CreateProductImagePayload struct {
	ImageName string `json:"imageName" validate:"required"`
	IsMain    bool   `json:"isMain"    validate:"required"`
}

type CreateProductSpecPayload struct {
	Label string `json:"label" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type ProductVariantAttributeSetPayload struct {
	AttributeId int `json:"attributeId" validate:"required"`
	OptionId    int `json:"optionId"    validate:"required"`
}

type CreateProductVariantPayload struct {
	Quantity      int                                 `json:"quantity"      validate:"required"`
	AttributeSets []ProductVariantAttributeSetPayload `json:"attributeSets" validate:"required"`
}

type CreateProductPayload struct {
	Base     CreateProductBasePayload      `json:"base"     validate:"required"`
	TagIds   []int                         `json:"tagIds"   validate:"required"`
	Images   []CreateProductImagePayload   `json:"images"   validate:"required"`
	Specs    []CreateProductSpecPayload    `json:"specs"    validate:"required"`
	Variants []CreateProductVariantPayload `json:"variants" validate:"required"`
}

type UpdateProductBasePayload struct {
	Name          *string  `json:"name"`
	Slug          *string  `json:"slug"`
	Price         *float64 `json:"price"`
	Description   *string  `json:"description"`
	SubcategoryId *int     `json:"subcategoryId"`
	IsActive      *bool    `json:"isActive"`
}

type UpdateProductSpecPayload struct {
	Label *string `json:"label"`
	Value *string `json:"value"`
}

type UpdateProductVariantPayload struct {
	Quantity         *int                                `json:"quantity"`
	NewAttributeSets []ProductVariantAttributeSetPayload `json:"newAttributeSets"`
	DelAttributeIds  []int                               `json:"delAttributeIds"`
}

type UpdatedProductSpecPayload struct {
	Id int
	UpdateProductSpecPayload
}

type UpdatedProductVariantPayload struct {
	Id int
	UpdateProductVariantPayload
}

type UpdateProductPayload struct {
	Base            *UpdateProductBasePayload      `json:"base"`
	NewTagIds       []int                          `json:"newTagIds"`
	DelTagIds       []int                          `json:"delTagIds"`
	NewImages       []CreateProductImagePayload    `json:"newImages"`
	NewMainImage    *int                           `json:"newMainImage"`
	DelImageIds     []int                          `json:"delImageIds"`
	NewSpecs        []CreateProductSpecPayload     `json:"newSpecs"`
	UpdatedSpecs    []UpdatedProductSpecPayload    `json:"updatedSpecs"`
	DelSpecIds      []int                          `json:"delSpecIds"`
	NewVariants     []CreateProductVariantPayload  `json:"newVariants"`
	UpdatedVariants []UpdatedProductVariantPayload `json:"updatedVariants"`
	DelVariantIds   []int                          `json:"delVariantIds"`
}

type ProductSearchQuery struct {
	Keyword       *string `json:"keyword"`
	Name          *string `json:"name"`
	Slug          *string `json:"slug"`
	MinQuantity   *int    `json:"minQuantity"`
	HasOffer      *bool   `json:"hasOffer"`
	CategoryId    *int    `json:"categoryId"`
	TagName       *string `json:"tagName"`
	TagId         *int    `json:"tagId"`
	PriceLessThan *int    `json:"priceLessThan"`
	PriceMoreThan *int    `json:"priceMoreThan"`
	StoreId       *int    `json:"storeId"`
	IsActive      *bool   `json:"isActive"`
	Limit         *int    `json:"limit"`
	Offset        *int    `json:"offset"`
}

type CreateProductOfferPayload struct {
	Discount  float64   `json:"discount"  validate:"required"`
	ExpireAt  time.Time `json:"expireAt"  validate:"required"`
	ProductId int       `json:"productId" validate:"required"`
}

type UpdateProductOfferPayload struct {
	Discount *float64   `json:"discount"`
	ExpireAt *time.Time `json:"expireAt"`
}

type ProductOfferSearchQuery struct {
	DiscountLessThan *float64   `json:"discountLessThan"`
	DiscountMoreThan *float64   `json:"discountMoreThan"`
	ExpireAtLessThan *time.Time `json:"expireAtLessThan"`
	ExpireAtMoreThan *time.Time `json:"expireAtMoreThan"`
	Limit            *int       `json:"limit"`
	Offset           *int       `json:"offset"`
}

type CreateProductAttributePayload struct {
	Label   string   `json:"label"   validate:"required"`
	Options []string `json:"options" validate:"required"`
}

type UpdateProductAttributeOptionPayload struct {
	Value *string `json:"value"`
}

type UpdatedProductAttributeOptionPayload struct {
	Id int
	UpdateProductAttributeOptionPayload
}

type UpdateProductAttributePayload struct {
	Label          *string                                `json:"label"`
	NewOptions     []string                               `json:"newOptions"`
	UpdatedOptions []UpdatedProductAttributeOptionPayload `json:"updatedOptions"`
	DelOptionIds   []int                                  `json:"delOptionIds"`
}

type ProductAttributeSearchQuery struct {
	Label  *string `json:"label"`
	Limit  *int    `json:"limit"`
	Offset *int    `json:"offset"`
}

type CreateProductCommentPayload struct {
	Scoring   int    `json:"scoring"   validate:"required"`
	Comment   string `json:"comment"   validate:"required"`
	ProductId int    `json:"productId" validate:"required"`
	UserId    int    `json:"userId"    validate:"required"`
}

type UpdateProductCommentPayload struct {
	Scoring *int    `json:"scoring"`
	Comment *string `json:"comment"`
}

type ProductCommentSearchQuery struct {
	ScoringLessThan *int `json:"scoringLessThan"`
	ScoringMoreThan *int `json:"scoringMoreThan"`
	Limit           *int `json:"limit"`
	Offset          *int `json:"offset"`
}
