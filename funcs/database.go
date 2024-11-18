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
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT NOT NULL UNIQUE,
        uname TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL ,
		token TEXT UNIQUE
    );
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL ,
        user_id INTEGER NOT NULL,
        content TEXT NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		content TEXT,
		category TEXT,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS post_interactions (
		user_id INTEGER,
		post_id INTEGER,
		interaction INTEGER,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS comment_interactions (
		user_id INTEGER,
		comment_id INTEGER,
		interaction INTEGER,
		FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	return nil
}
