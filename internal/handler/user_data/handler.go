package user_data

import (
	"github.com/gin-gonic/gin"
	"github.com/sourava/secfix/internal/service/user_data"
	"net/http"
)

type UserDataHandler struct {
	service user_data.UserDataServiceInterface
}

func NewUserDataHandler(service user_data.UserDataServiceInterface) *UserDataHandler {
	return &UserDataHandler{
		service: service,
	}
}

func (handler *UserDataHandler) GetLatestUserData(c *gin.Context) {
	userData, err := handler.service.GetLatestUserData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userData,
	})
}
