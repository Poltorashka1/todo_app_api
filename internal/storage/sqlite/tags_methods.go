package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"net/http"
	"strings"
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

func (s *StoreSqlite) GetAllTasksByTag(tagList []string) (*storage.Tasks, error) {
	const op = "sqlite.GetAllTasksByTag"

	// create query and args for this query
	query, args := buildQueryFull(tagList)
	// get rows from database
	rows, err := s.DataBase.Query(query, args...)

	if err != nil {
		s.Log.Error("%v: %v", op, err.Error())
		return nil, err
	}

	return getAllTasksFromRows(rows)
}

func (s *StoreSqlite) GetTaskByTag(tagList []string) (*storage.Tasks, error) {
	const op = "sqlite.GetTaskByTag"

	// create query and args for this query
	query, args := buildQueryShort(tagList)
	// get rows from database
	rows, err := s.DataBase.Query(query, args...)

	if err != nil {
		s.Log.Error("%v: %v", op, err.Error())
		return nil, err
	}

	return getAllTasksFromRows(rows)

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

func getAllTasksFromRows(rows *sql.Rows) (*storage.Tasks, error) {
	const op = "sqlite.getAllTasksFromRows"
	var allTasks storage.Tasks

	for rows.Next() {
		var id int
		var text string
		var tags string
		var due string
		err := rows.Scan(&id, &text, &tags, &due)
		if err != nil {
			return nil, err
		}
		allTasks.Tasks = append(allTasks.Tasks, *storage.NewTask(id, text, tags, due))
	}

	if len(allTasks.Tasks) == 0 {
		return nil, ErrorSqliteNew(http.StatusNotFound, "tasks not found")
	}
	return &allTasks, nil
}

// buildQuery builds query and args for GetAllTasksByTag
func buildQueryShort(tagList []string) (string, []interface{}) {
	var query string
	var args = make([]interface{}, len(tagList))

	switch len(tagList) {
	case 1:
		query = `
					SELECT t1.id, t1.text, t1.tags, t1.due
					FROM tasks t1
					JOIN task_tags t2 ON t1.id = t2.task_id
					WHERE tag_name = ?  AND NOT EXISTS (
						SELECT 1
						FROM task_tags t3
						WHERE t3.task_id = t1.id
						AND t3.tag_name <> ?
					)	
				`
		for n, v := range tagList {
			args[n] = v
		}
		args = append(args, args...)
	default:
		tagsString := strings.Trim(strings.Repeat("?,", len(tagList)), ",")
		query = fmt.Sprintf(`
					SELECT t1.id, t1.text, t1.tags, t1.due
					FROM tasks t1
					JOIN task_tags t2 ON t1.id = t2.task_id
					WHERE tag_name IN (%s)  AND NOT EXISTS (
						SELECT 1
						FROM task_tags t3
						WHERE t3.task_id = t1.id
						AND t3.tag_name NOT IN (%s) 
					)
					GROUP BY t1.id, t1.text, t1.tags, t1.due
					HAVING COUNT(DISTINCT tag_name) = ?
				`, tagsString, tagsString)

		for n, v := range tagList {
			args[n] = v
		}

		args = append(args, args...)
		args = append(args, len(tagList))
	}
	return query, args
}

func buildQueryFull(tagList []string) (string, []interface{}) {
	var query string
	var args = make([]interface{}, len(tagList))

	switch len(tagList) {
	case 1:
		query = `
					SELECT t1.id, t1.text, t1.tags, t1.due 
					FROM tasks t1 
				    JOIN task_tags t2 ON t1.id = t2.task_id 
					WHERE tag_name = ?
					`
		for n, v := range tagList {
			args[n] = v
		}
	default:
		tagString := strings.Trim(strings.Repeat("?,", len(tagList)), ",")
		// create query with tags
		query = fmt.Sprintf(`
					SELECT t1.id, t1.text, t1.tags, t1.due
					FROM tasks t1
					JOIN task_tags t2 ON t1.id = t2.task_id
					WHERE tag_name IN (%s)
					GROUP BY t1.id, t1.text, t1.tags, t1.due
					HAVING COUNT(DISTINCT tag_name) = ?;
				`, tagString)

		for n, v := range tagList {
			args[n] = v
		}
		args = append(args, len(tagList))
	}
	return query, args
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

// important job
//tagsString := strings.Trim(strings.Repeat("?,", len(tagList)), ",")
//query := fmt.Sprintf(`
//					SELECT t1.id, t1.text, t1.tags, t1.due
//					FROM tasks t1
//					JOIN task_tags t2 ON t1.id = t2.task_id
//					WHERE tag_name IN (%s)  AND NOT EXISTS (
//						SELECT 1
//						FROM task_tags t3
//						WHERE t3.task_id = t1.id
//						AND t3.tag_name NOT IN (%s)
//					)
//				`, tagsString, tagsString)
//args := make([]interface{}, len(tagList))
//for n, v := range tagList {
//args[n] = v
//}
//args = append(args, args...)
//
//rows, err = s.DataBase.Query(query, args...)
//{
//    "status": 200,
//    "data": [
//        {
//            "tasks": [
//                {
//                    "id": 12,
//                    "text": "eee",
//                    "tags": [
//                        "important"
//                    ],
//                    "due": "2024-03-24T00:00:00Z"
//                },
//                {
//                    "id": 13,
//                    "text": "rrr",
//                    "tags": [
//                        "important"
//                    ],
//                    "due": "2024-03-24T00:00:00Z"
//                },
//                {
//                    "id": 14,
//                    "text": "ggg",
//                    "tags": [
//                        "job"
//                    ],
//                    "due": "2024-03-24T00:00:00Z"
//                },
//                {
//                    "id": 15,
//                    "text": "aaa",
//                    "tags": [
//                        "job",
//                        "important"
//                    ],
//                    "due": "2024-03-24T00:00:00Z"
//                },
//                {
//                    "id": 15,
//                    "text": "aaa",
//                    "tags": [
//                        "job",
//                        "important"
//                    ],
//                    "due": "2024-03-24T00:00:00Z"
//                },
//                {
//                    "id": 16,
//                    "text": "ttt",
//                    "tags": [
//                        "job"
//                    ],
//                    "due": "2024-03-24T00:00:00Z"
//                },
//                {
//                    "id": 17,
//                    "text": "kkk",
//                    "tags": [
//                        "job"
//                    ],
//                    "due": "2024-03-24T00:00:00Z"
//                },
//                {
//                    "id": 18,
//                    "text": "lll",
//                    "tags": [
//                        "job"
//                    ],
//                    "due": "2024-03-24T00:00:00Z"
//                }
//            ]
//        }
//    ]
//}
