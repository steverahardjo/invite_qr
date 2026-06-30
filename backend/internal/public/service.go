package public

import (
	"context"
	db_gen "invite_qr/db/db_gen"

	"github.com/google/uuid"
)

// Service wraps the sqlc-generated database queries with additional
// business logic for public participant operations.
type Service struct {
	DB *db_gen.Queries
}

// NewService creates a Service backed by the given sqlc queries instance.
func NewService(db *db_gen.Queries) *Service {
	return &Service{DB: db}
}

// GetParticipantByExternalID parses the externalID string as a UUID and
// returns the matching participant from the database, or an error.
func (s *Service) GetParticipantByExternalID(ctx context.Context, externalID string) (*db_gen.Participant, error) {
	id, err := uuid.Parse(externalID)
	if err != nil {
		return nil, err
	}
	participant, err := s.DB.GetParticipantByExternalID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &participant, nil
}
