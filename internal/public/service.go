package public

import (
	"context"
	db_gen "invite_qr/db/db_gen"

	"github.com/google/uuid"
)

type Service struct {
	DB *db_gen.Queries
}

func NewService(db *db_gen.Queries) *Service {
	return &Service{DB: db}
}

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

func (s *Service) UpdateParticipantAccessed(ctx context.Context, externalID string, email string, waNumber string) error {
	id, err := uuid.Parse(externalID)
	if err != nil {
		return err
	}
	participant, err := s.DB.GetParticipantByExternalID(ctx, id)
	if err != nil {
		return err
	}
	_, err = s.DB.UpdateParticipantAccessed(ctx, db_gen.UpdateParticipantAccessedParams{
		ID:       participant.ID,
		Email:    email,
		WaNumber: waNumber,
	})
	return err
}
