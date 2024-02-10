package sqlite

import (
	"contagio/contagio/config/logging"
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

var initst = []string{
	"CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, login TEXT, password TEXT)",
	"CREATE TABLE IF NOT EXISTS allowed (id INTEGER PRIMARY KEY, ip TEXT)",
	"CREATE TABLE IF NOT EXISTS pids(pid TEXT)",
	"CREATE TABLE IF NOT EXISTS sessions(id INTEGER PRIMARY KEY, count TEXT)",
	"CREATE TABLE IF NOT EXISTS bots(id INTEGER PRIMARY KEY, count TEXT)",
	"CREATE TABLE IF NOT EXISTS stats(id INTEGER PRIMARY KEY, inc TEXT, out TEXT)",
}

func InitDb() {
	db, err := sql.Open("sqlite3", "./sqlite/database.db")

	if err != nil {
		os.Create("./sqlite/database.db")
		InitDb()
	}

	Db = db

	for _, i := range initst {
		statement, err := db.Prepare(i)

		if err != nil {
			logging.PrintError("Can't init db: " + err.Error())
			return
		}

		statement.Exec()
	}

}
