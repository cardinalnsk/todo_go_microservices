package repository

import (
	"github.com/jmoiron/sqlx"
	"todo/model"
)

type TodoList interface {
	Create(userId int, list model.TodoList) (int, error)
	GetAll(userId int) ([]model.TodoList, error)
	GetById(userId, listId int) (model.TodoList, error)
	DeleteById(userId, listId int) error
	UpdateById(userId int, listId int, input model.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, input model.TodoItem) (int, error)
	GetAll(userId, listId int) ([]model.TodoItem, error)
	GetById(userId, itemId int) (model.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input model.UpdateItemInput) error
}

type Repository struct {
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		TodoList: NewTodoListPostgres(db),
		TodoItem: NewTodoItemPostgres(db),
	}
}
