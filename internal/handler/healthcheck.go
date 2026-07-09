package handler

import (
	"net/http"

	"github.com/TNJKL/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

// HealthCheck defines the handler contract for the health check endpoint.
type HealthCheck interface {
	HealthCheck(ctx *gin.Context)
}

// healthCheckHandler implements the HealthCheck interface.
type healthCheckHandler struct {
	healthCheckService service.HealthChecker
}

// healthcheck constructor
// NewHealthCheck creates a new HealthCheck handler instance.
func NewHealthCheck(heathCheckSvc service.HealthChecker) HealthCheck {
	return &healthCheckHandler{
		healthCheckService: heathCheckSvc,
	}
}

// Các phần liên quan tới Swagger  và MakeFile em nhờ AI gen thử còn  code bài tập là em dựa vào bài học rồi tự viết lại ạ
// HealthCheck godoc
// @Summary      System health-check
// @Description  Trả về trạng thái hoạt động hiện tại, tên dịch vụ và instance ID
// @Tags         health-check
// @Produce      json
// @Success      200  {object}  model.HealthCheckResponse
// @Failure      500  {object}  map[string]string "Lỗi hệ thống nội bộ"
// @Router       /health-check [get]
func (h *healthCheckHandler) HealthCheck(ctx *gin.Context) {
	// Gọi xuống tầng Service để xử lý logic
	result, err := h.healthCheckService.HealthCheck(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	// Trả kết quả về client
	ctx.JSON(http.StatusOK, result)
}
