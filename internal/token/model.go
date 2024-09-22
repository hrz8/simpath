package token

import (
	"time"
)

type OauthAccessToken struct {
	ID          uint32
	ClientID    uint32
	UserID      uint32
	AccessToken string
	Scope       string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type OauthRefreshToken struct {
	ID           uint32
	ClientID     uint32
	UserID       uint32
	RefreshToken string
	Scope        string
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
