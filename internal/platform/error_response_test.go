package platform

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		status          int
		code            string
		message         string
		details         any
		expectedDetails any
	}{
		{
			name:            "without details",
			status:          http.StatusBadRequest,
			code:            "INVALID_REQUEST",
			message:         "invalid request",
			details:         nil,
			expectedDetails: []any{},
		},
		{
			name:    "with details",
			status:  http.StatusUnauthorized,
			code:    "UNAUTHORIZED",
			message: "authentication required",
			details: map[string]any{
				"field": "token",
			},
			expectedDetails: map[string]any{
				"field": "token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			Error(c, tt.status, tt.code, tt.message, tt.details)

			require.Equal(t, tt.status, w.Code)

			var resp ErrorResponse
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

			assert.False(t, resp.Success)
			assert.Equal(t, tt.code, resp.Error.Code)
			assert.Equal(t, tt.message, resp.Error.Message)
			assert.Equal(t, tt.expectedDetails, resp.Error.Details)
		})
	}
}
