package user

import (
	"github.com/hrz8/simpath/internal/role"
)

type OauthUser struct {
	Role              *role.OauthRole
	RoleID            string
	Email             string
	EncryptedPassword string
}
