package model

// User represents the user account data model in the database
type User struct {
	Base
	DisplayName string `json:"display_name"`
	Username    string `json:"username" gorm:"unique"`
	Password    string `json:"-"`
	Email       string `json:"email" gorm:"unique"`
}
