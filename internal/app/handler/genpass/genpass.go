package genpass

import (
	"net/http"

	"github.com/TNJKL/bookmark-management/internal/app/service/genpass"
	"github.com/gin-gonic/gin"
)

const passwordLength = 12

// Interface: Handler có thể làm gì
type GenPass interface {
	GeneratePassword(c *gin.Context)
}

// Struct: Handler giữ reference tới Service
type genPassHandler struct {
	genPassService genpass.GenPass // ← Đây là "ổ cắm" để nhận Service
}

// Constructor: Nhận service từ bên ngoài truyền vào
func NewGenPass(genPassSvc genpass.GenPass) GenPass {
	return &genPassHandler{
		genPassService: genPassSvc, // ← "Cắm" service vào handler
	}
}

// GeneratePassword godoc
// @Summary      Generate a random password
// @Tags         password
// @Produce      json
// @Success      200  {object}  model.HealthCheckResponse
// @Router       /genpass [get]
func (s *genPassHandler) GeneratePassword(c *gin.Context) {
	// Gọi xuống tầng Service để xử lý logic
	pass, err := s.genPassService.GeneratePassword(passwordLength)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Trả kết quả về client
	c.JSON(http.StatusOK, gin.H{"password": pass})
}
