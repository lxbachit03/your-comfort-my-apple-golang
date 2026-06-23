package errorcode

type ErrorCode string

const (
	ErrInternalServer      ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"
	ErrCodeValidation      ErrorCode = "VALIDATION_ERROR"
)
