package authcode

import "time"

type OauthAuthorizationCode struct {
	ID          uint32
	ClientID    uint32
	UserID      uint32
	Code        string
	RedirectURI string
	Scope       string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
