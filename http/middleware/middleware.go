package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"saqrware.com/chat/data"
)

func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		token := c.Request().Header.Get("Authorization")

		// Check token header
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "MISSING_TOKEN"})
		}

		// Check token with Redis.
		exists, err := data.RedisClient.Exists(data.RCTX, "session:"+token).Result()
		if err != nil || exists == 0 {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "INVALID_TOKEN"})
		}

		// Save redis data in the context
		userData, err := data.RedisClient.HGetAll(data.RCTX, "session:"+token).Result()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "FAILED_TO_RETRIEVE_USER_DATA"})
		}

		c.Set("user", userData)

		return next(c)
	}
}
