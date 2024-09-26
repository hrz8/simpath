package consent

import "time"

type OauthConsent struct {
	ID        uint32
	ClientID  uint32
	UserID    uint32
	Consent   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
