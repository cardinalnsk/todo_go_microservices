package handler

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func respondOr500(c *gin.Context, payload interface{}, err error) {
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, payload)
}

func isNotFoundErr(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func bindJSON(c *gin.Context, out interface{}) bool {
	if err := c.BindJSON(out); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

func okResponse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
