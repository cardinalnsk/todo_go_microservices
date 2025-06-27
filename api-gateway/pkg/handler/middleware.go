package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

const (
	authorizationHeader = "Authorization"
	contentTypeHeader   = "Content-Type"
	userCtx             = "userID"
	xUserIdHeader       = "X-User-Id"
	authProxyPrefix     = "/auth"
	apiProxyPrefix      = "/api"
)

// Middleware — вытаскивает userID из токена и кладёт в c.Set().
func (h *Handler) AuthMiddleware(c *gin.Context) {
	tokenStr := c.GetHeader(authorizationHeader)
	if !strings.HasPrefix(tokenStr, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no bearer token"})
		return
	}
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	var cl claims
	token, err := jwt.ParseWithClaims(tokenStr, &cl, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			logrus.Error(err)
			return nil, err
		}
		return h.publicKey, nil
	})
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	c.Set(userCtx, cl.UserID)
	c.Next()
}

func (h *Handler) proxyToAuth(c *gin.Context) {
	path := strings.TrimPrefix(c.Request.URL.Path, authProxyPrefix)
	url := h.authURL + path
	proxyRequest(c, url)
}

func (h *Handler) proxyToTodo(c *gin.Context) {
	userID, _ := c.Get(userCtx)
	uidStr := strconv.Itoa(userID.(int))
	origPath := c.Request.URL.Path

	path := strings.TrimPrefix(origPath, apiProxyPrefix)
	url := h.todoURL + path
	req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "request error"})
		return
	}
	req.Header = c.Request.Header.Clone()
	req.Header.Set(xUserIdHeader, uidStr)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(502, gin.H{"error": "todo service error"})
		return
	}
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, resp.Header.Get(contentTypeHeader), bytes)
}

func (h *Handler) notFoundHandler(c *gin.Context) {
	c.JSON(404, gin.H{
		"error": "not found",
	})
}
func proxyRequest(c *gin.Context, url string) {
	req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "request error"})
		return
	}
	req.Header = c.Request.Header.Clone()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(502, gin.H{"error": "service error"})
		return
	}
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, resp.Header.Get(contentTypeHeader), bytes)
}
