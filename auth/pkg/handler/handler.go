package handler

import (
	"auth/pkg/logger"
	"auth/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(logger.LogMiddleware())

	router.Use(gin.Recovery())

	router.POST("/sign-up", h.signUp)
	router.POST("/sign-in", h.signIn)

	return router
}
