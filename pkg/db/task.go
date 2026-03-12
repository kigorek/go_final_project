package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// добавление задачи
func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`

	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

// получение всех задач
func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler 
			  ORDER BY date ASC LIMIT :limit`

	rows, err := db.Query(query, sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		var task Task
		var id int64

		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}

		task.ID = fmt.Sprintf("%d", id)
		tasks = append(tasks, &task)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// получение одной задачи
func GetTask(id int64) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler where id = :id`
	row := db.QueryRow(query, sql.Named("id", id))

	var task Task
	var idInt int64

	err := row.Scan(&idInt, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		return nil, err
	}
	task.ID = fmt.Sprintf("%d", idInt)

	return &task, nil
}

// обновление задачи
func UpdateTask(task *Task) error {

	query := `UPDATE scheduler SET date = :date , title = :title, comment = :comment, repeat = :repeat where id = :id`
	res, err := db.Exec(query,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	query := `DELETE FROM scheduler where id = :id`
	_, err = db.Exec(query, sql.Named("id", idInt))

	if err != nil {
		return err
	}
	return nil
}

func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = :date where id = :id`

	res, err := db.Exec(query,
		sql.Named("date", next),
		sql.Named("id", id),
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}
