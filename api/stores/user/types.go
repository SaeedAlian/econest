package user_store

import "time"

type User struct {
	Id            int
	Username      string
	Email         string
	EmailVerified bool
	Password      string
	FullName      string
	BirthDate     time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	RoleId        int
}

type UserPhoneNumber struct {
	Id          int
	CountryCode string
	Number      string
	Verified    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserId      int
}

type UserAddress struct {
	Id        int
	State     string
	City      string
	Street    string
	Zipcode   string
	Details   string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    int
}
