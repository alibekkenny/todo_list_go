package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Todo struct {
	Id          int
	Title       string
	Description string
	Created     time.Time
	Expires     time.Time
	User        User
}

type TodoModel struct {
	DB *pgxpool.Pool
}

// Insert new task into db
func (m *TodoModel) Insert(userId int, title string, description string, expires time.Time) error {
	// returnedInfo := &Todo{}
	stmt := `INSERT INTO task (title, description, creation_date, expire_date, userId)
		VALUES($1,$2,$3,now(),$4) `
	// RETURNING taskId, title, description, creation_date,expire_date,userId
	rowsCreated, err := m.DB.Exec(context.Background(), stmt, title, description, expires, userId)
	//.Scan(&returnedInfo.Id,
	// 	&returnedInfo.Title,
	// 	&returnedInfo.Description,
	// 	&returnedInfo.Created,
	// 	&returnedInfo.Expires,
	// 	&returnedInfo.User.Id,
	// )
	fmt.Println(rowsCreated)
	if err != nil {
		return err
	}
	return nil
}

// Return task by given taskid
func (m *TodoModel) GetById(id int) (*Todo, error) {
	returnedInfo := &Todo{}
	stmt := `SELECT taskId, title, description, creation_date, expire_date, userId FROM task WHERE taskId=$1`
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&returnedInfo.Id,
		&returnedInfo.Title,
		&returnedInfo.Description,
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
func (m *TodoModel) GetByUserId(userId int) ([]*Todo, error) {
	stmt := `SELECT taskId, title, description, creation_date, expire_date, userId FROM task WHERE userId=$1 ORDER BY creation_date ASC`
	rows, err := m.DB.Query(context.Background(), stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := []*Todo{}
	for rows.Next() {
		item := &Todo{}
		err = rows.Scan(&item.Id,
			&item.Title,
			&item.Description,
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
