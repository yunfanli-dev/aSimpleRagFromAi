package service

type HealthService struct{}

// NewHealthService builds the lightweight liveness/readiness service.
func NewHealthService() *HealthService {
	return &HealthService{}
}

// Health reports a basic liveness signal.
func (s *HealthService) Health() map[string]string {
	return map[string]string{
		"status": "ok",
	}
}

// Ready reports a basic readiness signal.
func (s *HealthService) Ready() map[string]string {
	return map[string]string{
		"status": "ready",
	}
}
