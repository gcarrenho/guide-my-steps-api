package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"project/guidemysteps/src/internal/core/models"
	"project/guidemysteps/src/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type RoutingHandler struct {
	routingSvc ports.RoutingSvc
	userSvc    ports.UserSvc
}

func NewRoutingHandler(rg *gin.RouterGroup, routingSvc ports.RoutingSvc, userSvc ports.UserSvc) {
	rh := RoutingHandler{
		routingSvc: routingSvc,
		userSvc:    userSvc,
	}

	rg.POST("/route", rh.getRoute)
}

func (r *RoutingHandler) getRoute(c *gin.Context) {
	var routesRequest models.RoutesRequest

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
	}

	err := json.Unmarshal(bodyBytes, &routesRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := r.userSvc.Get(routesRequest.UserEmail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	routesRequest.DrivingMode = "routed-foot" // tambien viene en la request

	rout, err := r.routingSvc.GetRouting(routesRequest, *user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, rout)
}
