package handler

import (
	"net/http"

	"github.com/TNJKL/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

const passwordLength = 12

type GenPass interface {
	GeneratePassword(c *gin.Context)
}

type genPassHandler struct {
	genPassService service.GenPass
}

func NewGenPass(genPassSvc service.GenPass) GenPass {
	return &genPassHandler{
		genPassService: genPassSvc,
	}
}

func (s *genPassHandler) GeneratePassword(c *gin.Context) {
	pass, err := s.genPassService.GeneratePassword(passwordLength)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
	c.JSON(http.StatusOK, gin.H{"password": pass})
}
