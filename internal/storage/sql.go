package sql

import (
	"database/sql"
	"fmt"
	"log"
	

	_ "modernc.org/sqlite"
)

func Load() {
	// Connect to or create the database
	db, err := sql.Open("sqlite", "./mydata.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to SQLite!")

	// Create a sample table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS snippets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		code TEXT,
		tags TEXT,
		reference_url TEXT);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %s", err)
	}

	fmt.Println("Table created successfully!")
}
