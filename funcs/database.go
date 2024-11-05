package forum

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Createbase() error {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        email TEXT NOT NULL UNIQUE,
        uname TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL ,
		is_active INTEGER DEFAULT 0,
		token TEXT UNIQUE,
		tokenTime TEXT
    );
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL ,
        uname TEXT NOT NULL ,
        content TEXT NOT NULL
    );
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uname TEXT,
		title TEXT,
		content TEXT,
		category TEXT
	);
	CREATE TABLE IF NOT EXISTS interactions (
		username TEXT,
		post_id INTEGER,
		type TEXT,
		interaction INTEGER
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	return nil
}
