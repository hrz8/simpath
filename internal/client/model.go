package client

type OauthClient struct {
	ID           uint32
	ClientID     string // uuid
	ClientSecret string
	RedirectURI  string
	AppName      string
}
