package storage

import (
	"database/sql"
	"fmt"
)

// AddTask добавляет задачу в БД
func (s *TaskStore) AddTask(task Task) (int64, error) {
	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		return 0, fmt.Errorf("Add task failed: %w", err)
	}
	return res.LastInsertId()
}

// GetTasks возвращает список задач
func (s *TaskStore) GetTasks(limit int) ([]*Task, error) {
	var tasks []*Task

	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?"
	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("Get tasks failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		); err != nil {
			return nil, fmt.Errorf("Failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Get tasks failed: %w", err)
	}

	if tasks == nil {
		return []*Task{}, nil
	}

	return tasks, nil
}

// SearchTasks поиск задач
func (s *TaskStore) SearchTasks(search string, limit int) ([]*Task, error) {
	var tasks []*Task

	query := `SELECT date, title, comment, repeat FROM scheduler 
    		WHERE title LIKE :param OR comment LIKE :param 
            ORDER BY date ASC LIMIT :limit`

	rows, err := s.db.Query(query, sql.Named("param", "%"+search+"%"), sql.Named("limit", limit))
	if err != nil {
		return nil, fmt.Errorf("Search tasks failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		); err != nil {
			return nil, fmt.Errorf("Failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Search tasks failed: %w", err)
	}

	if tasks == nil {
		return []*Task{}, nil
	}

	return tasks, nil
}

// SearchTasksByDate поиск задач
func (s *TaskStore) SearchTasksByDate(search string, limit int) ([]*Task, error) {
	var tasks []*Task

	query := `SELECT date, title, comment, repeat FROM scheduler 
			WHERE date LIKE :param 
			ORDER BY date ASC LIMIT :limit`

	rows, err := s.db.Query(query, sql.Named("param", "%"+search+"%"), sql.Named("limit", limit))
	if err != nil {
		return nil, fmt.Errorf("Search tasks failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		); err != nil {
			return nil, fmt.Errorf("Failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Search tasks failed: %w", err)
	}

	if tasks == nil {
		return []*Task{}, nil
	}

	return tasks, nil
}

// GetTaskByID возвращает задачу
func (s *TaskStore) GetTaskByID(id int) (*Task, error) {
	var task Task

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err := s.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Task not found: %w", err)
		}
		return nil, fmt.Errorf("Get task failed: %w", err)
	}

	return &task, nil
}

func (s *TaskStore) UpdateTaskByID(task *Task) error {
	query := "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id"
	res, err := s.db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		return fmt.Errorf("Update task failed: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Update task failed: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("Task not found")
	}

	return nil
}

func (s *TaskStore) DeleteTaskByID(id int) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	res, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Delete task failed: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Delete task failed: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("Task not found")
	}

	return nil
}
