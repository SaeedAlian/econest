package store_store

import "time"

type Store struct {
	Id          int
	Name        string
	Description string
	Verified    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	OwnerId     int
}

type StorePhoneNumber struct {
	Id          int
	CountryCode string
	Number      string
	IsPublic    bool
	Verified    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	StoreId     int
}

type StoreAddress struct {
	Id        int
	State     string
	City      string
	Street    string
	Zipcode   string
	Details   string
	IsPublic  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	StoreId   int
}

type StoreProduct struct {
	StoreId   int
	ProductId int
}
