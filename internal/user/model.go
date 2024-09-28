package user

import (
	"time"

	"github.com/hrz8/simpath/internal/role"
)

type OauthUser struct {
	ID                uint32          `json:"id"`
	Role              *role.OauthRole `json:"role,omitempty"`
	RoleID            uint32          `json:"role_id"`
	RoleName          string          `json:"role_name,omitempty"`
	Email             string          `json:"email"`
	EncryptedPassword string          `json:"-"`
	CreatedAt         time.Time       `json:"-"`
	UpdatedAt         time.Time       `json:"-"`
	PublicID          string          `json:"public_id"`
}
