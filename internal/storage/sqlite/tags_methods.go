package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"net/http"
	"web/internal/storage"
)

// INFO: docs of this function in web/internal/storage/storage.go

func (s *StoreSqlite) GetAllTags() (*storage.Tags, error) {
	const op = "sqlite.GetAllTags"

	rows, err := s.DataBase.Query(`SELECT id, name FROM tags`)
	if err != nil {
		s.Log.Error("%v: %v", op, err.Error())
		return nil, err
	}

	var allTags storage.Tags

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			s.Log.Error("%v: %v", op, err.Error())
			return nil, err
		}
		allTags.Tags = append(allTags.Tags, *storage.NewTag(id, name))
	}

	if len(allTags.Tags) == 0 {
		return nil, fmt.Errorf("no tags found")
	}

	return &allTags, nil
}

func (s *StoreSqlite) GetTag(name string) (*storage.Tag, error) {
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
		return storage.NewTag(id, name), nil
	}
	return nil, fmt.Errorf("no tags found")
}

func (s *StoreSqlite) CreateTag(name string) error {
	const op = "sqlite.CreateTag"

	_, err := s.DataBase.Exec(`INSERT INTO tags (name) VALUES (?)`, name)
	if err != nil {
		if errSql := err.(sqlite3.Error); errSql.ExtendedCode == sqlite3.ErrConstraintUnique {
			s.Log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
			return ErrorSqliteNew(http.StatusConflict, "tag already exists")
		}
		s.Log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
		return err
	}
	return nil
}

func (s *StoreSqlite) DeleteTag(name ...string) error {
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
		return ErrorSqliteNew(http.StatusNotFound, "task not found")
	}
	return nil
}
