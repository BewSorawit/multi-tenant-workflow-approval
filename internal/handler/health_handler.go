package handler

import (
	"net/http"

	"github.com/BewSorawit/multi-tenant-workflow-approval/internal/platform"
	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	platform.Success(c, http.StatusOK, "ok", nil)
}
