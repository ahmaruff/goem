package health

import (
	"context"
	"database/sql"
	"fmt"

	"maka-go/internal/logger"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Check(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		// This is where the error is born, so this is where the stack is
		// captured. The HTTP layer logs it once and turns it into a response.
		return logger.WithStack(fmt.Errorf("ping database: %w", err))
	}

	return nil
}
