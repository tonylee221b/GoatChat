package domain

const (
	// User Aggregate Errors
	ErrInvalidUsername          = "invalid username"
	ErrInvalidPhoneNumberFormat = "invalid phone number format"
	ErrInvalidPasswordFormat    = "invalid password format"
	ErrInvalidEmailFormat       = "invalid email format"
	ErrInvalidAudit             = "invalid audit"

	ErrUserNotFound     = "user not found"
	ErrUserAlreadyExist = "user already exists"

	// Session Aggregate Errors
	ErrTokenIssuerRequired    = "issuer is required"
	ErrTokenAudienceRequired  = "audience is required"
	ErrTokenTTLNotPositive    = "ttl must be positive"
	ErrTokenUserIDRequired    = "user id is required"
	ErrTokenSessionIDRequired = "session id is required"
	ErrTokenMissingClaim      = "claim is missing"

	ErrInvalidAccessToken = "invalid access token"

	ErrEmptyRefreshToken       = "refresh token is empty"
	ErrInvalidRefreshToken     = "invalid refresh token"
	ErrInvalidRefreshTokenHash = "invalid refresh token hash"
	ErrGenerateRefreshToken    = "failed to generate refresh token"

	ErrInvalidUserID   = "invalid user id"
	ErrSessionExpired  = "session expired"
	ErrSessionRevoked  = "session revoked"
	ErrSessionNotFound = "session not found"

	ErrInvalidCredentials = "username or password is invalid"

	// General Errors
	ErrDB            = "db error"
	ErrInvalidConfig = "invalid config value"
)
