package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var createStmt = `
CREATE TABLE IF NOT EXISTS todos (
	id INTEGER PRIMARY KEY,
	title TEXT NOT NULL,
	notes TEXT NOT NULL,
	location INTEGER NOT NULL,
	checked INTEGER NOT NULL
);
`

var getStmt *sql.Stmt
var getSrc = "SELECT title, notes, location, checked FROM todos WHERE title = ?"

func StoreGetTodo(title string) (*Todo, error) {
	rows, err := getStmt.Query(title)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, rows.Err()
	}

	if !rows.Next() {
		// No results. Okay.
		return nil, nil
	}
	err = rows.Err()
	if err != nil {
		return nil, rows.Err()
	}

	var todo *Todo = &Todo{}

	err = rows.Scan(&todo.Title, &todo.Notes, &todo.Location, &todo.Checked)
	if err != nil {
		rows.Close()
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, rows.Err()
	}

	if rows.Next() {
		return nil, errors.New("multiple rows")
	}
	err = rows.Err()
	if err != nil {
		return nil, rows.Err()
	}

	rows.Close()
	return todo, nil
}

var insertStmt *sql.Stmt
var insertSrc = `
INSERT INTO todos(title, notes, location, checked)
       VALUES(?, ?, ?, ?)
	   ON CONFLICT(id)
	       DO UPDATE
		   SET title=excluded.title,
		       notes=excluded.notes,
			   location=excluded.location,
			   checked=excluded.checked
`

func StoreSetTodo(todo Todo) error {
	_, err := insertStmt.Exec(todo.Title, todo.Notes, todo.Location, todo.Checked)
	return err
}

func StoreConstruct() error {
	db_, err := sql.Open("sqlite3", "gowiki.db")
	if err != nil {
		return err
	}

	_, err = db_.Exec(createStmt)
	if err != nil {
		return err
	}

	getStmt_, err := db_.Prepare(getSrc)
	if err != nil {
		return err
	}

	insertStmt_, err := db_.Prepare(insertSrc)
	if err != nil {
		return err
	}

	db = db_
	getStmt = getStmt_
	insertStmt = insertStmt_
	return nil
}

func StoreDestroy() {
	getStmt.Close()
	db.Close()
}
