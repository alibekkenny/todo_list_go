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
func (m *UserModel) Insert(nickname string, email string, password string) (*User, error) {
	returnedInfo := &User{}
	stmt := `INSERT INTO userinfo (nickname, email, password)
	VALUES($1,$2,$3) RETURNING userId, nickname, email`
	err := m.DB.QueryRow(context.Background(), stmt, nickname, email, password).Scan(&returnedInfo.Id, &returnedInfo.Nickname, &returnedInfo.Email)
	if err != nil {
		return nil, err
	}
	return returnedInfo, nil
}

// Get user by given userid
func (m *UserModel) Get(id int) (*User, error) {
	returnedInfo := &User{}
	stmt := `SELECT userId,nickname,email FROM userinfo WHERE userId = $1`
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&returnedInfo.Id, &returnedInfo.Nickname, &returnedInfo.Email)
	if err != nil {
		return nil, err
	}
	return returnedInfo, nil
}
