package serving

import (
	"context"
	db_gen "invite_qr/db/db_gen"
	pkg "invite_qr/pkg"
)

type Service struct {
	DB        *db_gen.Queries
	encryptor *pkg.IDEncryptor
}

func (s *Service) GetParticipantName(ctx context.Context, id string) (string, error) {

	decode_id, err := s.encryptor.Decode(id)
	if err != nil {
		return "", err
	}
	participant, err := s.DB.GetParticipantByID(ctx, decode_id)
	if err != nil {
		return "", err
	}
	return participant.Name, nil
}

func (s *Service) UpdateParticipantAccessed(ctx context.Context, participantID string, email string, waNumber string) error {
	int_id, err := s.encryptor.Decode(participantID)
	if err != nil {
		return err
	}
	_, err = s.DB.UpdateParticipantAccessed(ctx, db_gen.UpdateParticipantAccessedParams{
		ID:       int32(int_id),
		Email:    email,
		WaNumber: waNumber,
	})
	return err
}

func (s *Service) MakeInviteQRCode(ctx context.Context, participantID string, qrCodePath string) error {

}
