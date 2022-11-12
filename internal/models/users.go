package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
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
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	err = m.DB.QueryRow(context.Background(), stmt, nickname, email, hash).Scan(&returnedInfo.Id, &returnedInfo.Nickname, &returnedInfo.Email)
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

func (m *UserModel) CheckUserExists(email string) (bool, error) {
	userExists := 0
	stmt := `SELECT count(*) FROM userinfo WHERE email LIKE $1`
	err := m.DB.QueryRow(context.Background(), stmt, email).Scan(&userExists)
	if err != nil {
		return false, err
	}
	if userExists != 0 {
		return true, nil
	}
	return false, nil
}

// We'll use the Authenticate method to verify whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {

	// Retrieve the id and hashed password associated with the given email. If
	// no matching email exists we return the ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte
	stmt := "SELECT userid, password FROM userinfo WHERE email = $1"
	err := m.DB.QueryRow(context.Background(), stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}
