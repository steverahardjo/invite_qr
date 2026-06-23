package invitation

import (
	db "invite_qr/db/db_gen"
	"net/url"
	"time"
	"context"
	"go.uber.org/zap"
	"invite_qr/cmd"
	"fmt"
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

	const batchSize int32 = 50
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		guests, err := s.queries.GetUnsentInvites(ctx, db.GetUnsentInvitesParams{
			Limit:  batchSize,
			Offset: 0,
		})
		if err != nil {
			logger.Error("failed to fetch unsent invites", zap.Error(err))
			return err
		}

		if len(guests) == 0 {
			logger.Info("no more unsent invites")
			return nil
		}

		for _, g := range guests {
			if !s.checkAllowedSend() {
				logger.Info("send limit reached")
				return nil
			}

			select {
			case <-ctx.Done():
				logger.Info("context cancelled")
				return ctx.Err()
			case <-ticker.C:
			}

			eventURL := fmt.Sprintf("%sevent/%s/%s",
				s.BaseWebURL.String(),
				url.PathEscape(eventTitle),
				url.PathEscape(g.Name),
			)

			var sendErr error
			var channel string

			if g.WaNumber != "" {
				channel = "whatsapp"
				msg := fmt.Sprintf(
					"Hello, %s! You are invited to %s. Please click the link below to access the event website:\n%s",
					g.Name,
					eventTitle,
					eventURL,
				)
				sendErr = s.whatsappSender.WhatsappSend(ctx, g.WaNumber, msg)
			} else if g.Email != "" {
				channel = "email"
				html := GenHTML(eventURL, g.Name, eventTitle)
				sendErr = s.emailSender.SendEmailInvitation(ctx, g.Email, html)
			} else {
				logger.Warn("guest has no supported communication channel",
					zap.Int32("guest_id", g.ID),
					zap.String("name", g.Name),
				)
				continue
			}

			if sendErr != nil {
				logger.Error("failed to send invite",
					zap.Error(sendErr),
					zap.Int32("guest_id", g.ID),
					zap.String("channel", channel),
				)
				continue
			}

			if err := s.queries.(ctx, g.ID); err != nil {
				logger.Error("invite sent but failed to mark as sent",
					zap.Error(err),
					zap.Int32("guest_id", g.ID),
				)
				continue
			}

			logger.Info("invite sent successfully",
				zap.Int32("guest_id", g.ID),
				zap.String("channel", channel),
			)
		}
	}
}
