package http

func (s *server) SetupRoutes() {
	auth := s.di.Echo.Group("/v1/auth")
	{
		auth.POST("/signup", s.h.authHandler.SignUp)
		auth.POST("/login", s.h.authHandler.Login)
	}

	onboard := s.di.Echo.Group("/v1/onboard", s.m.JWTMiddleware)
	{
		onboard.GET("/swipe/profiles", s.h.onboardHandler.GetSwipeableProfiles)
		onboard.POST("/swipe/pass", s.h.onboardHandler.SwipePass)
		onboard.POST("/swipe/like", s.h.onboardHandler.SwipeLike)

		onboard.POST("/premium/buy", s.h.onboardHandler.BuyPremiumFeature)
	}
}
