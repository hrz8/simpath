package handler

import (
	"database/sql"

	"github.com/hrz8/simpath/internal/client"
	"github.com/hrz8/simpath/internal/scope"
	"github.com/hrz8/simpath/internal/token"
	"github.com/hrz8/simpath/internal/user"
	"github.com/hrz8/simpath/session"
)

type Handler struct {
	DB         *sql.DB
	sessionSvc *session.Service
	userSvc    *user.Service
	clientSvc  *client.Service
	scopeSvc   *scope.Service
	tokenSvc   *token.Service
}

func NewHandler(
	db *sql.DB,
	sesSvc *session.Service,
	uSvc *user.Service,
	cSvc *client.Service,
	sSvc *scope.Service,
	tSvc *token.Service,
) *Handler {
	return &Handler{
		db,
		sesSvc,
		uSvc,
		cSvc,
		sSvc,
		tSvc,
	}
}
