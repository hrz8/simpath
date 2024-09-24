package user

import (
	"github.com/hrz8/simpath/internal/role"
)

type OauthUser struct {
	ID                uint32
	Role              *role.OauthRole
	RoleID            uint32
	Email             string
	EncryptedPassword string
	PublicID          string
}
