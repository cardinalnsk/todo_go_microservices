package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"todo/model"
)

type getAllListsResponse struct {
	Data []model.TodoList `json:"data"`
}

func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)

	if err != nil {
		switch {
		case isNotFoundErr(err):
			newErrorResponse(c, http.StatusNotFound, "not found")
		default:
			newErrorResponse(c, http.StatusInternalServerError, "internal error")
		}
		return
	}

	logrus.Info(lists)
	c.JSON(http.StatusOK, getAllListsResponse{
		Data: lists,
	})
}

func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input model.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"list_id": id,
	})
}

func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id, ok := getIntParam(c, "id", "invalid list id")
	if !ok {
		return
	}

	list, err := h.services.TodoList.GetById(userId, id)
	if err != nil {
		switch {
		case isNotFoundErr(err):
			newErrorResponse(c, http.StatusNotFound, "not found")
		default:
			newErrorResponse(c, http.StatusInternalServerError, "internal error")
		}
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, ok := getIntParam(c, "id", "invalid list id")
	if !ok {
		return
	}

	var input model.UpdateListInput
	if !bindJSON(c, &input) {
		return
	}

	if err := h.services.TodoList.UpdateById(userId, listId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	okResponse(c)
}

func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, ok := getIntParam(c, "id", "invalid list id")
	if !ok {
		return
	}

	err = h.services.TodoList.DeleteById(userId, listId)
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
