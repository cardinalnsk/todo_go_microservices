package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"todo/model"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item model.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var id int
	createItemQuery := `INSERT INTO todo_items (title, description) VALUES ($1, $2) RETURNING id`
	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemQuery := `INSERT INTO lists_items (list_id, item_id) VALUES ($1, $2)`
	_, err = tx.Exec(createListItemQuery, listId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]model.TodoItem, error) {
	var items []model.TodoItem
	query := `SELECT ti.id, ti.title, ti.description, ti.done 
				FROM todo_items ti 
    			INNER JOIN lists_items li on ti.id = li.item_id 
				INNER JOIN users_lists ul on li.list_id = ul.list_id
				WHERE ul.user_id = $1 AND ul.list_id = $2`
	err := r.db.Select(&items, query, userId, listId)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, itemId int) (model.TodoItem, error) {
	var item model.TodoItem
	query := `SELECT ti.id, ti.title, ti.description, ti.done 
				FROM todo_items ti 
				INNER JOIN lists_items li on ti.id = li.item_id
				INNER JOIN users_lists ul on li.list_id = ul.list_id
				WHERE ul.user_id = $1 AND ti.id = $2`
	err := r.db.Get(&item, query, userId, itemId)
	if err != nil {
		return model.TodoItem{}, err
	}
	return item, nil
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	query := `DELETE FROM todo_items ti 
       			USING lists_items li, users_lists ul
       			WHERE ti.id = li.item_id 
       			  AND li.list_id = ul.list_id
       			  AND ul.user_id = $1
       			  AND ti.id = $2`
	_, err := r.db.Exec(query, userId, itemId)
	return err
}
func (r *TodoItemPostgres) Update(userId, itemId int, input model.UpdateItemInput) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done = $%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE todo_items ti 
								SET %s 
								FROM lists_items li, users_lists ul 
								WHERE ti.id = li.item_id 
								    AND li.list_id = ul.list_id 
								    AND ul.user_id = $%d 
								    AND ti.id = $%d`,
		setQuery, argId, argId+1)
	args = append(args, userId, itemId)
	_, err := r.db.Exec(query, args...)
	return err
}
