package apperrors

// Defines our user-specific error codes.
const (
	CodeUserInvalidEmailFormat ErrorCode = "USER_INVALID_EMAIL_FORMAT"
	CodeUserEmptyPassword      ErrorCode = "USER_EMPTY_PASSWORD"
	CodeUserEmptyEmail         ErrorCode = "USER_EMPTY_EMAIL"
	CodeUserNotFound           ErrorCode = "USER_NOT_FOUND"
	CodeUserIncorrectPassword  ErrorCode = "USER_INCORRECT_PASSWORD"
	CodeUserInactive           ErrorCode = "USER_INACTIVE"
)

// --- User-specific Error Constructors ---

func NewInvalidEmailFormat(err error) *AppError {
	return NewInvalidInput(err, CodeUserInvalidEmailFormat, "Please enter a valid email format.")
}

func NewEmptyPassword(err error) *AppError {
	return NewInvalidInput(err, CodeUserEmptyPassword, "Password cannot be empty.")
}

func NewEmptyEmail(err error) *AppError {
	return NewInvalidInput(err, CodeUserEmptyEmail, "Email cannot be empty.")
}

func NewUserNotFound(err error) *AppError {
	return &AppError{
		Code:    CodeUserNotFound,
		Message: "Invalid email or password.", // Generic message for security
		Err: err,
	}
}

func NewIncorrectPassword(err error) *AppError {
	return &AppError{
		Code:    CodeUserIncorrectPassword,
		Message: "Invalid email or password.", // Generic message for security
		Err:     err,
	}
}

func NewUserInactive(err error) *AppError {
	return &AppError{
		Code:    CodeUserInactive,
		Message: "This account is inactive.",
		Err:     err,
	}
}
