package sqlite

import (
	"fmt"
	"strings"
	"time"
	"web/internal/storage"
)

// Methods for get tasks by tag

func (s *StoreSqlite) GetTasksByTagFull(tagList []string) (*storage.Tasks, error) {
	const op = "sqlite.GetAllTasksByTag"

	// create query and args for this query
	query, args := buildQueryFull(tagList)
	// get rows from database
	rows, err := s.DataBase.Query(query, args...)

	if err != nil {
		s.Log.Error("%v: %v", op, err.Error())
		return nil, err
	}

	return getTasksFromRows(rows)
}

func (s *StoreSqlite) GetTasksByTagShort(tagList []string) (*storage.Tasks, error) {
	const op = "sqlite.GetTaskByTag"

	// create query and args for this query
	query, args := buildQueryShort(tagList)
	// get rows from database
	rows, err := s.DataBase.Query(query, args...)

	if err != nil {
		s.Log.Error("%v: %v", op, err.Error())
		return nil, err
	}

	return getTasksFromRows(rows)

}

// buildQuery builds query and args for GetAllTasksByTag
func buildQueryShort(tagList []string) (string, []interface{}) {
	var query string
	var args = make([]interface{}, len(tagList))

	// Todo if len == 0
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

// Methods for get tasks by tag and due date

func (s *StoreSqlite) GetTasksByDueAndTagFull(tags []string, dueDate *time.Time) (*storage.Tasks, error) {
	const op = "sqlite.GetTasksByDueAndTag"
	query, args := buildQueryTagDueFull(tags, dueDate)

	rows, err := s.DataBase.Query(query, args...)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
		return nil, err
	}

	return getTasksFromRows(rows)
}

func (s *StoreSqlite) GetTasksByDueAndTagShort(tagList []string, dueDate *time.Time) (*storage.Tasks, error) {
	const op = "sqlite.GetTaskByDueAndTag"

	query, args := buildQueryTagDueShort(tagList, dueDate)

	rows, err := s.DataBase.Query(query, args...)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
		return nil, err
	}
	return getTasksFromRows(rows)
}

func buildQueryTagDueFull(tagList []string, dueDate *time.Time) (string, []interface{}) {
	const op = "sqlite.buildQueryTagDueShort"

	args := make([]interface{}, len(tagList), len(tagList)+2)
	var query string

	switch len(tagList) {
	case 1:
		query = fmt.Sprintf(`
									SELECT t1.id, t1.text, t1.tags, t1.due
									FROM tasks t1
									Join task_tags t2 On t2.task_id = t1.id
									WHERE t2.tag_name = ? AND t1.due = ?;
							`)
		args[0] = tagList[0]
		args = append(args, dueDate)
	default:
		tagsString := strings.Trim(strings.Repeat("?,", len(tagList)), ",")
		query = fmt.Sprintf(`
									SELECT t1.id, t1.text, t1.tags, t1.due
									FROM tasks t1
									Join task_tags t2 On t2.task_id = t1.id
									WHERE t2.tag_name IN (%s) AND t1.due = ?
									GROUP BY t1.id
									HAVING COUNT(DISTINCT t2.tag_name) = ?;
							`, tagsString)
		for k, v := range tagList {
			args[k] = v
		}
		args = append(args, dueDate)
		args = append(args, len(tagList))
	}
	return query, args
}

func buildQueryTagDueShort(tagList []string, dueDate *time.Time) (string, []interface{}) {
	const op = "sqlite.buildQueryTagDueShort"

	args := make([]interface{}, len(tagList), len(tagList)+2)
	var query string

	switch len(tagList) {
	case 1:
		query = fmt.Sprintf(`
									SELECT t1.id, t1.text, t1.tags, t1.due
									FROM tasks t1
									JOIN task_tags t2 ON t1.id = t2.task_id
									WHERE t2.tag_name = ? AND t1.due = ? AND NOT EXISTS (
									SELECT 1
									FROM task_tags t3
									WHERE t3.task_id = t1.id
									AND t3.tag_name <> ?);
							`)
		args[0] = tagList[0]
		args = append(args, dueDate)
		args = append(args, tagList[0])
	default:
		tagsString := strings.Trim(strings.Repeat("?,", len(tagList)), ",")
		//goland:noinspection ALL
		query = fmt.Sprintf(`
								    SELECT t1.id, t1.text, t1.tags, t1.due
									FROM tasks t1
									JOIN task_tags t2 ON t1.id = t2.task_id
									WHERE t2.tag_name IN (%s) AND t1.due = ? AND NOT EXISTS (
										SELECT 1
										FROM task_tags t3
										WHERE t3.task_id = t1.id
										AND t3.tag_name NOT IN (%s))
									GROUP BY t1.id
									HAVING COUNT(DISTINCT t2.tag_name) = ?;
							`, tagsString, tagsString)
		for k, v := range tagList {
			args[k] = v
		}
		args = append(args, dueDate)
		args = append(args, args[:len(tagList)]...)
		args = append(args, len(tagList))
	}
	return query, args
}

func (s *StoreSqlite) GetTasksByTag(tagList []string) (*storage.Tasks, error) {
	tagString := strings.Trim(strings.Repeat("?,", len(tagList)), ",")
	query := fmt.Sprintf(`
										SELECT DISTINCT t1.id,t1.text,t1.tags,t1.due
										FROM tasks t1
										JOIN task_tags t2
										ON t1.id = t2.task_id
										WHERE t2.tag_name IN (%s);
									`, tagString)
	args := make([]interface{}, len(tagList))
	for k, v := range tagList {
		args[k] = v
	}
	rows, err := s.DataBase.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return getTasksFromRows(rows)
}

func (s *StoreSqlite) GetTasksByTagAndDue(tagList []string, dueDate *time.Time) (*storage.Tasks, error) {
	tagString := strings.Trim(strings.Repeat("?,", len(tagList)), ",")

	query := fmt.Sprintf(`
								SELECT DISTINCT t1.id,t1.text,t1.tags,t1.due
								FROM tasks t1
								JOIN task_tags t2
								ON t1.id = t2.task_id
								WHERE t2.tag_name IN (%s) AND t1.due = ?;
								`, tagString)
	args := make([]interface{}, len(tagList))
	for k, v := range tagList {
		args[k] = v
	}
	args = append(args, dueDate)
	rows, err := s.DataBase.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return getTasksFromRows(rows)
}

// Mb realizyu
//--returns tasks that have one of the specified tags but do not have any other.
//SELECT DISTINCT t1.id,t1.text,t1.tags,t1.due
//FROM tasks t1
//JOIN task_tags t2 ON t1.id = t2.task_id
//WHERE t2.tag_name IN ('job', 'important') AND NOT EXISTS (SELECT 1
//FROM task_tags t3
//WHERE t3.task_id = t1.id
//AND t3.tag_name NOT IN ('job', 'important'));
