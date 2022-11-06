package main

import (
	"context"
	"log"
	"net/http"
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

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
		return
	}

	defer db.Close()

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		users:    &models.UserModel{DB: db},
		tasks:    &models.TaskModel{DB: db},
	}

	server := http.Server{
		Addr:     ":4000",
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
