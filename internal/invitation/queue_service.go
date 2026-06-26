package invitation

import (
	"context"
	sql "database/sql"
	"fmt"
	"invite_qr/cmd"
	db "invite_qr/db/db_gen"
	util "invite_qr/pkg"
	"net/url"
	"time"

	zap "go.uber.org/zap"
)

type Service struct {
	queries        *db.Queries
	whatsappSender *WhatsappSender
	emailSender    *EmailSender
	date_limit     time.Time
	hour_limit     time.Time
	BaseWebURL     *url.URL
	decryptor      *util.IDEncryptor
}

// NewService creates a Service for sending invitation messages.
//
// The service uses the provided database connection to retrieve and
// update invitation records. Invitations may be sent through WhatsApp
// or email depending on the recipient's available contact information.
//
// dateLimit and hourLimit define the allowed sending window.
// baseWebURL is used to generate invitation links included in outgoing
// messages.
func NewService(dbConn *sql.DB, whatsappSender *WhatsappSender, emailSender *EmailSender, dateLimit time.Time, hourLimit time.Time, baseWebURL *url.URL, decryptor *util.IDEncryptor) *Service {
	return &Service{
		queries:        db.New(dbConn),
		whatsappSender: whatsappSender,
		emailSender:    emailSender,
		date_limit:     dateLimit,
		hour_limit:     hourLimit,
		BaseWebURL:     baseWebURL,
		decryptor:      decryptor,
	}
}

func (s *Service) checkAllowedSend() bool {
	now := time.Now()
	return now.Before(s.date_limit) && now.Before(s.hour_limit)
}

func (s *Service) SendInviteOnetime(
	ctx context.Context,
	guestID int32,
	eventTitle string,
	email string,
	waNumber string,
	name string,
) error {

	token := s.decryptor.Encode(guestID, waNumber, email)

	eventURL := fmt.Sprintf(
		"%sevent/%s",
		s.BaseWebURL.String(),
		token,
	)

	if waNumber != "" {
		msg := fmt.Sprintf(
			"Hello, %s! You are invited to %s. Please click the link below to access the event website:\n%s",
			name,
			eventTitle,
			eventURL,
		)

		if err := s.whatsappSender.WhatsappSend(ctx, waNumber, msg); err != nil {
			return err
		}

		return nil
	}

	if email != "" {
		html := GenHTML(eventURL, name, eventTitle)
		return s.emailSender.SendEmailInvitation(ctx, email, html)
	}

	return fmt.Errorf("guest has no supported communication channel")
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

		var sentIDs []int32

		for _, g := range guests {
			select {
			case <-ctx.Done():
				logger.Info("context cancelled")
				return ctx.Err()
			case <-ticker.C:
			}

			if !s.checkAllowedSend() {
				logger.Info("send limit reached")
				return nil
			}

			token := s.decryptor.Encode(g.ID, g.WaNumber, g.Email)

			eventURL := fmt.Sprintf(
				"%sevent/%s",
				s.BaseWebURL.String(),
				token,
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
				logger.Warn(
					"guest has no supported communication channel",
					zap.Int32("guest_id", g.ID),
					zap.String("name", g.Name),
				)
				continue
			}

			if sendErr != nil {
				logger.Error(
					"failed to send invite",
					zap.Error(sendErr),
					zap.Int32("guest_id", g.ID),
					zap.String("channel", channel),
				)
				continue
			}

			sentIDs = append(sentIDs, g.ID)

			logger.Info(
				"invite sent successfully",
				zap.Int32("guest_id", g.ID),
				zap.String("channel", channel),
			)
		}

		if len(sentIDs) > 0 {
			if err := s.queries.MarkInvitesAsSent(ctx, sentIDs); err != nil {
				logger.Error("failed to mark invites as sent", zap.Error(err))
				return err
			}
		}
	}
}
