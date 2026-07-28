package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	databaseUrl := "postgres://postgres:password@localhost:5432/approval?sslmode=disable"

	dbPool, err := pgxpool.New(context.Background(), databaseUrl)

	if err != nil {
		logger.Error("failed to create database pool", "error", err)
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

		api.GET("ready", func(c *gin.Context) {
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
				})

				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "ready",
			})
		})
	}

	addr := ":8080"

	logger.Info(
		"server starting", "address", addr,
	)

	if err := router.Run(addr); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}

}
