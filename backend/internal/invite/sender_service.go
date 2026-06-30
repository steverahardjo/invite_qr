package invite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"invite_qr/internal/server"
	"net/http"
	"os"
	"time"

	resend "github.com/resend/resend-go/v3"
	"go.uber.org/zap"
)

// WhatsappSender sends text messages via the Meta/Facebook Graph API.
type WhatsappSender struct {
	apiKey   string
	ourPhone string
}

// EmailSender sends HTML emails via the Resend API.
type EmailSender struct {
	Client   *resend.Client
	ourEmail string
}

// InitWhatsappSender creates a WhatsappSender using the WA_API_KEY environment
// variable. Returns nil if the API key is not set.
func InitWhatsappSender(ourPhone string) *WhatsappSender {
	api_key := os.Getenv("WA_API_KEY")
	if api_key == "" {
		return nil
	}
	return &WhatsappSender{
		apiKey:   api_key,
		ourPhone: ourPhone,
	}
}

// InitEmailSender creates an EmailSender using the RESEND_API_KEY environment
// variable and initializes a Resend client. Returns nil if the key is not set.
func InitEmailSender(ourEmail string) *EmailSender {
	api_key := os.Getenv("RESEND_API_KEY")
	if api_key == "" {
		return nil
	}
	return &EmailSender{
		Client:   resend.NewClient(api_key),
		ourEmail: ourEmail,
	}
}

// WhatsappSend sends a plain-text message to the specified user phone number
// via the Meta Graph API v23.0 messages endpoint using a Bearer token.
func (w *WhatsappSender) WhatsappSend(ctx context.Context, userPhone string, msg string) error {
	logger := server.LoggerFromContext(ctx)
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

// SendEmailInvitation sends an HTML email invitation to the given address
// using the Resend API with a configured subject line.
func (e *EmailSender) SendEmailInvitation(ctx context.Context, emailAddr string, html string) error {

	logger := server.LoggerFromContext(ctx)

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

// GenHTML produces an inline-styled HTML email body with the event title,
// recipient name, and a styled "Open Invitation" link button.
func GenHTML(inviteLink string, recName string, title string) string {

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>%s</title>
</head>

<body>
	<div style="text-align:center; font-family:Arial, sans-serif;">

		<h1>%s</h1>

		<p>
			Hi %s,
		</p>

		<p>
			You are invited to join our special event.
		</p>

		<a
			href="%s"
			style="
				display:inline-block;
				padding:12px 20px;
				background:#333;
				color:white;
				text-decoration:none;
				border-radius:5px;
			"
		>
			Open Invitation
		</a>

	</div>
</body>
</html>
`,
		title,
		title,
		recName,
		inviteLink,
	)
}
