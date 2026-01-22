package apperrors

type ErrorCode string

const(
	CodeUnknown ErrorCode = "UNKNOWN"

	CodeInvalideInput ErrorCode = "INVALID_INPUT"

	CodeUnauthorized ErrorCode = "UNAUTHORIZEd"
	CodeForbidden ErrorCode = "FORBIDDEN"

	CodeNotFound ErrorCode = "NOT_FOUND"
	CodeConflict ErrorCode = "CONFILICT"

	CodeInternal ErrorCode = "INTERNAL"
)

//Http status returns the HTTP status code frot eh given error code
func (c ErrorCode) HTTPStatus() int{
	switch c{
	case CodeInvalideInput:
		return 400 // BAd Request
	case CodeUnauthorized:
		return 401 // UnAuthorized
	case CodeForbidden:
		return 403 // Forbidden
	case CodeNotFound:
		return 404 // Not found
	case CodeConflict:
		return 409 // conflict
	case CodeInternal, CodeUnknown:
		return 500 // Internal server error
	default:
		return 500
	}
}


// apperror is our custom error type
type AppError struct{
	Code ErrorCode `json:"code"`
	Message string `json:"message"`
	Err error
}

// error returns the user facing message
func (e *AppError) Error() string{
	return e.Message
}

// prefined error or custom errors
func NewUnauthorized(err error, message string) *AppError{
	return &AppError{
		Code: CodeUnauthorized,
		Err: err,
		Message: message,
	}
}

func NewInvalidInput(err error,  code ErrorCode, message string) *AppError{
	return &AppError{
		Code: code,
		Message: message,
		Err: err,
	}
}

func NewInternal(err error, message string) *AppError{
	return &AppError{
		Code: CodeInternal,
		Message: message,
		Err: err,
	}
}

func NewNotFound(err error, message string) *AppError{
	return &AppError{
		Code: CodeNotFound,
		Message: message,
		Err: err,
	}
}

func NewConflict(err error, message string) *AppError{
	return &AppError{
		Code: CodeConflict,
		Message: message,
		Err: err,
	}
}