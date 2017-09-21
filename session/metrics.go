package session

import (
	"context"

	"github.com/go-kit/kit/metrics"
)

type metricsService struct {
	service         Service
	sessionsCreated metrics.Counter
	sessionsCleared metrics.Counter
}

func NewMetricsService(service Service, sessionsCreated metrics.Counter, sessionsCleared metrics.Counter) Service {
	// Initialize counters with 0
	sessionsCreated.Add(0)
	sessionsCleared.Add(0)

	return &metricsService{
		service:         service,
		sessionsCreated: sessionsCreated,
		sessionsCleared: sessionsCleared,
	}
}

func (s *metricsService) CreateSession(ctx context.Context, userID, userUsername string) (*Session, error) {
	sess, err := s.service.CreateSession(ctx, userID, userUsername)

	if err == nil {
		s.sessionsCreated.Add(1)
	}

	return sess, err
}

func (s *metricsService) FindSession(ctx context.Context, id string) (*Session, error) {
	return s.service.FindSession(ctx, id)
}

func (s *metricsService) ClearSessions(ctx context.Context) (int64, error) {
	num, err := s.service.ClearSessions(ctx)

	if err == nil {
		s.sessionsCleared.Add(float64(num))
	}

	return num, err
}
