package client

import "time"

type OauthClient struct {
	ID           uint32
	ClientID     string // uuid
	ClientSecret string
	RedirectURI  string
	AppName      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
