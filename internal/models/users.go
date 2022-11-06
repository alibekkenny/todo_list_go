package models

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	Id       int
	Nickname string
	Email    string
	Password string
}

type UserModel struct {
	DB *pgxpool.Pool
}

// Insert new user into db
func (m *UserModel) Insert(nickname string, email string, password string) (int, error) {
	stmt := `INSERT INTO user (nickname, email, password)
	VALUES($1,$2,$3) RETURNING userId`
	var returnedId int
	err := m.DB.QueryRow(context.Background(), stmt, nickname, email, password).Scan(&returnedId)
	if err != nil {
		return 0, err
	}
	return returnedId, nil
}

// Get user by given userid
func (m *UserModel) Get(id int) (*User, error) {
	returnedInfo := &User{}
	stmt := `SELECT * FROM user WHERE id = &1`
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&returnedInfo)
	if err != nil {
		return nil, err
	}
	return returnedInfo, nil
}
