package config

const (
	DatabaseURL = "postgres://postgres:toor@localhost:5432/simpath?sslmode=disable"
	// token
	AccessTokenLifetime  = 7200    // 2 hours
	RefreshTokenLifetime = 1209600 // 14 days
	AuthCodeLifetime     = 7200
	// session
	SessionSecretKey = "kcp?l2Qh39{89Wq2"
	SessionPath      = "/"
	// SessionMaxAge       = 604800 // 7 days
	SessionMaxAge       = 3600 // 1 hour
	SessionHttpOnly     = true
	SessionName         = "simpath_session"
	UserSessionName     = "simpath_session_user"
	UserDataSessionKey  = "userdata"
	CSRFTokenSessionKey = "csrftoken"
)
