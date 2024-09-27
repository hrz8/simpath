package user

import (
	"time"

	"github.com/hrz8/simpath/internal/role"
)

type OauthUser struct {
	ID                uint32
	Role              *role.OauthRole
	RoleID            uint32
	RoleName          string
	Email             string
	EncryptedPassword string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	PublicID          string
}
