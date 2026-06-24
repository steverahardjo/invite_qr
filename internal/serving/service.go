package serving

import (
	"context"
	db_gen "invite_qr/db/db_gen"
	pkg "invite_qr/pkg"
	"strconv"
)

type Service struct {
	DB        *db_gen.DB
	encryptor *pkg.IDEncryptor
}

func (s *Service) GetParticipantName(ctx context.Context, id int) (string, error) {
	decode_id, err := s.encryptor.Decode(strconv.Itoa(id))
	if err != nil {
		return "", err
	}
	participant, err := s.DB.Participants.GetParticipant(ctx, decode_id)
	if err != nil {
		return "", err
	}
	return participant.Name, nil
}
