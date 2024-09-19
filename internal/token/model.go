package token

import (
	"time"
)

type OauthAccessToken struct {
	ID          uint32
	ClientID    string
	UserID      string
	AccessToken string
	Scope       string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
