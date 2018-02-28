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

func (s *metricsService) Create(ctx context.Context, userID, userUsername string) (*Session, error) {
	sess, err := s.service.Create(ctx, userID, userUsername)

	if err == nil {
		s.sessionsCreated.Add(1)
	}

	return sess, err
}

func (s *metricsService) Find(ctx context.Context, id string) (*Session, error) {
	return s.service.Find(ctx, id)
}

func (s *metricsService) Delete(ctx context.Context, id string) error {
	return s.service.Delete(ctx, id)
}

func (s *metricsService) DeleteExpired(ctx context.Context) (int64, error) {
	num, err := s.service.DeleteExpired(ctx)

	if err == nil {
		s.sessionsCleared.Add(float64(num))
	}

	return num, err
}
