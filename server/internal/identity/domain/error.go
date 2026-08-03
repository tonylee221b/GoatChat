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
	ErrInvalidRefreshTokenHash = "invalid refresh token hash"
	ErrGenerateRefreshToken    = "failed to generate refresh token"

	ErrInvalidUserID  = "invalid user id"
	ErrSessionExpired = "session expired"
	ErrSessionRevoked = "session revoked"

	// General Errors
	ErrDB = "db error"
)
