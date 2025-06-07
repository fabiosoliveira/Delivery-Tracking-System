package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func initDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	// deliveries (id, status: "pendente" | "em andamento" | "finalizada", empresa_id, motorista_id)
	// locations (id, motorista_id, latitude, longitude, timestamp)
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS Users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		company_id INTEGER
	);
		
	CREATE TABLE IF NOT EXISTS Deliveries (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		status INTEGER NOT NULL,
		recipient TEXT NOT NULL,
		address TEXT NOT NULL,
		company_id INTEGER NOT NULL,
		driver_id INTEGER
	);

	CREATE TABLE IF NOT EXISTS Locations (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		driver_id INTEGER NOT NULL,
		latitude REAL NOT NULL,
		longitude REAL NOT NULL,
		timestamp TEXT NOT NULL
	);
	`

	_, err = DB.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, sqlStmt)
	}
}

func init() {
	initDB()
}
