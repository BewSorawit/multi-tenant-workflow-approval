package platform

import "github.com/gin-gonic/gin"

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func Success(c *gin.Context, status int, message string, details any) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Message: message,
		Details: details,
	})
}
