package platform

import "github.com/gin-gonic/gin"

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

func Error(c *gin.Context, status int, code string, message string, details any) {
	if details == nil {
		details = []any{}
	}
	c.AbortWithStatusJSON(status, ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
