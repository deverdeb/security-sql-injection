package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // importer le drive sqlite (go-sqlite3)
)

const databaseFileName = "./tasklist.db"

var database *sql.DB = nil

func GetConnection() *sql.DB {
	if database == nil {
		database = Initialize()
	}
	return database
}

func Initialize() *sql.DB {
	if database != nil {
		return database
	}

	os.Remove(databaseFileName)
	// SQLite is a file based db.
	database, _ = sql.Open("sqlite3", databaseFileName)
	createTables(database)

	return database
}

func Release() {
	if database != nil {
		database.Close()
	}
	database = nil
}

func createTables(db *sql.DB) {
	createUsersRequest := `CREATE TABLE users (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"firstname" TEXT,
		"lastname" TEXT,
		"login" TEXT,
		"password" TEXT,
		"isAdmin" BOOL		
	  );`
	createTable(db, "users", createUsersRequest)

	createTasksRequest := `CREATE TABLE tasks (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"userId" integer,	
		"name" TEXT,
		"description" TEXT,
		"priority" TEXT,
		"status" TEXT,
		"archived" BOOL		
	  );`
	createTable(db, "tasks", createTasksRequest)
}

func createTable(db *sql.DB, name string, sqlRequest string) {
	log.Printf("create %s table", name)
	statement, err := db.Prepare(sqlRequest)
	if err != nil {
		log.Fatalf("failed to prepare request \"%s\"\nerror: %v", sqlRequest, err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalf("failed to execute request \"%s\"\nerror: %v", sqlRequest, err.Error())
	}
	log.Printf("%s table created", name)
}
