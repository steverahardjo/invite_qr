// list of service fuunc to send invitations through user preferences
package invitation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"invite_qr/cmd"
	db "invite_qr/db/db_gen"
	"net/http"
	"time"

	resend "github.com/resend/resend-go/v3"
	"go.uber.org/zap"
)

type Service struct {
	queries        *db.Queries
	whatsappSender *WhatsappSender
	emailSender    *EmailSender
	date_limit     time.Time
	hour_limit     time.Time
}

type WhatsappSender struct {
	apiKey   string
	ourPhone string
}

type EmailSender struct {
	Client   *resend.Client
	ourEmail string
}

func GenWebsite() string {

}

func (s *Service) checkAllowedSend() bool {
	now := time.Now()
	return now.Before(s.date_limit) && now.Before(s.hour_limit)
}

// func to send whatsapp message, run inside (s *Service) BulkSend
func (w *WhatsappSender) WhatsappSend(ctx context.Context, userPhone string, msg string) error {
	logger := cmd.LoggerFromContext(ctx)
	logger.Info("sending whatsapp message", zap.String("userPhone", userPhone), zap.Time("time_sent", time.Now()))
	payload := map[string]any{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                userPhone,
		"type":              "text",
		"text": map[string]string{
			"body": msg,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://graph.facebook.com/v23.0/"+w.ourPhone+"/messages",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+w.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf(
			"whatsapp api error: %s",
			resp.Status,
		)
	}

	return nil
}

func (e *EmailSender) SendEmailInvitation(
	ctx context.Context,
	emailAddr string,
	html string,
) error {

	logger := cmd.LoggerFromContext(ctx)

	logger.Info(
		"sending email invitation",
		zap.String("email", emailAddr),
	)

	params := resend.SendEmailRequest{
		From: e.ourEmail,
		To: []string{
			emailAddr,
		},
		Subject: "QR Invitation to our wedding",
		Html:    html,
	}

	sent, err := e.Client.Emails.Send(&params)

	if err != nil {
		logger.Error(
			"failed to send email invitation",
			zap.Error(err),
		)
		return err
	}

	logger.Info(
		"email invitation sent",
		zap.String("id", sent.Id),
		zap.String("receiver", emailAddr),
		zap.Time("time_sent", time.Now()),
	)

	return nil
}
