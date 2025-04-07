package user_data

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sourava/secfix/internal/models"
	"github.com/sourava/secfix/internal/service/user_data"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLatestUserData_ShouldReturn500_WhenServiceReturnsError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := new(user_data.MockUserDataService)
	mockService.On("GetLatestUserData").Return(nil, errors.New("some error"))
	router.GET("/latest_data", NewUserDataHandler(mockService).GetLatestUserData)

	req, err := http.NewRequest("GET", "/latest_data", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"success":false`)
	assert.Contains(t, w.Body.String(), `"error":"some error"`)
}

func TestGetLatestUserData_ShouldReturn200_WhenServiceReturnsNoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := new(user_data.MockUserDataService)
	mockService.On("GetLatestUserData").Return(&models.UserData{
		OSVersion:      "1.0.0",
		OSQueryVersion: "1.0.0",
		AppsInstalled:  []string{"App1"},
	}, nil)
	router.GET("/latest_data", NewUserDataHandler(mockService).GetLatestUserData)

	req, err := http.NewRequest("GET", "/latest_data", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	assert.Contains(t, w.Body.String(), `"os_version":"1.0.0"`)
	assert.Contains(t, w.Body.String(), `"os_query_version":"1.0.0"`)
	assert.Contains(t, w.Body.String(), `"apps_installed":["App1"]`)
}
