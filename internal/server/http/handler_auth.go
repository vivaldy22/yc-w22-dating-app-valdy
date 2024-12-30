package http

import (
	"log"

	"github.com/labstack/echo/v4"

	"yc-w22-dating-app-valdy/internal/model"
	"yc-w22-dating-app-valdy/internal/usecase/auth"
	ierror "yc-w22-dating-app-valdy/pkg/error"
)

type authHandler struct {
	authService auth.Service
}

func newAuthHandler(authService auth.Service) *authHandler {
	if authService == nil {
		panic("authService is nil")
	}

	return &authHandler{
		authService: authService,
	}
}

func (h *authHandler) SignUp(c echo.Context) error {
	req := model.SignUpRequest{}
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		return errorResponse(c, ierror.ErrInvalidRequest, nil)
	}

	err = req.Validate()
	if err != nil {
		log.Println(err.Error())
		return errorResponse(c, ierror.ErrInvalidRequest, nil)
	}

	res, err := h.authService.SignUp(c.Request().Context(), req)
	if err != nil {
		return errorResponse(c, err, res)
	}

	return successResponse(c, res)
}

func (h *authHandler) Login(c echo.Context) error {
	req := model.LoginRequest{}
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		return errorResponse(c, ierror.ErrInvalidRequest, nil)
	}

	err = req.Validate()
	if err != nil {
		log.Println(err.Error())
		return errorResponse(c, ierror.ErrInvalidRequest, nil)
	}

	res, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		return errorResponse(c, err, res)
	}

	return successResponse(c, res)
}
