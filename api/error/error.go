package error

import (
	"fmt"
	"net/http"
)

//AppError is struct of application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Err     error  `json:"-"`
}

//Error is error interface impl
func (e *AppError) Error() string {
	if e.Message != "" && e.Err != nil {
		return fmt.Sprintf("%s. Details: %v\n", e.Message, e.Err)
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

//AppErrorf app error format
func AppErrorf(code int, format string, a ...interface{}) *AppError {
	appE := &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
	return appE
}

//NewAppError create new error
func NewAppError(value ...interface{}) *AppError {
	appE := &AppError{}

	for i, v := range value {
		if i >= 3 {
			break
		}

		switch d := v.(type) {
		case int:
			appE.Code = d
		case string:
			appE.Message = d
		case error:
			appE.Detail = d.Error()
			appE.Err = d
		case nil:
		default:
			appE.Message = "Unsupport AppError type"
		}
	}

	if appE.Code == 0 {
		appE.Code = http.StatusInternalServerError
	}
	return appE
}
