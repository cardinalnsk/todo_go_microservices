package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
	"todo/model"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list model.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var id int
	createListQuery := `INSERT INTO todo_lists (title, description) VALUES ($1, $2) RETURNING id`
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUserListQuery := `INSERT INTO users_lists (user_id, list_id) VALUES ($1, $2)`
	_, err = tx.Exec(createUserListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]model.TodoList, error) {
	var lists []model.TodoList
	query := `SELECT tl.id, tl.title, tl.description 
				FROM todo_lists tl 
				INNER JOIN users_lists ul on tl.id = ul.list_id 
				WHERE ul.user_id = $1`
	err := r.db.Select(&lists, query, userId)
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (r *TodoListPostgres) GetById(userId, id int) (model.TodoList, error) {
	var list model.TodoList
	query := `SELECT tl.id, tl.title, tl.description 
				FROM todo_lists tl 
				    INNER JOIN users_lists ul on tl.id = ul.list_id 
				WHERE ul.user_id = $1 AND tl.id = $2`
	err := r.db.Get(&list, query, userId, id)
	if err != nil {
		return model.TodoList{}, err
	}
	return list, nil
}

func (r *TodoListPostgres) DeleteById(userId, listId int) error {
	query := `DELETE FROM todo_lists tl 
		   USING users_lists ul 
		   WHERE tl.id = ul.list_id 
		     AND ul.user_id = $1 AND ul.list_id = $2`
	_, err := r.db.Exec(query, userId, listId)
	return err
}

func (r *TodoListPostgres) UpdateById(userId, listId int, input model.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE todo_lists tl 
									SET %s 
									FROM users_lists ul 
								WHERE tl.id = ul.list_id 
									AND ul.user_id = $%d AND ul.list_id = $%d`, setQuery, argId, argId+1)
	args = append(args, userId, listId)
	logrus.Debug("updateQuery", query)
	_, err := r.db.Exec(query, args...)
	return err
}
