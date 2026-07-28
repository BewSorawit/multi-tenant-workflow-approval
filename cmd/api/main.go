package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/BewSorawit/multi-tenant-workflow-approval/internal/platform"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Warn(".env file not found, using system environment variables")
	}

	config, err := platform.LoadConfig()
	if err != nil {
		logger.Error(
			"failed to load configuration",
			"error", err,
		)
		os.Exit(1)
	}

	dbPool, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		logger.Error("failed to create database pool", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})

		api.GET("/ready", func(c *gin.Context) {
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

				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status": "not_ready",
					"reason": "database unavailable",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "ready",
			})
		})
	}

	addr := ":" + config.AppPort

	logger.Info("server starting", "address", addr)

	if err := router.Run(addr); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
