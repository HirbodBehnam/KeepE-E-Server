package database

import (
	"context"
	"github.com/go-faster/errors"
)

func (db Database) GetTodos(ctx context.Context, userID uint64) ([]Todo, error) {
	rows, err := db.db.Query(ctx, "SELECT id, title, created_at FROM todo WHERE for_user=$1 ORDER BY todo_order", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Read each todo_list
	todos := make([]Todo, 0)
	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID, &todo.Title, &todo.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "cannot read todo row")
		}
		// Read items
		todo.Tasks, err = db.getTodoItems(ctx, todo.ID)
		if err != nil {
			return nil, errors.Wrap(err, "cannot read todo list items")
		}
		// Done
		todos = append(todos, todo)
	}
	return todos, nil
}

func (db Database) GetTodo(ctx context.Context, userID uint64, todoID int32) (Todo, error) {
	var todo Todo
	err := db.db.QueryRow(ctx, "SELECT id, title, created_at FROM todo WHERE id=$1 AND for_user=$2", todoID, userID).
		Scan(&todo.ID, &todo.Title, &todo.CreatedAt)
	if err != nil {
		return Todo{}, err
	}
	todo.Tasks, err = db.getTodoItems(ctx, todoID)
	if err != nil {
		return Todo{}, errors.Wrap(err, "cannot read todo list items")
	}
	return todo, nil
}

// getTodoItems will get all todo_ list items for a todo_ entry
func (db Database) getTodoItems(ctx context.Context, todoID int32) ([]TodoItem, error) {
	rows, err := db.db.Query(ctx, "SELECT todo_text, deadline, done FROM todo_items WHERE for_todo=$1", todoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Read all items
	items := make([]TodoItem, 0)
	for rows.Next() {
		var item TodoItem
		err = rows.Scan(&item.Text, &item.Deadline, &item.Done)
		if err != nil {
			return nil, errors.Wrap(err, "cannot read todo item row")
		}
		items = append(items, item)
	}
	return items, nil
}

// AddTodo will add a new todo_ to database and return the newly added todo_
func (db Database) AddTodo(ctx context.Context, userID uint64, todo Todo) (Todo, error) {
	// Create a transaction to make everything atomic
	tx, err := db.db.Begin(ctx)
	if err != nil {
		return Todo{}, err
	}
	defer tx.Rollback(ctx)
	// Insert the todo_ itself
	err = tx.QueryRow(ctx, "INSERT INTO todo (for_user, title, created_at, todo_order) VALUES ($1, $2, NOW(), (SELECT COALESCE(MAX(todo_order)+1, 0) FROM todo WHERE for_user=$1)) RETURNING id, created_at",
		userID, todo.Title).Scan(&todo.ID, &todo.CreatedAt)
	if err != nil {
		return Todo{}, errors.Wrap(err, "cannot insert todo itself")
	}
	// Insert items
	for _, task := range todo.Tasks {
		_, err = tx.Exec(ctx, "INSERT INTO todo_items (for_todo, todo_text, done, deadline) VALUES ($1, $2, $3, $4)",
			todo.ID, task.Text, task.Done, task.Deadline)
		if err != nil {
			return Todo{}, errors.Wrap(err, "cannot insert todo item")
		}
	}
	// Done
	err = tx.Commit(ctx)
	if err != nil {
		return Todo{}, errors.Wrap(err, "cannot commit")
	}
	// Small fix
	if todo.Tasks == nil {
		todo.Tasks = make([]TodoItem, 0)
	}
	return todo, nil
}

func (db Database) DeleteTodo(ctx context.Context, userID uint64, id int32) error {
	_, err := db.db.Exec(ctx, "DELETE FROM todo WHERE id=$1 AND for_user=$2", id, userID)
	return err
}

func (db Database) ReorderTodo(ctx context.Context, forUser uint64, todo1, todo2 int32) error {
	// Get old orders
	var todo1Order, todo2Order int32
	err := db.db.QueryRow(ctx, "SELECT todo_order FROM todo WHERE id=$1 AND for_user=$2",
		todo1, forUser).Scan(&todo1Order)
	if err != nil {
		return err
	}
	err = db.db.QueryRow(ctx, "SELECT todo_order FROM todo WHERE id=$1 AND for_user=$2",
		todo2, forUser).Scan(&todo2Order)
	if err != nil {
		return err
	}
	// Update them in database. Also, we don't care about errors
	_, _ = db.db.Exec(ctx, "UPDATE todo SET todo_order=$1 WHERE id=$2", todo1Order, todo2)
	_, _ = db.db.Exec(ctx, "UPDATE todo SET todo_order=$1 WHERE id=$2", todo2Order, todo1)
	return nil
}

func (db Database) EditTodo(ctx context.Context, userID uint64, todo Todo) error {
	// Create a transaction to make everything atomic
	tx, err := db.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	// Get ready to see the most bullshit thing I've written.
	// At first, we start by updating the note itself
	_, err = tx.Exec(ctx, "UPDATE todo SET title=$1 WHERE id=$2 AND for_user=$3", todo.Title, todo.ID, userID)
	if err != nil {
		return err
	}
	// Now for the shitty part. We delete every item in notes
	_, err = tx.Exec(ctx, "DELETE FROM todo_items WHERE for_todo=$1", todo.ID)
	if err != nil {
		return errors.Wrap(err, "cannot delete items of todo")
	}
	// Put them again
	for _, task := range todo.Tasks {
		_, err = tx.Exec(ctx, "INSERT INTO todo_items (for_todo, todo_text, done, deadline) VALUES ($1, $2, $3, $4)",
			todo.ID, task.Text, task.Done, task.Deadline)
		if err != nil {
			return errors.Wrap(err, "cannot insert todo item")
		}
	}
	// Done
	return tx.Commit(ctx)
}
