package apperror

var (
	ErrValidation = NewError("ERR_VALIDATION", "validation error")

	ErrLinkNotFound = NewError("ERR_LINK_NOT_FOUND", "link not found error")
)

type Error struct {
	Code    string
	Message string
}

func NewError(code string, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}
