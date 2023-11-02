package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	errInvalidInput  = "invalid input body"
	errInvalidUserId = "invalid user id"
)

type errorResponse struct {
	Message string `json:"message"`
}

type dataResponse struct {
	Data interface{} `json:"data"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
