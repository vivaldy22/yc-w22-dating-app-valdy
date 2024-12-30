package http

import (
	"log"

	"github.com/labstack/echo/v4"

	"yc-w22-dating-app-valdy/internal/model"
	"yc-w22-dating-app-valdy/internal/usecase/swipe"
	ierror "yc-w22-dating-app-valdy/pkg/error"
)

type onboardHandler struct {
	swipeService swipe.Service
}

func newOnboardHandler(swipeService swipe.Service) *onboardHandler {
	if swipeService == nil {
		panic("swipe service is nil")
	}

	return &onboardHandler{
		swipeService: swipeService,
	}
}

func (h *onboardHandler) GetSwipeableProfiles(c echo.Context) error {
	req := model.GetSwipeableProfilesRequest{}
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

	res, err := h.swipeService.GetSwipeableProfiles(c.Request().Context(), req)
	if err != nil {
		return errorResponse(c, err, res)
	}

	return successResponse(c, res)
}

func (h *onboardHandler) SwipePass(c echo.Context) error {
	return nil
}

func (h *onboardHandler) SwipeLike(c echo.Context) error {
	return nil
}
