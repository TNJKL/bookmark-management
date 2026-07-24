package user

import (
	"unicode"

	"github.com/TNJKL/bookmark-management/internal/service/user"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Handler defines the contract for handling HTTP requests related to users
type Handler interface {
	Register(ctx *gin.Context)
}

// userHandler implements the Handler interface and communicates with the user Service
type userHandler struct {
	svc user.Service
}

// NewHandler creates a new user Handler instance and registers the custom password validation
func NewHandler(svc user.Service) Handler {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("strong_password", validateStrongPassword)
	}
	return &userHandler{svc: svc}
}

// validateStrongPassword checks if the input password meets the minimum strength requirements
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}

	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}
