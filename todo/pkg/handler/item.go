package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todo/model"
)

func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, ok := getIntParam(c, "id", "invalid list id")
	if !ok {
		return
	}

	var input model.TodoItem
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoItem.Create(userId, listId, input)
	if err != nil {
		switch {
		case isNotFoundErr(err):
			newErrorResponse(c, http.StatusNotFound, "not found")
		default:
			newErrorResponse(c, http.StatusInternalServerError, "internal error")
		}
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	listId, ok := getIntParam(c, "id", "invalid list id")
	if !ok {
		return
	}

	items, err := h.services.TodoItem.GetAll(userId, listId)
	if err != nil {
		switch {
		case isNotFoundErr(err):
			newErrorResponse(c, http.StatusNotFound, "not found")
		default:
			newErrorResponse(c, http.StatusInternalServerError, "internal error")
		}
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, ok := getIntParam(c, "id", "invalid item id")
	if !ok {
		return
	}

	item, err := h.services.TodoItem.GetItemById(userId, itemId)
	if err != nil {
		switch {
		case isNotFoundErr(err):
			newErrorResponse(c, http.StatusNotFound, "not found")
		default:
			newErrorResponse(c, http.StatusInternalServerError, "internal error")
		}
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	itemId, ok := getIntParam(c, "id", "invalid item id")
	if !ok {
		return
	}
	err = h.services.TodoItem.Delete(userId, itemId)
	if err != nil {
		switch {
		case isNotFoundErr(err):
			newErrorResponse(c, http.StatusNotFound, "not found")
		default:
			newErrorResponse(c, http.StatusInternalServerError, "internal error")
		}
		return
	}
	okResponse(c)
}

func (h *Handler) updateItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	itemId, ok := getIntParam(c, "id", "invalid item id")
	if !ok {
		return
	}

	var input model.UpdateItemInput
	if !bindJSON(c, &input) {
		return
	}

	err = h.services.TodoItem.Update(userId, itemId, input)
	if err != nil {
		switch {
		case isNotFoundErr(err):
			newErrorResponse(c, http.StatusNotFound, "not found")
		default:
			newErrorResponse(c, http.StatusInternalServerError, "internal error")
		}
		return
	}

	okResponse(c)
}
