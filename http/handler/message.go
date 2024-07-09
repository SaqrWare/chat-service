package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"saqrware.com/chat/http/dto"
	"saqrware.com/chat/http/middleware"
	"saqrware.com/chat/service"
	"strconv"
)

func RegisterMessageRoutes(e *echo.Echo) {
	messageGroup := e.Group("/api/v1/message", middleware.AuthenticationMiddleware)
	messageGroup.POST("/send", sendMessageHandler)
	messageGroup.GET("", messageHistoryHandler)

}

// send message
func sendMessageHandler(c echo.Context) error {
	// get user data from context

	userMap, ok := c.Get("user").(map[string]string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "NO_VALID_SENDER"})
	}

	var sendMessageDto dto.SendMessageDto
	if err := c.Bind(&sendMessageDto); err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	// Sender id in session
	senderID, ok := userMap["id"]
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ID_NOT_FOUND"})
	}
	err := service.CreateMessage(sendMessageDto, senderID)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, map[string]string{"status": "SUCCESS"})

}

// get message history
func messageHistoryHandler(c echo.Context) error {
	// Get user data from context
	userMap, ok := c.Get("user").(map[string]string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "NO_VALID_RECEIVER"})
	}

	// Retrieve query parameters
	lastID := c.QueryParam("lastID")
	limitParam := c.QueryParam("limit")

	// Convert limit from string to int
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		// Set a default limit or return an error
		limit = 10 // Example default value
	}

	// Assuming userID is needed and obtained from the userMap
	userID, ok := userMap["id"]
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "USER_ID_NOT_FOUND"})
	}

	// Call the service to get message history
	messages, err := service.GetMessageHistory(userID, lastID, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Return the messages as JSON
	return c.JSON(http.StatusOK, messages)
}
