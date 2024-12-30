package http

import (
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"

	"yc-w22-dating-app-valdy/config"
	"yc-w22-dating-app-valdy/internal/model"
	"yc-w22-dating-app-valdy/pkg/constant"
	ierror "yc-w22-dating-app-valdy/pkg/error"
	"yc-w22-dating-app-valdy/pkg/jwt"
)

type middleware struct {
	cfg *config.Configuration
}

func (s *server) SetupMiddlewares(cfg *config.Configuration) {
	if cfg == nil {
		panic("config is nil")
	}

	s.m = &middleware{cfg: cfg}
}

func (m *middleware) JWTMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the token from the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("Authorization header is empty")
			return errorResponse(c, ierror.ErrUnauthorized, nil)
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			log.Printf("Authorization header is invalid")
			return errorResponse(c, ierror.ErrUnauthorized, nil)
		}
		token := tokenParts[1]

		// Validate the token
		claims, err := jwt.ValidateJWT(m.cfg.JWTSecret, token)
		if err != nil {
			log.Printf("Validate error: %v", err)
			return errorResponse(c, ierror.ErrUnauthorized, nil)
		}

		// Store user in context
		c.Set(constant.User, model.UserLogin{
			ID:         cast.ToString(claims["id"]),
			Name:       cast.ToString(claims["name"]),
			Gender:     cast.ToString(claims["gender"]),
			IsVerified: cast.ToBool(claims["is_verified"]),
		})

		// Continue to the next handler
		return h(c)
	}
}
