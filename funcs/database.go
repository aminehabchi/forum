package forum

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db            *sql.DB
	allCategories = map[string]bool{
		"general":       true,
		"news":          true,
		"entertainment": true,
		"hobbies":       true,
		"lifestyle":     true,
		"technology":    true,
	}
)

const (
	usersTables = `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT NOT NULL UNIQUE,
        uname TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	tokensTables = `
    CREATE TABLE IF NOT EXISTS tokens (
        user_id INTEGER PRIMARY KEY AUTOINCREMENT,
        token TEXT UNIQUE,
        created_at DATETIME,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
    `
	postsTable = `
    CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NULL UNIQUE,
        content TEXT NULL UNIQUE,
        user_id INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
    `
	categoriesTable = `
    CREATE TABLE IF NOT EXISTS post_categories (
        post_id INTEGER NOT NULL,
       	category VARCHAR(255) NOT NULL,
        PRIMARY KEY (post_id,category),
        FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
    );
    `
	commentsTable = `
    CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
    `

	postInteractionsTable = `
    CREATE TABLE IF NOT EXISTS post_interactions (
        user_id INTEGER,
        post_id INTEGER,
        interaction INTEGER,
        PRIMARY KEY (user_id, post_id), 
        FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`

	commentInteractionsTable = `
    CREATE TABLE IF NOT EXISTS comment_interactions (
        user_id INTEGER,
        comment_id INTEGER,
        interaction INTEGER,
        PRIMARY KEY (user_id, comment_id), 
        FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`
)

func CreateDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	// Enable foreign key constraints
	if _, err = db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return err
	}

	tables := []struct {
		name   string
		schema string
	}{
		{"users", usersTables},
		{"tokens", tokensTables},
		{"posts", postsTable},
		{"post_categories", categoriesTable},
		{"comments", commentsTable},
		{"post_interactions", postInteractionsTable},
		{"comment_interactions", commentInteractionsTable},
	}

	for _, table := range tables {
		if _, err := db.Exec(table.schema); err != nil {
			return fmt.Errorf("failed to create %s table: %v", table.name, err)
		}
	}

	return nil
}
