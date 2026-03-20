package service

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) Health() map[string]string {
	return map[string]string{
		"status": "ok",
	}
}

func (s *HealthService) Ready() map[string]string {
	return map[string]string{
		"status": "ready",
	}
}
