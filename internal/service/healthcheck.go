package service

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/model"
	"github.com/TNJKL/bookmark-management/internal/repository/ping"
)

//healthcheck interface

// HealthChecker defines the service operations for checking system health.
//
//go:generate mockery --name HealthChecker --filename healthcheck.go
type HealthChecker interface {
	// HealthCheck runs the health checks and returns the system status.
	HealthCheck(ctx context.Context) (*model.HealthCheckResponse, error)
}

// healthcheck struct
// healthCheckService is the default implementation of the HealthChecker interface.
type healthCheckService struct {
	serviceName string
	instanceID  string
	healthRepo  ping.HealthRepository
}

// constructor khởi tạo
// NewHealthCheck creates a new HealthChecker service instance.
func NewHealthCheck(serviceName string, instanceID string, healthRepo ping.HealthRepository) HealthChecker {
	return &healthCheckService{
		serviceName: serviceName,
		instanceID:  instanceID,
		healthRepo:  healthRepo,
	}
}

// đây là method của struct healthCheckService , cách nhận biết method của 1 struct là nhìn vào receiver của method đó đúng không anh ?
// ví dụ method này có receiver là (h *healthCheckService) nên nó là method của "healthCheckService struct"
// HealthCheck runs the health checks (e.g. database ping) and returns the system status.
func (h *healthCheckService) HealthCheck(ctx context.Context) (*model.HealthCheckResponse, error) {
	err := h.healthRepo.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &model.HealthCheckResponse{
		Message:     "OK",
		ServiceName: h.serviceName,
		InstanceID:  h.instanceID,
	}, nil
}
