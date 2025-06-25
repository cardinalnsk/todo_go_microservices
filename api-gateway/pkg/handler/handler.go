package handler

import (
	"api-gateway/pkg/config"
	"api-gateway/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	authURL       string
	todoURL       string
	jwtSigningKey string
}

func NewHandler(
	gatewayCfg config.GatewayConfig,
	authCfg config.UpstreamServiceConfig,
	todoCfg config.UpstreamServiceConfig,
) *Handler {
	return &Handler{
		authURL:       authCfg.Url,
		todoURL:       todoCfg.Url,
		jwtSigningKey: gatewayCfg.JwtSigningKey,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(logger.LogMiddleware())
	r.Use(gin.Recovery())

	// AUTH proxy
	r.POST("/auth/sign-up", h.proxyToAuth)
	r.POST("/auth/sign-in", h.proxyToAuth)

	// TODO proxy с аутентификацией
	api := r.Group("/api", h.AuthMiddleware)
	{
		lists := api.Group("/lists")
		{
			lists.GET("/", h.proxyToTodo)
			lists.POST("/", h.proxyToTodo)
			lists.GET("/:id", h.proxyToTodo)
			lists.PUT("/:id", h.proxyToTodo)
			lists.DELETE("/:id", h.proxyToTodo)

			items := lists.Group("/:id/items")
			{
				items.GET("/", h.proxyToTodo)
				items.POST("/", h.proxyToTodo)
			}
		}
		items := api.Group("/items")
		{
			items.GET("/:id", h.proxyToTodo)
			items.PUT("/:id", h.proxyToTodo)
			items.DELETE("/:id", h.proxyToTodo)
		}
	}

	r.NoRoute(h.notFoundHandler)
	return r
}
