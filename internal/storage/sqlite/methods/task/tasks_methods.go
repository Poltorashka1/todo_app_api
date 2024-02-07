package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"
	"web/internal/storage"

	//"web/internal/sqlite"
	"web/internal/storage/models"
)

// INFO: docs of this function in web/internal/storage/storage.go

func getTasksFromRows(rows *sql.Rows) (*models.Tasks, error) {
	const op = "sqlite.getAllTasksFromRows"
	var allTasks models.Tasks

	for rows.Next() {
		var id int
		var text string
		var tags string
		var due string
		err := rows.Scan(&id, &text, &tags, &due)
		if err != nil {
			return nil, err
		}
		allTasks.Tasks = append(allTasks.Tasks, *models.NewTask(id, text, tags, due))
	}

	if len(allTasks.Tasks) == 0 {
		return nil, storage.ErrorSqlNew(http.StatusNotFound, "tasks not found")
	}
	allTasks.Total = len(allTasks.Tasks)
	return &allTasks, nil
}

func (s *TaskTagMethods) CreateTask(text string, tags []string, dueDate *time.Time) error {
	const op = "sqlite.CreateTask"
	// add task
	res, err := s.DataBase.Exec(`INSERT INTO tasks(text, tags, due) VALUES (?, ?, ?)`, text, strings.Join(tags, "; "), dueDate)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
		return err
	}

	// add tags to task
	for _, tagName := range tags {
		_, err := s.DataBase.Exec(`INSERT INTO task_tags (task_id, tag_name) VALUES (?, ?)`, id, tagName)
		if err != nil {
			s.Log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
			return err
		}
	}
	return nil
}

func (s *TaskTagMethods) GetTasksByDueDate(due *time.Time) (*models.Tasks, error) {
	const op = "sqlite.GetTasksByDueDate"

	query := fmt.Sprintf(`SELECT id, text, tags, due FROM tasks WHERE due = ?`)
	rows, err := s.DataBase.Query(query, due)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
		return nil, err
	}

	return getTasksFromRows(rows)

}

// GetAllTasks returns all tasks
func (s *TaskTagMethods) GetAllTasks() (*models.Tasks, error) {
	const op = "sqlite.GetAllTasks"

	rows, err := s.DataBase.Query(`SELECT id, text, tags, due FROM tasks`)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
		return nil, err
	}

	return getTasksFromRows(rows)
}

func (s *TaskTagMethods) GetTask(id int) (*models.Task, error) {
	const op = "sqlite.GetTask"

	rows, err := s.DataBase.Query(`SELECT id, text, tags, due FROM tasks WHERE id = ?`, id)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
		return nil, err
	}

	for rows.Next() {
		var id int
		var text string
		var tags string
		var due string
		err := rows.Scan(&id, &text, &tags, &due)
		if err != nil {
			s.Log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
			return nil, err
		}
		task := models.NewTask(id, text, tags, due)
		return task, nil
	}
	return nil, storage.ErrorSqlNew(http.StatusNotFound, "task not found")
}

func (s *TaskTagMethods) DeleteTask(args ...string) error {
	const op = "sqlite.Delete"
	var err error
	var result sql.Result

	switch len(args) {
	case 1:
		// delete task by id
		result, err = s.DataBase.Exec(`DELETE FROM tasks WHERE id = ?; DELETE FROM task_tags WHERE task_id = ?;`, args[0], args[0])
	case 0:
		// delete all tasks
		result, err = s.DataBase.Exec(`DELETE FROM tasks; DELETE FROM task_tags;`)
	default:
		return fmt.Errorf("database: delete: expect 0 or 1 args, got %d", len(args))
	}
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return storage.ErrorSqlNew(http.StatusNotFound, "task not found")
	}

	return nil
}
