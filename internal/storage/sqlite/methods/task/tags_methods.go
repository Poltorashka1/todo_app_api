package repository

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"net/http"
	"web/internal/storage"
	"web/internal/storage/models"
)

// INFO: docs of this function in web/internal/storage/storage.go

func (s *TaskTagMethods) GetAllTags() (*models.Tags, error) {
	const op = "sqlite.GetAllTags"
	rows, err := s.DataBase.Query(`SELECT id, name FROM tags`)
	if err != nil {
		s.Log.Error("%v: %v", op, err.Error())
		return nil, err
	}

	var allTags models.Tags

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			s.Log.Error("%v: %v", op, err.Error())
			return nil, err
		}
		allTags.Tags = append(allTags.Tags, *models.NewTag(id, name))
	}

	if len(allTags.Tags) == 0 {
		return nil, fmt.Errorf("no tags found")
	}

	return &allTags, nil
}

func (s *TaskTagMethods) GetTag(name string) (*models.Tag, error) {
	const op = "sqlite.GetTag"

	rows, err := s.DataBase.Query(`SELECT id, name FROM tags WHERE name = ?`, name)
	if err != nil {
		s.Log.Error("%v: %v", op, err.Error())
		return nil, err
	}

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			s.Log.Error("%v: %v", op, err.Error())
			return nil, err
		}
		return models.NewTag(id, name), nil
	}
	return nil, fmt.Errorf("no tags found")
}

func (s *TaskTagMethods) CreateTag(name string) error {
	const op = "sqlite.CreateTag"

	_, err := s.DataBase.Exec(`INSERT INTO tags (name) VALUES (?)`, name)
	if err != nil {
		if errSql := err.(sqlite3.Error); errSql.ExtendedCode == sqlite3.ErrConstraintUnique {
			s.Log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
			return storage.ErrorSqlNew(http.StatusConflict, "tag already exists")
		}
		s.Log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
		return err
	}
	return nil
}

func (s *TaskTagMethods) DeleteTag(name ...string) error {
	const op = "sqlite.DeleteTag"

	var err error
	var result sql.Result

	switch len(name) {
	case 0:
		result, err = s.DataBase.Exec(`DELETE FROM tags`)
	case 1:
		result, err = s.DataBase.Exec(`DELETE FROM tags WHERE name = ?`, name[0])
	default:
		return fmt.Errorf("database: delete: expect 0 or 1 args, got %d", len(name))
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
