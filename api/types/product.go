package types

import (
	"database/sql"
	"time"
)

type Product struct {
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
	Id               int           `json:"id"               exposure:"public"`
	Name             string        `json:"name"             exposure:"public"`
	ImageName        string        `json:"imageName"        exposure:"public"`
	CreatedAt        time.Time     `json:"createdAt"        exposure:"public"`
	UpdatedAt        time.Time     `json:"updatedAt"        exposure:"public"`
	ParentCategoryId sql.NullInt32 `json:"parentCategoryId" exposure:"public"`
}

type ProductCategoryWithParents struct {
	Id             int                         `json:"id"                       exposure:"public"`
	Name           string                      `json:"name"                     exposure:"public"`
	CreatedAt      time.Time                   `json:"createdAt"                exposure:"public"`
	UpdatedAt      time.Time                   `json:"updatedAt"                exposure:"public"`
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

type ProductSpecInfo struct {
	Id    int    `json:"id"    exposure:"public"`
	Label string `json:"label" exposure:"public"`
	Value string `json:"value" exposure:"public"`
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
	Id        int    `json:"id"        exposure:"public"`
	Label     string `json:"label"     exposure:"public"`
	ProductId int    `json:"productId" exposure:"public"`
}

type ProductAttributeOption struct {
	Id          int    `json:"id"          exposure:"public"`
	Value       string `json:"value"       exposure:"public"`
	AttributeId int    `json:"attributeId" exposure:"public"`
}

type ProductAttributeOptionInfo struct {
	Id    int    `json:"id"    exposure:"public"`
	Value string `json:"value" exposure:"public"`
}

type ProductAttributeWithOptions struct {
	Id      int                          `json:"id"      exposure:"public"`
	Label   string                       `json:"label"   exposure:"public"`
	Options []ProductAttributeOptionInfo `json:"options" exposure:"public"`
}

type ProductVariant struct {
	Id        int `json:"id"        exposure:"public"`
	Quantity  int `json:"quantity"  exposure:"public"`
	ProductId int `json:"productId" exposure:"public"`
}

type ProductVariantOption struct {
	VariantId   int `json:"variantId"   exposure:"public"`
	AttributeId int `json:"attributeId" exposure:"public"`
	OptionId    int `json:"optionId"    exposure:"public"`
}

type ProductVariantOptionInfo struct {
	AttributeId int `json:"attributeId" exposure:"public"`
	OptionId    int `json:"optionId"    exposure:"public"`
}

type ProductVariantInfo struct {
	Id       int                        `json:"id"       exposure:"public"`
	Quantity int                        `json:"quantity" exposure:"public"`
	Options  []ProductVariantOptionInfo `json:"options"  exposure:"public"`
}

type ProductComment struct {
	Id        int            `json:"id"        exposure:"public"`
	Scoring   int            `json:"scoring"   exposure:"public"`
	Comment   sql.NullString `json:"comment"   exposure:"public"`
	CreatedAt time.Time      `json:"createdAt" exposure:"public"`
	UpdatedAt time.Time      `json:"updatedAt" exposure:"public"`
	ProductId int            `json:"productId" exposure:"public"`
	UserId    int            `json:"userId"    exposure:"public"`
}

type ProductWithMainInfo struct {
	Product       `              json:"product"             exposure:"public"`
	TotalQuantity int           `json:"totalQuantity"       exposure:"public"`
	Offer         *ProductOffer `json:"offer,omitempty"     exposure:"public"`
	MainImage     *ProductImage `json:"mainImage,omitempty" exposure:"public"`
	Store         StoreInfo     `json:"store"               exposure:"public"`
}

type ProductWithAllInfo struct {
	Product     `                              json:"product"         exposure:"public"`
	Subcategory ProductCategoryWithParents    `json:"subcategory"     exposure:"public"`
	Specs       []ProductSpecInfo             `json:"specs"           exposure:"public"`
	Tags        []ProductTag                  `json:"tags"            exposure:"public"`
	Attributes  []ProductAttributeWithOptions `json:"attributes"      exposure:"public"`
	Variants    []ProductVariantInfo          `json:"variants"        exposure:"public"`
	Offer       *ProductOffer                 `json:"offer,omitempty" exposure:"public"`
	Images      []ProductImage                `json:"images"          exposure:"public"`
	Store       StoreInfo                     `json:"store"           exposure:"public"`
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

type CreateProductTagAssignment struct {
	ProductId int `json:"productId" validate:"required"`
	TagId     int `json:"tagId"     validate:"required"`
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

type CreateProductPayload struct {
	Name          string  `json:"name"          validate:"required"`
	Slug          string  `json:"slug"`
	Price         float64 `json:"price"         validate:"required"`
	Description   string  `json:"description"`
	SubcategoryId int     `json:"subcategoryId" validate:"required"`
	Quantity      int     `json:"quantity"      validate:"min=0,required"`
	StoreId       int     `json:"storeId"       validate:"required"`
}

type UpdateProductPayload struct {
	Name          *string  `json:"name"`
	Slug          *string  `json:"slug"`
	Price         *float64 `json:"price"`
	Description   *string  `json:"description"`
	SubcategoryId *int     `json:"subcategoryId"`
	IsActive      *bool    `json:"isActive"`
}

type ProductSearchQuery struct {
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

type CreateProductImagePayload struct {
	ImageName string `json:"imageName" validate:"required"`
	IsMain    bool   `json:"isMain"    validate:"required"`
	ProductId int    `json:"productId" validate:"required"`
}

type UpdateProductImagePayload struct {
	IsMain *bool `json:"isMain"`
}

type CreateProductSpecPayload struct {
	Label     string `json:"label"     validate:"required"`
	Value     string `json:"value"     validate:"required"`
	ProductId int    `json:"productId" validate:"required"`
}

type UpdateProductSpecPayload struct {
	Label *string `json:"label"`
	Value *string `json:"value"`
}

type CreateProductAttributePayload struct {
	Label     string `json:"label"     validate:"required"`
	ProductId int    `json:"productId" validate:"required"`
}

type UpdateProductAttributePayload struct {
	Label *string `json:"label"`
}

type CreateProductAttributeOptionPayload struct {
	Value       string `json:"value"       validate:"required"`
	AttributeId int    `json:"attributeId" validate:"required"`
}

type UpdateProductAttributeOptionPayload struct {
	Value *string `json:"value"`
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
