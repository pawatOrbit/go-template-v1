package exception

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// ExceptionErrors is used as our project error response.
// All error response will be in this format.
type ExceptionError struct {
	Code           int32
	HttpStatusCode int
	APIStatusCode  int
	GlobalMessage  string
	DebugMessage   string
	// ErrItems       []*ErrItem
	ErrFields        []string
	ErrWithDatas     map[string]string
	Level            Level
	OverrideLogLevel bool

	// This is used to store the stack trace of the error.
	// This is not a field that will be marshalled to JSON.
	StackErrors error
	StackCaller []byte
}

// Error implements go built-in error interface.
// This will output to CommonResponse for our project.
func (cErr *ExceptionError) Error() string {
	// print the entire cErr as json string without stack errors
	cpyCErr := cErr.Copy()
	cpyCErr.StackErrors = nil
	cpyCErr.DebugMessage = ""
	strCpyCErr, err := json.Marshal(cpyCErr)
	if err != nil {
		return cErr.GlobalMessage
	}
	return string(strCpyCErr)
}

// MarshalJSON implements JSON marshaller interface.
// This will marshal only property ErrItems.
//func (cErr *ExceptionError) MarshalJSON() ([]byte, error) {
//	return json.Marshal(cErr.ErrItems)
//}

type ErrItem struct {
	Title       string
	Description string
	Tag         string
}

func (cErr *ErrItem) Error() string {
	return cErr.Description
}

// NewExceptionErrors allocates new empty error item ExceptionErrors
func NewExceptionError(apiStatusCode int, errCode int32, globalMessage string, httpStatusCode int) *ExceptionError {
	return &ExceptionError{
		Code:           errCode,
		HttpStatusCode: httpStatusCode,
		GlobalMessage:  globalMessage,
		APIStatusCode:  apiStatusCode,
	}
}

// This method modifies a new copy of the ExceptionError.
//func (cErr *ExceptionError) WithDetail(errItems ...*ErrItem) *ExceptionError {
//	newErr := cErr.Copy()
//	newErr.ErrItems = errItems
//	return newErr
//}

func (cErr *ExceptionError) WithFields(errItems []string) *ExceptionError {
	newErr := cErr.Copy()
	newErr.ErrFields = errItems
	return newErr
}

func (cErr *ExceptionError) WithDatas(in map[string]string) *ExceptionError {
	newErr := cErr.Copy()
	if newErr.ErrWithDatas == nil {
		newErr.ErrWithDatas = make(map[string]string)
	}

	newErr.ErrWithDatas = in

	return newErr
}

func (cErr *ExceptionError) WithAPIStatusCode(apiStatusCode int) *ExceptionError {
	newErr := cErr.Copy()
	newErr.APIStatusCode = apiStatusCode
	return newErr
}

func (cErr *ExceptionError) WithMessage(globalMessage string) *ExceptionError {
	newErr := cErr.Copy()
	newErr.GlobalMessage = globalMessage
	return newErr
}

func (cErr *ExceptionError) WithDebugMessage(debugMessage string) *ExceptionError {
	newErr := cErr.Copy()
	newErr.DebugMessage = debugMessage
	newErr.StackErrors = errors.New(newErr.GlobalMessage)
	return newErr
}

// This method creates a deep copy of the ExceptionError.
func (cErr *ExceptionError) Copy() *ExceptionError {
	// Copy primitive fields
	copy := *cErr

	// Deep copy the ErrItems slice
	copy.ErrFields = make([]string, len(cErr.ErrFields))
	for i, item := range cErr.ErrFields {
		copy.ErrFields[i] = item
	}

	copy.DebugMessage = cErr.DebugMessage

	return &copy
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type errorField struct {
	Kind    string `json:"kind"`
	Stack   string `json:"stack"`
	Message string `json:"message"`
}

func GetStackField(err error) errorField {
	var stack string

	if serr, ok := err.(stackTracer); ok {
		st := serr.StackTrace()
		stack = fmt.Sprintf("%+v", st)
		if len(stack) > 0 && stack[0] == '\n' {
			stack = stack[1:]
		}
	}
	return errorField{
		Kind:    reflect.TypeOf(err).String(),
		Stack:   stack,
		Message: err.Error(),
	}
}
