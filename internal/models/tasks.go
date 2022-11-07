package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Task struct {
	Id          int
	Title       string
	Description string
	Tag         string
	Created     time.Time
	Expires     time.Time
	User        User
}

type TaskModel struct {
	DB *pgxpool.Pool
}

// Insert new task into db
func (m *TaskModel) Insert(title string, description string, tag string, expires time.Time, userId int) (int, error) {
	stmt := `INSERT INTO task (title, description, tag, creation_date, expire_date, userId
		VALUES($1,$2,$3,now(),$4,$5) RETURNING taskId`
	var returnedId int
	err := m.DB.QueryRow(context.Background(), stmt, title, description, tag, expires, userId).Scan(&returnedId)
	if err != nil {
		return 0, err
	}
	return returnedId, nil
}

// Return task by given taskid
func (m *TaskModel) GetById(id int) (*Task, error) {
	returnedInfo := &Task{}
	stmt := `SELECT * FROM task WHERE taskId=$1`
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&returnedInfo.Id,
		&returnedInfo.Title,
		&returnedInfo.Description,
		&returnedInfo.Tag,
		&returnedInfo.Created,
		&returnedInfo.Expires,
		&returnedInfo.User.Id,
	)
	if err != nil {
		return nil, err
	}
	return returnedInfo, nil
}

// Return tasks by given userid
func (m *TaskModel) GetByUserId(userId int) ([]*Task, error) {
	stmt := `SELECT title,description,tag,creation_date,expire_date,userId FROM task WHERE userId=$1 ORDER BY creation_date DESC`
	rows, err := m.DB.Query(context.Background(), stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := []*Task{}
	for rows.Next() {
		item := &Task{}
		err = rows.Scan(&item.Id,
			&item.Title,
			&item.Description,
			&item.Tag,
			&item.Created,
			&item.Expires,
			&item.User.Id,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
