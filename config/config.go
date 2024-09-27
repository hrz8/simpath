package config

const (
	DatabaseURL = "postgres://postgres:toor@localhost:5432/simpath?sslmode=disable"
	AutoMigrate = true
	AllowClient = "http://localhost:8088"
	// token
	AccessTokenLifetime  = 7200    // 2 hours
	RefreshTokenLifetime = 1209600 // 14 days
	IDTokenLifetime      = 7200
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
	// JWT
	JWTIssuer          = "http://localhost:5001"
	JWTAccessTokenAud  = "resource-server-xyz"
	JWTRefreshTokenAud = "http://localhost:5001"
)
