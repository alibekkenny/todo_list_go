package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"todo_list/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	users    *models.UserModel
	tasks    *models.TaskModel
}

func main() {
	dsn := "postgres://web:admin@localhost:5432/todo_list"

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		users:    &models.UserModel{},
		tasks:    &models.TaskModel{},
	}

	db, err := openDB(dsn)
	if err != nil {
		app.errorLog.Fatal(err)
		return
	}

	fmt.Println(db)

	fmt.Println("success")
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
