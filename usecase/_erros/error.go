package errors

import "errors"

// Errors related to user authentication
var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrTokenCreationFailed    = errors.New("failed to create token")
)

// Errors related to password update
var (
	ErrPasswordUpdateFailed = errors.New("failed to update password")
	ErrPasswordIncorrect    = errors.New("incorrect old password")
)

// Errors related to user requests
var (
	ErrRequestUserNotFound  = errors.New("user not found")
	ErrTokenCreationFailure = errors.New("failed to create reset token")
)

// Errors related to password reset
var (
	ErrInvalidResetToken = errors.New("invalid or expired reset token")
	ErrUserNotFound      = errors.New("user not found")
)

// Errors related to user creation
var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrEmailNotValid          = errors.New("email not valid")
	ErrCreatePasswordHash     = errors.New("failed to hash password")
	ErrCreateVerificationCode = errors.New("failed to create verification code")
	ErrCreateId               = errors.New("failed to create id")
	ErrUserCreationFailed     = errors.New("failed to create user")
	ErrIdIsMissing            = errors.New("ID parameter is missing")
)

// ErrUserUpdateFailed Errors related to user update
var (
	ErrUserUpdateFailed = errors.New("failed to update user")
)

// Errors related to user permissions
var (
	ErrUpdateFailed = errors.New("failed to update user")
	ErrForbidden    = errors.New("permission denied")
)

// Errors related to verification code generation
var (
	ErrGenerateCode = errors.New("error generating verification code")
	ErrTimeInterval = errors.New("please wait before requesting a new code")
)

// ErrInvalidCode Errors related to verification code validation
var (
	ErrInvalidCode = errors.New("invalid email or verification code")
)
