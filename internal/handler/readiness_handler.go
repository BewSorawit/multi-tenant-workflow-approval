package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/BewSorawit/multi-tenant-workflow-approval/internal/platform"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ReadinessHandler(dbPool *pgxpool.Pool, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(
			c.Request.Context(),
			2*time.Second,
		)
		defer cancel()

		if err := dbPool.Ping(ctx); err != nil {
			logger.Error(
				"database readiness check failed",
				"error", err,
			)

			platform.Error(c, http.StatusServiceUnavailable, "DATABASE_UNAVAILABLE", "Database is not ready", nil)
			return
		}

		platform.Success(c, http.StatusOK, "ready", nil)
	}
}
