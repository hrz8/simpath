package client

type OauthClient struct {
	ID           uint32
	ClientID     string
	ClientSecret string
	RedirectURI  string
	AppName      string
}
