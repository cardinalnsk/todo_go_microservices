package service

import (
	"todo/model"
	"todo/pkg/repository"
)

type TodoList interface {
	Create(userId int, list model.TodoList) (int, error)
	GetAll(userId int) ([]model.TodoList, error)
	GetById(userId, listId int) (model.TodoList, error)
	DeleteById(userId, listId int) error
	UpdateById(userId, listId int, input model.UpdateListInput) error
}

type TodoItem interface {
	Create(userId, listId int, input model.TodoItem) (int, error)
	GetAll(userId, listId int) ([]model.TodoItem, error)
	GetItemById(userId, itemId int) (model.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input model.UpdateItemInput) error
}

type Service struct {
	TodoList
	TodoItem
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		TodoList: NewTodoListService(repository.TodoList),
		TodoItem: NewTodoItemService(repository.TodoItem, repository.TodoList),
	}
}
