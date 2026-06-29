package invite

import (
	"context"
	sql "database/sql"
	"fmt"
	"invite_qr/internal/server"
	db "invite_qr/db/db_gen"
	"net/url"
	"time"

	zap "go.uber.org/zap"
)

// Service orchestrates the sending of event invitations through WhatsApp
// and/or email, with rate limiting, date/time send windows, and batch processing.
type Service struct {
	queries        *db.Queries
	whatsappSender *WhatsappSender
	emailSender    *EmailSender
	date_limit     time.Time
	hour_limit     time.Time
	BaseWebURL     *url.URL
}

// NewService creates an invitation Service with the given database connection,
// sender implementations, send window limits, and base URL for invite links.
func NewService(dbConn *sql.DB, whatsappSender *WhatsappSender, emailSender *EmailSender, dateLimit time.Time, hourLimit time.Time, baseWebURL *url.URL) *Service {
	return &Service{
		queries:        db.New(dbConn),
		whatsappSender: whatsappSender,
		emailSender:    emailSender,
		date_limit:     dateLimit,
		hour_limit:     hourLimit,
		BaseWebURL:     baseWebURL,
	}
}

// checkAllowedSend returns true if the current time is before both the
// configured date and hour send limits.
func (s *Service) checkAllowedSend() bool {
	now := time.Now()
	return now.Before(s.date_limit) && now.Before(s.hour_limit)
}

// SendInviteOnetime looks up a participant by ID, constructs an event invite
// link using their external UUID, and sends the invitation via WhatsApp (if
// waNumber is provided) or email (if email is provided).
func (s *Service) SendInviteOnetime(
	ctx context.Context,
	guestID int32,
	eventTitle string,
	email string,
	waNumber string,
	name string,
) error {

	participant, err := s.queries.GetParticipantByID(ctx, guestID)
	if err != nil {
		return fmt.Errorf("guest not found: %w", err)
	}

	eventURL := fmt.Sprintf(
		"%sevent/%s",
		s.BaseWebURL.String(),
		participant.ExternalID.String(),
	)

	var sent bool

	if waNumber != "" && s.whatsappSender != nil {
		msg := fmt.Sprintf(
			"Hello, %s! You are invited to %s. Please click the link below to access the event website:\n%s",
			name,
			eventTitle,
			eventURL,
		)

		if err := s.whatsappSender.WhatsappSend(ctx, waNumber, msg); err != nil {
			return err
		}

		sent = true
	}

	if !sent && email != "" && s.emailSender != nil {
		html := GenHTML(eventURL, name, eventTitle)
		return s.emailSender.SendEmailInvitation(ctx, email, html)
	}

	if !sent {
		return fmt.Errorf("no sender available for guest")
	}

	return nil
}

// BulkSendInvite fetches unsent participants in batches of 50, sends each
// invitation with a 500ms throttle between sends, respects the configured
// date/time send windows, and marks successfully sent invites in the database.
// The method respects context cancellation for graceful shutdown.
func (s *Service) BulkSendInvite(ctx context.Context, eventTitle string) error {
	logger := server.LoggerFromContext(ctx)

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

			eventURL := fmt.Sprintf(
				"%sevent/%s",
				s.BaseWebURL.String(),
				g.ExternalID.String(),
			)

			var sent bool

			if g.WaNumber != "" && s.whatsappSender != nil {
				msg := fmt.Sprintf(
					"Hello, %s! You are invited to %s. Please click the link below to access the event website:\n%s",
					g.Name,
					eventTitle,
					eventURL,
				)

				if err := s.whatsappSender.WhatsappSend(ctx, g.WaNumber, msg); err != nil {
					logger.Error("whatsapp failed, falling back to email",
						zap.Error(err),
						zap.Int32("guest_id", g.ID),
					)
				} else {
					sent = true
				}
			}

			if !sent && g.Email != "" && s.emailSender != nil {
				html := GenHTML(eventURL, g.Name, eventTitle)

				if err := s.emailSender.SendEmailInvitation(ctx, g.Email, html); err != nil {
					logger.Error("email also failed",
						zap.Error(err),
						zap.Int32("guest_id", g.ID),
					)
					continue
				}
				sent = true
			}

			if !sent {
				logger.Warn(
					"guest has no supported communication channel",
					zap.Int32("guest_id", g.ID),
					zap.String("name", g.Name),
				)
				continue
			}

			sentIDs = append(sentIDs, g.ID)

			logger.Info(
				"invite sent successfully",
				zap.Int32("guest_id", g.ID),
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
