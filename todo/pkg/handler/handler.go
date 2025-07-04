package handler

import (
	"github.com/gin-gonic/gin"
	"todo/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	lists := router.Group("/lists")
	{
		lists.GET("/", h.getAllLists)
		lists.POST("/", h.createList)
		lists.GET("/:id", h.getListById)
		lists.PUT("/:id", h.updateList)
		lists.DELETE("/:id", h.deleteList)

		items := lists.Group("/:id/items")
		{
			items.GET("/", h.getAllItems)
			items.POST("/", h.createItem)
		}
	}
	items := router.Group("items")
	{
		items.GET("/:id", h.getItemById)
		items.PUT("/:id", h.updateItem)
		items.DELETE("/:id", h.deleteItem)
	}

	return router
}
