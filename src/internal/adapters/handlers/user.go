package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"project/guidemysteps/src/internal/core/models"
	"project/guidemysteps/src/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc ports.UserSvc
}

func NewUserHandler(rg *gin.RouterGroup, userSvc ports.UserSvc) {
	uh := UserHandler{
		userSvc: userSvc,
	}

	rg.POST("/user/register", uh.createUser)
	rg.PATCH("/user/update", uh.updateUser)
}

func (uh *UserHandler) createUser(c *gin.Context) {
	var userRequest models.User

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
	}

	err := json.Unmarshal(bodyBytes, &userRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user := models.NewUser(userRequest)

	err = uh.userSvc.Create(user) // Create user with default information
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (uh *UserHandler) updateUser(c *gin.Context) {
	var userRequest models.User

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
	}

	err := json.Unmarshal(bodyBytes, &userRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = uh.userSvc.Update(userRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
