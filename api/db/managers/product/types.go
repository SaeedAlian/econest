package product_db_manager

import "time"

type Product struct {
	Id            int
	Name          string
	Slug          string
	Price         float32
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	SubcategoryId int
}

type ProductCategory struct {
	Id               int
	Name             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ParentCategoryId int
}

type ProductImage struct {
	Id        int
	ImageName string
	IsMain    string
	CreatedAt time.Time
	UpdatedAt time.Time
	ProductId int
}

type ProductSpec struct {
	Id        int
	Label     string
	Value     string
	ProductId int
}

type ProductTag struct {
	Id        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductTagAssignment struct {
	ProductId int
	TagId     int
}

type ProductAttribute struct {
	Id        int
	Label     string
	ProductId int
}

type ProductAttributeOption struct {
	Id    int
	Value string

	AttributeId int
}

type ProductVariant struct {
	Id       int
	Quantity int

	ProductId int
}

type ProductVariantOption struct {
	Id          int
	VariantId   int
	AttributeId int
	OptionId    int
}

type ProductComments struct {
	Id        int
	Scoring   int
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
	ProductId int
	UserId    int
}
