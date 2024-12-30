package http

func (s *server) SetupRoutes() {
	auth := s.di.Echo.Group("/v1/auth")
	{
		auth.POST("/signup", s.h.authHandler.SignUp)
		auth.POST("/login", s.h.authHandler.Login)
	}

	// TODO: middleware not yet
	swipe := s.di.Echo.Group("/v1/onboard/swipe")
	{
		swipe.GET("/profiles", s.h.onboardHandler.GetSwipeableProfiles)
		swipe.POST("/pass", s.h.onboardHandler.SwipePass)
		swipe.POST("/like", s.h.onboardHandler.SwipeLike)
	}
}
