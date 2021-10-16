package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	UserName      string    `db:"username" json:"username"`
	PasswordHash  []byte    `db:"phash" json:"phash"`
	Name          string    `db:"name" json:"name"`
	SurName       string    `db:"surname" json:"surname"`
	BirthDate     time.Time `db:"birthdate" json:"birthdate"`
	Gender        bool      `db:"gender" json:"gender"`
	City          string    `db:"city" json:"city"`
	Interests     []string  `json:"interests"`
	Friends       []string  `json:"friends"`
	Subscriptions []string  `json:"subscriptions"`
	Subscribers   []string  `json:"subscribers"`
	Deleted       bool      `db:"deleted" json:"deleted"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.UserName, validation.Required, validation.Length(2, 500).Error("length must be 2 to 500")),
		validation.Field(&u.PasswordHash, validation.Required),
	)
}

type UserFilter struct {
	UserName      string    `json:"username"`
	Name          string    `json:"name"`
	SurName       string    `json:"surname"`
	BirthDateFrom time.Time `json:"bd_from"`
	BirthDateTo   time.Time `json:"bd_to"`
	AgeFrom       string    `json:"age_from"`
	AgeTo         string    `json:"age_to"`
	Gender        string    `json:"gender"`
	City          string    `json:"city"`
	Interests     string    `json:"interests"`
	Friends       string    `json:"friends"`
	PageNum       int       `json:"pagenum"`
	PageSize      int       `json:"pagesize"`
}

func (u *UserFilter) Validate() error {
	return validation.ValidateStruct(
		u,
	)
}
