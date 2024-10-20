package usecase_error

type IUseCaseError interface {
	StatusCode() int
	Message() string
}

type UseCaseError struct {
	statusCode int
	message    string
}

func NewUseCaseError(statusCode int, message string) IUseCaseError {
	return &UseCaseError{
		statusCode: statusCode,
		message:    message,
	}
}

func (e UseCaseError) StatusCode() int {
	return e.statusCode
}

func (e UseCaseError) Message() string {
	return e.message
}
