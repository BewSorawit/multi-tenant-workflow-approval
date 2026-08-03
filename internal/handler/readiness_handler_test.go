package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BewSorawit/multi-tenant-workflow-approval/internal/platform"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPinger struct {
	err error
}

func (m mockPinger) Ping(ctx context.Context) error {
	return m.err
}

func TestReadinessHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name        string
		db          mockPinger
		status      int
		isError     bool
		message     string
		errorCode   string
		errorDetail any
	}{
		{
			name:    "database ready",
			db:      mockPinger{},
			status:  http.StatusOK,
			isError: false,
			message: "ready",
		}, {
			name: "database unavailable",
			db: mockPinger{
				err: errors.New("database unavailable"),
			},
			status:      http.StatusServiceUnavailable,
			isError:     true,
			message:     "Database is not ready",
			errorCode:   "DATABASE_UNAVAILABLE",
			errorDetail: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET(
				"/api/ready",
				ReadinessHandler(tt.db, logger),
			)

			req := httptest.NewRequest(
				http.MethodGet,
				"/api/ready",
				nil,
			)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, tt.status, w.Code)

			if tt.isError {
				var resp platform.ErrorResponse

				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

				assert.False(t, resp.Success)
				assert.Equal(t, tt.errorCode, resp.Error.Code)
				assert.Equal(t, tt.message, resp.Error.Message)
				assert.Equal(t, tt.errorDetail, resp.Error.Details)

				return
			}

			var resp platform.SuccessResponse

			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.True(t, resp.Success)
			assert.Equal(t, tt.message, resp.Message)
			assert.Nil(t, resp.Details)
		})

	}
}
