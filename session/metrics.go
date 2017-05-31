package session

import "github.com/go-kit/kit/metrics"

type metricsService struct {
	service         Service
	sessionsCreated metrics.Counter
	sessionsCleared metrics.Counter
}

func NewMetricsService(service Service, sessionsCreated metrics.Counter, sessionsCleared metrics.Counter) Service {
	return &metricsService{
		service:         service,
		sessionsCreated: sessionsCreated,
		sessionsCleared: sessionsCleared,
	}
}

func (s *metricsService) CreateSession(userID, userUsername string) (*Session, error) {
	sess, err := s.service.CreateSession(userID, userUsername)

	if err == nil {
		s.sessionsCreated.Add(1)
	}

	return sess, err
}

func (s *metricsService) FindSession(id string) (*Session, error) {
	return s.service.FindSession(id)
}

func (s *metricsService) ClearSessions() (int64, error) {
	num, err := s.service.ClearSessions()

	if err == nil {
		s.sessionsCleared.Add(float64(num))
	}

	return num, err
}
