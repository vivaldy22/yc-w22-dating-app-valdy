package http

import (
	"log"

	"github.com/labstack/echo/v4"

	"yc-w22-dating-app-valdy/internal/model"
	"yc-w22-dating-app-valdy/internal/usecase/onboard"
	"yc-w22-dating-app-valdy/pkg/constant"
	ierror "yc-w22-dating-app-valdy/pkg/error"
)

type onboardHandler struct {
	swipeService onboard.Service
}

func newOnboardHandler(swipeService onboard.Service) *onboardHandler {
	if swipeService == nil {
		panic("onboard service is nil")
	}

	return &onboardHandler{
		swipeService: swipeService,
	}
}

func (h *onboardHandler) GetSwipeableProfiles(c echo.Context) error {
	u, ok := c.Get(constant.User).(model.UserLogin)
	if !ok {
		log.Println("User not found in context")
		return errorResponse(c, ierror.ErrInvalidRequest, nil)
	}

	req := model.GetSwipeableProfilesRequest{
		UserID: u.ID,
		Gender: u.Gender,
	}

	err := req.Validate()
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
	u, ok := c.Get(constant.User).(model.UserLogin)
	if !ok {
		log.Println("User not found in context")
		return errorResponse(c, ierror.ErrInvalidRequest, nil)
	}

	req := model.SwipeRequest{}
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		return errorResponse(c, ierror.ErrInvalidRequest, nil)
	}

	req.SwiperID = u.ID

	err = req.Validate()
	if err != nil {
		log.Println(err.Error())
		return errorResponse(c, ierror.ErrInvalidRequest, nil)
	}

	res, err := h.swipeService.Swipe(c.Request().Context(), req, constant.ActionPass)
	if err != nil {
		return errorResponse(c, err, res)
	}

	return successResponse(c, res)
}

func (h *onboardHandler) SwipeLike(c echo.Context) error {
	u, ok := c.Get(constant.User).(model.UserLogin)
	if !ok {
		log.Println("User not found in context")
		return errorResponse(c, ierror.ErrInvalidRequest, nil)
	}

	req := model.SwipeRequest{}
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		return errorResponse(c, ierror.ErrInvalidRequest, nil)
	}

	req.SwiperID = u.ID

	err = req.Validate()
	if err != nil {
		log.Println(err.Error())
		return errorResponse(c, ierror.ErrInvalidRequest, nil)
	}

	res, err := h.swipeService.Swipe(c.Request().Context(), req, constant.ActionLike)
	if err != nil {
		return errorResponse(c, err, res)
	}

	return successResponse(c, res)
}

func (h *onboardHandler) BuyPremiumFeature(c echo.Context) error {
	u, ok := c.Get(constant.User).(model.UserLogin)
	if !ok {
		log.Println("User not found in context")
		return errorResponse(c, ierror.ErrInvalidRequest, nil)
	}

	req := model.BuyPremiumFeatureRequest{}
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		return errorResponse(c, ierror.ErrInvalidRequest, nil)
	}

	req.UserID = u.ID

	err = req.Validate()
	if err != nil {
		log.Println(err.Error())
		return errorResponse(c, ierror.ErrInvalidRequest, nil)
	}

	res, err := h.swipeService.BuyPremiumFeature(c.Request().Context(), req)
	if err != nil {
		return errorResponse(c, err, res)
	}

	return successResponse(c, res)
}
