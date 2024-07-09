package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"os"
	"saqrware.com/chat/data"
	handler2 "saqrware.com/chat/http/handler"
)

func main() {

	e := echo.New()

	// Initiate database session
	data.InitiateCassandraSession()

	// Initiate Redis client
	data.InitiateRedisClient()

	defer data.CassandraSession.Close()
	// Register routes
	handler2.RegisterMessageRoutes(e)
	handler2.RegisterUserRoutes(e)

	// api lists all routes
	e.GET("/api-list", func(c echo.Context) error {
		return c.JSON(200, e.Routes())
	})
	//healthcheck
	e.GET("/health", func(c echo.Context) error { return c.String(200, "OK") })

	// get port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))

}
