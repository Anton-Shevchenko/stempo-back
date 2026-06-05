package email

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const mailjetSendURL = "https://api.mailjet.com/v3.1/send"

type EmailService interface {
	SendInviteEmail(toEmail, toName, inviteLink string) error
}

type mailjetService struct {
	apiKey    string
	secretKey string
	fromEmail string
	fromName  string
	client    *http.Client
}

func NewMailjetService() EmailService {
	return &mailjetService{
		apiKey:    os.Getenv("MAILJET_API_KEY"),
		secretKey: os.Getenv("MAILJET_SECRET_KEY"),
		fromEmail: os.Getenv("MAILJET_FROM_EMAIL"),
		fromName:  os.Getenv("MAILJET_FROM_NAME"),
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

type mjRecipient struct {
	Email string `json:"Email"`
	Name  string `json:"Name,omitempty"`
}

type mjMessage struct {
	From     mjRecipient   `json:"From"`
	To       []mjRecipient `json:"To"`
	Subject  string        `json:"Subject"`
	TextPart string        `json:"TextPart"`
	HTMLPart string        `json:"HTMLPart"`
}

type mjRequest struct {
	Messages []mjMessage `json:"Messages"`
}

func (s *mailjetService) SendInviteEmail(toEmail, toName, inviteLink string) error {
	if toName == "" {
		toName = toEmail
	}

	textBody := fmt.Sprintf(
		"You have been invited to join Stempo as an employee.\n\nTo set up your password and access the app, tap the link below:\n%s\n\nThis link expires in 48 hours.\n\nIf you did not expect this invitation, you can safely ignore this email.",
		inviteLink,
	)

	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;padding:24px;color:#111827;">
  <h2 style="color:#111827;margin-bottom:8px;">You're invited to Stempo</h2>
  <p style="color:#6B7280;margin-bottom:24px;">You have been added as an employee. Tap the button below to set your password and start using the app.</p>
  <a href="%s"
     style="display:inline-block;background:#2563EB;color:#FFFFFF;text-decoration:none;padding:12px 24px;border-radius:8px;font-weight:600;font-size:16px;">
    Set my password
  </a>
  <p style="margin-top:24px;color:#9CA3AF;font-size:13px;">This link expires in 48 hours. If you did not expect this invitation, you can safely ignore this email.</p>
  <hr style="border:none;border-top:1px solid #E5E7EB;margin:32px 0;">
  <p style="color:#9CA3AF;font-size:12px;">Stempo — Loyalty made simple.</p>
</body>
</html>`, inviteLink)

	payload := mjRequest{
		Messages: []mjMessage{
			{
				From: mjRecipient{Email: s.fromEmail, Name: s.fromName},
				To:   []mjRecipient{{Email: toEmail, Name: toName}},
				Subject:  "You're invited to join Stempo",
				TextPart: textBody,
				HTMLPart: htmlBody,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mailjet: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, mailjetSendURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("mailjet: failed to create request: %w", err)
	}

	credentials := base64.StdEncoding.EncodeToString([]byte(s.apiKey + ":" + s.secretKey))
	req.Header.Set("Authorization", "Basic "+credentials)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("mailjet: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("mailjet: unexpected status %d", resp.StatusCode)
	}

	return nil
}
