package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"yc-w22-dating-app-valdy/internal/model"
	ierror "yc-w22-dating-app-valdy/pkg/error"
)

type handler struct {
	authHandler    *authHandler
	onboardHandler *onboardHandler
}

func (s *server) SetupHandlers() {
	s.h = &handler{
		authHandler:    newAuthHandler(s.di.AuthService),
		onboardHandler: newOnboardHandler(s.di.OnboardService),
	}
}

func successResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, model.ApiResponse{
		Code:    "000",
		Message: "success",
		Data:    data,
	})
}

func errorResponse(c echo.Context, err error, data interface{}) error {
	iErr := ierror.ExtractError(err)

	return c.JSON(iErr.HttpCode, model.ApiResponse{
		Code:    iErr.Code,
		Message: iErr.Message,
		Error:   iErr.Description,
		Data:    data,
	})
}
