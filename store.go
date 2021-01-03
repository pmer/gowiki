package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var db *sql.DB

var createStmt = `
CREATE TABLE IF NOT EXISTS versions (
	created INTEGER PRIMARY KEY,
	data BLOB NOT NULL
);
`

var insertStmt *sql.Stmt
var insertSrc = "INSERT INTO versions(created, data) VALUES(?, ?)"

func StoreSet(data []byte) error {
	created := time.Now().UnixNano()
	_, err := insertStmt.Exec(created, data)
	return err
}

var getStmt *sql.Stmt
var getSrc = "SELECT data FROM versions ORDER BY created DESC LIMIT 1"

func StoreGet() ([]byte, error) {
	rows, err := getStmt.Query()
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

	var data []byte

	err = rows.Scan(&data)
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
	return data, nil
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

	insertStmt_, err := db_.Prepare(insertSrc)
	if err != nil {
		return err
	}

	getStmt_, err := db_.Prepare(getSrc)
	if err != nil {
		return err
	}

	db = db_
	getStmt = getStmt_
	insertStmt = insertStmt_
	return nil
}

func StoreDestroy() {
	insertStmt.Close()
	getStmt.Close()
	db.Close()
}
