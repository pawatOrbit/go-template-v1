package exception

import "net/http"

type ApplicationErrors struct {
	Debug        bool
	MemberErrors *MockDataServiceErrors
}

func NewApplicationErrors() *ApplicationErrors {
	return &ApplicationErrors{
		Debug:        false,
		MemberErrors: NewMockDataServiceErrors(),
	}
}

// Application Errors Interface
type CommonApplicationErrors interface {
	ThrowUnauthorized() *ExceptionError
	ThrowPermissionDenied() *ExceptionError
	ThrowInvalidRequest() *ExceptionError
}

type MockDataServiceErrors struct {
	CommonApplicationErrors
	ErrUnauthorized     *ExceptionError
	ErrPermissionDenied *ExceptionError
	ErrNotFound         *ExceptionError
	ErrUnableToProceed  *ExceptionError
	ErrInvalidRequest   *ExceptionError
}

func NewMockDataServiceErrors() *MockDataServiceErrors {
	return &MockDataServiceErrors{
		ErrUnauthorized:     NewExceptionError(400, 200000, "Unauthorized", http.StatusUnauthorized),
		ErrPermissionDenied: NewExceptionError(400, 200001, "Permission Denied (Forbidden error)", http.StatusForbidden),
		ErrNotFound:         NewExceptionError(400, 200002, "Not found", http.StatusNotFound),
		ErrUnableToProceed:  NewExceptionError(500, 209999, "Unable to proceed", http.StatusInternalServerError),
		ErrInvalidRequest:   NewExceptionError(500, 210000, "Invalid Request", http.StatusInternalServerError),
	}
}
