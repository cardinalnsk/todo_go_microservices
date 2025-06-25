package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	userIdHeader = "X-User-Id"
)

func getUserId(c *gin.Context) (int, error) {
	//id, ok := c.Get(userCtx)
	s := c.GetHeader(userIdHeader)
	logrus.Debugf("user id: %s", s)
	id, err := strconv.Atoi(s)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	return id, nil
}

func getIntParam(c *gin.Context, param, message string) (int, bool) {
	id, err := strconv.Atoi(c.Param(param))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s%s", message, param))
		return 0, false
	}
	return id, true
}
