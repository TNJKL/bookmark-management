package model

// Bookmark represents a user's saved web bookmark entity in the database
type Bookmark struct {
	Base
	Description string `json:"description"`
	URL         string `json:"url"`
	Code        string `json:"code" gorm:"unique"`
	CodeInt     int    `json:"-" gorm:"column:code_int;autoIncrement"`
	UserID      string `json:"-"`
	User        *User  `gorm:"references:ID" json:"-"`
}
