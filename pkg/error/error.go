package error

import (
	"net/http"

	"github.com/joomcode/errorx"
)

type Error struct {
	HttpCode    int
	Code        string
	Message     string
	Description string
}

var (
	namespace        = errorx.NewNamespace("yc-w22-dating-app-valdy")
	propHttpCode     = errorx.RegisterProperty("code")
	propCode         = errorx.RegisterProperty("code")
	propMessage      = errorx.RegisterProperty("code")
	propDescription  = errorx.RegisterProperty("code")
	errBase          = errorx.NewType(namespace, "base")
	errBaseTimeout   = errBase.NewSubtype("timeout", errorx.Timeout()).ApplyModifiers()
	errBaseNotFound  = errBase.NewSubtype("notfound", errorx.NotFound()).ApplyModifiers()
	errBaseDuplicate = errBase.NewSubtype("duplicate", errorx.Duplicate()).ApplyModifiers()
)

var (
	ErrInvalidPassword = NewError(http.StatusOK, "101", "invalid password")

	ErrDataNotFound = NewErrorNotFound(http.StatusOK, "801", "data not found")
	ErrDatabase     = NewError(http.StatusOK, "899", "database error")

	ErrInvalidRequest = NewError(http.StatusOK, "901", "invalid request")
	ErrTimeout        = NewErrorTimeout(http.StatusOK, "998", "timeout")
	ErrGeneral        = NewError(http.StatusOK, "999", "general error")
)

func NewError(hc int, code, message string) *errorx.Error {
	return errBase.New(message).WithProperty(propCode, code).WithProperty(propHttpCode, hc)
}

func NewErrorTimeout(hc int, code, message string) *errorx.Error {
	return errBaseTimeout.New(message).WithProperty(propCode, code).WithProperty(propHttpCode, hc)
}

func NewErrorNotFound(hc int, code, message string) *errorx.Error {
	return errBaseNotFound.New(message).WithProperty(propCode, code).WithProperty(propHttpCode, hc)
}

func NewErrorDuplicate(hc int, code, message string) *errorx.Error {
	return errBaseDuplicate.New(message).WithProperty(propCode, code).WithProperty(propHttpCode, hc)
}

func ExtractError(err error) Error {
	if err == nil {
		return Error{
			HttpCode:    http.StatusOK,
			Code:        "000",
			Message:     "success",
			Description: "success",
		}
	}

	var (
		e, ok = err.(*errorx.Error)
	)

	if ok && namespace.IsNamespaceOf(e.Type()) {
		code, httpCode := "500", 500
		hc, ok := errorx.ExtractProperty(e, propHttpCode)
		if ok {
			httpCode = hc.(int)
		}

		c, ok := errorx.ExtractProperty(e, propCode)
		if ok {
			code = c.(string)
		}

		return Error{httpCode, code, e.Message(), e.Error()}
	}

	return Error{
		Code:        "500",
		Message:     "general error",
		Description: err.Error(),
		HttpCode:    500,
	}
}
