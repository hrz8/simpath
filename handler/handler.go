package handler

import (
	"database/sql"

	"github.com/hrz8/simpath/internal/user"
)

type Handler struct {
	DB      *sql.DB
	UserSvc *user.Service
}
