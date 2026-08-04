package model

// Bookmark represents a user's saved web bookmark entity in the database
type Bookmark struct {
	Base
	Description string `json:"description"`
	URL         string `json:"url"`
	Code        string `json:"code" gorm:"unique"`
	UserID      string `json:"user_id"`
	User        *User  `gorm:"references:ID" json:"-"`
}
