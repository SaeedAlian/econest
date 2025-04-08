package vendor_store

import "time"

type Vendor struct {
	Id          int
	Name        string
	Description string
	Verified    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	OwnerId     int
}

type VendorPhoneNumber struct {
	Id          int
	CountryCode string
	Number      string
	Verified    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	VendorId    int
}

type VendorAddress struct {
	Id        int
	State     string
	City      string
	Street    string
	Zipcode   string
	Details   string
	CreatedAt time.Time
	UpdatedAt time.Time
	VendorId  int
}

type VendorProduct struct {
	VendorId  int
	ProductId int
}
