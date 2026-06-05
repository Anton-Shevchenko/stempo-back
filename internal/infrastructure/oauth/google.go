package oauth

import (
	"context"
	"errors"
	"os"
	"strings"

	"google.golang.org/api/idtoken"
)

type GoogleUserInfo struct {
	GoogleID string
	Email    string
	Name     string
	Picture  string
}

type GoogleTokenVerifier interface {
	Verify(idToken string) (*GoogleUserInfo, error)
}

type googleTokenVerifier struct {
	clientIDs []string
}

func NewGoogleTokenVerifier() GoogleTokenVerifier {
	raw := os.Getenv("GOOGLE_CLIENT_IDS")
	if raw == "" {
		raw = os.Getenv("GOOGLE_CLIENT_ID")
	}

	var clientIDs []string
	for _, id := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(id); trimmed != "" {
			clientIDs = append(clientIDs, trimmed)
		}
	}

	return &googleTokenVerifier{clientIDs: clientIDs}
}

func (v *googleTokenVerifier) Verify(idToken string) (*GoogleUserInfo, error) {
	if len(v.clientIDs) == 0 {
		return nil, errors.New("Google OAuth is not configured")
	}

	var lastErr error
	for _, clientID := range v.clientIDs {
		payload, err := idtoken.Validate(context.Background(), idToken, clientID)
		if err != nil {
			lastErr = err
			continue
		}

		email, _ := payload.Claims["email"].(string)
		if email == "" {
			return nil, errors.New("Google token does not contain email")
		}

		name, _ := payload.Claims["name"].(string)
		picture, _ := payload.Claims["picture"].(string)

		return &GoogleUserInfo{
			GoogleID: payload.Subject,
			Email:    strings.ToLower(email),
			Name:     name,
			Picture:  picture,
		}, nil
	}

	if lastErr != nil {
		return nil, errors.New("invalid Google token")
	}

	return nil, errors.New("invalid Google token")
}
