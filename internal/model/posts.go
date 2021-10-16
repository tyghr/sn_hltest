package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Post struct {
	UserName string    `db:"username" json:"username"`
	Header   string    `db:"header" json:"header"`
	Text     string    `db:"text" json:"text"`
	Created  time.Time `db:"created" json:"created"`
	Updated  time.Time `db:"updated" json:"updated"`
}

func (p *Post) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.Header, validation.Required, validation.Length(10, 500).Error("length must be 10 to 500")),
	)
}

type PostFilter struct {
	UserName    string `json:"username"`
	Header      string `json:"header"`
	Text        string `json:"text"`
	UpdatedFrom string `json:"updated_from"`
	UpdatedTo   string `json:"updated_to"`
	PageNum     int    `json:"pagenum"`
	PageSize    int    `json:"pagesize"`
}

func (p *PostFilter) Validate() error {
	return validation.ValidateStruct(
		p,
	)
}

type PostBacket struct {
	Post        Post     `json:"post"`
	Subscribers []string `json:"subscribers"`
}
