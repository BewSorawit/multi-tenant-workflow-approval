package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	logger.Info(
		"configuration loaded",
		"app_env", config.AppEnv,
		"app_port", config.AppPort,
		"shutdown_timeout", config.ShutdownTimeout.String(),
	)

	dbPool, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		logger.Error(
			"failed to create database pool",
			"error", err,
		)
		os.Exit(1)
	}

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

	server := &http.Server{
		Addr:              ":" + config.AppPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("server starting", "address", server.Addr)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	shutdownSignal, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	select {
	case err := <-serverErrors:
		logger.Error("server stopped unexpectedly", "error", err)
		logger.Info("closing database pool")
		dbPool.Close()
		os.Exit(1)

	case <-shutdownSignal.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		config.ShutdownTimeout,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(
			"graceful shutdown failed",
			"error", err,
		)

		if closeErr := server.Close(); closeErr != nil {
			logger.Error(
				"failed to force close server",
				"error", closeErr,
			)
		}

		logger.Info("closing database pool")
		dbPool.Close()
		os.Exit(1)
	}

	logger.Info("closing database pool")
	dbPool.Close()

	logger.Info("server shutdown completed")
}
