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

func TestSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		status         int
		message        string
		details        any
		expectedDetail any
	}{
		{
			name:           "without details",
			status:         http.StatusOK,
			message:        "ok",
			details:        nil,
			expectedDetail: nil,
		},
		{
			name:    "with details",
			status:  http.StatusOK,
			message: "created",
			details: map[string]any{
				"id":   1,
				"name": "workflow",
			},
			expectedDetail: map[string]any{
				"id":   float64(1), // json.Unmarshal แปลง number เป็น float64
				"name": "workflow",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			Success(c, tt.status, tt.message, tt.details)

			require.Equal(t, tt.status, w.Code)

			var resp SuccessResponse
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

			assert.True(t, resp.Success)
			assert.Equal(t, tt.message, resp.Message)

			if tt.expectedDetail == nil {
				assert.Nil(t, resp.Details)
				return
			}

			assert.Equal(t, tt.expectedDetail, resp.Details)
		})
	}
}
