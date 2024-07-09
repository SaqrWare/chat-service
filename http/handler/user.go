package handler

import (
	"github.com/labstack/echo/v4"
	dto2 "saqrware.com/chat/http/dto"
	"saqrware.com/chat/service"
)

func RegisterUserRoutes(e *echo.Echo) {
	userGroup := e.Group("/api/v1/user")
	userGroup.POST("/register", registerUserHandler)
	userGroup.POST("/login", userLoginHandler)

}

// register user
func registerUserHandler(c echo.Context) error {
	// parse dto.RegisterUserDto from request body
	var userDto dto2.RegisterUserDto
	if err := c.Bind(&userDto); err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	// call service.RegisterUser
	err := service.RegisterUser(userDto)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	return c.JSON(200, map[string]string{"status": "SUCCESS"})
}

func userLoginHandler(c echo.Context) error {
	var loginDto dto2.UserLoginDto
	if err := c.Bind(&loginDto); err != nil {
		return c.JSON(400, map[string]string{"error": "INVALID_REQUEST"})
	}
	token, err := service.UserLogin(loginDto)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	return c.JSON(200, map[string]string{"token": token})

}
