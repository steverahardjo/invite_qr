package invitation

import (
	"context"
	"invite_qr/cmd"
	db "invite_qr/db/db_gen"
	"net/url"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	queries        *db.Queries
	whatsappSender *WhatsappSender
	emailSender    *EmailSender
	date_limit     time.Time
	hour_limit     time.Time
	BaseWebURL     *url.URL
}

func (s *Service) checkAllowedSend() bool {
	now := time.Now()
	return now.Before(s.date_limit) && now.Before(s.hour_limit)
}

func (s *Service) BulkSendInvite(ctx context.Context, eventTitle string) error {
	logger := cmd.LoggerFromContext(ctx)

	batch_size := int32(50)
	offset := int32(0)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		params := db.GetUnsentInvitesParams{
			Limit:  batch_size,
			Offset: offset,
		}
		guests, err := s.queries.GetUnsentInvites(ctx, params)
		if err != nil {
			logger.Error("failed to fetch unsent participant from db, ", zap.Error(err))
			return err
		}
		if len(guests) == 0 {
			logger.Info("no more unsent invites to send")
			break
		}

		for _, g := range guests {
			if !s.checkAllowedSend() {
				logger.Info("send time limit reached")
				return nil
			}
			select {
			case <-ctx.Done():
				logger.Info("context cancelled db conn")
				return ctx.Err()
			case <-ticker.C:
				eventURL := s.BaseWebURL.String() + "event/" + eventTitle + "/" + g.Name
				if g.WaNumber != "" && g.Sent == false {
					wa_msg := "Hello, " + g.Name + "! You are invited to " + eventTitle + "." + " Please click the link below to access the event website: " + eventURL
					s.whatsappSender.WhatsappSend(ctx, g.WaNumber, wa_msg)
					g.Sent = true
				}

				if g.Email != "" && g.Sent == false {
					GenHTML()
					s.emailSender.SendEmailInvitation()
				}

			}

		}

		offset += batch_size
	}

	return nil
}
