package service

import "github.com/TNJKL/bookmark-management/internal/model"

//healthcheck interface

//go:generate mockery --name HealthChecker --filename healthcheck.go
type HealthChecker interface {
	HealthCheck() (*model.HealthCheckResponse, error)
}

// healthcheck struct
type healthCheckService struct {
	serviceName string
	instanceID  string
}

// constructor khởi tạo
func NewHealthCheck(serviceName string, instanceID string) HealthChecker {
	return &healthCheckService{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

// đây là method của struct healthCheckService , cách nhận biết method của 1 struct là nhìn vào receiver của method đó đúng không anh ?
// ví dụ method này có receiver là (h *healthCheckService) nên nó là method của "healthCheckService struct"
func (h *healthCheckService) HealthCheck() (*model.HealthCheckResponse, error) {
	return &model.HealthCheckResponse{
		Message:     "OK",
		ServiceName: h.serviceName,
		InstanceID:  h.instanceID,
	}, nil
}
