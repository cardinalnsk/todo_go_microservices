package service

import (
	"todo/model"
	"todo/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId int, list model.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]model.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetById(userId, id int) (model.TodoList, error) {
	return s.repo.GetById(userId, id)
}

func (s *TodoListService) DeleteById(userId, listId int) error {
	return s.repo.DeleteById(userId, listId)
}

func (s *TodoListService) UpdateById(userId, listId int, input model.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateById(userId, listId, input)
}
