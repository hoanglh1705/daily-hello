package middlewares

import (
	"fmt"
	"strings"

	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/pkg/response"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Error(c, appErrors.ErrUnauthorized)
			}

			// Extract Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Error(c, appErrors.ErrUnauthorized)
			}

			tokenStr := parts[1]

			// Parse and validate JWT
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return response.Error(c, appErrors.ErrInvalidToken)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return response.Error(c, appErrors.ErrInvalidToken)
			}

			// Set user info to context
			userID := uint(claims["user_id"].(float64))
			role := claims["role"].(string)

			c.Set("user_id", userID)
			c.Set("role", role)

			return next(c)
		}
	}
}

func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok {
				return response.Error(c, appErrors.ErrUnauthorized)
			}

			for _, r := range allowedRoles {
				if r == role {
					return next(c)
				}
			}

			return response.Error(c, appErrors.ErrForbidden)
		}
	}
}
