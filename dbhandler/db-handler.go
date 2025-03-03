package dbhandler

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const database string = "database.db"

// First run takes a couple of minutes. Create the db and table if they don't exist
func CreateDB() {
  db, err := sql.Open("sqlite3", database)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  _, err = db.Exec("CREATE TABLE IF NOT EXISTS times (ID INTEGER PRIMARY KEY, length TEXT, date TEXT)")
  if err != nil {
    log.Fatal(err)
  }
}

// Saving the time into the database
func InsertTime(length string, currentDate string) {
  db, err := sql.Open("sqlite3", database)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  // I use prepared statements for good practice
  statement, err := db.Prepare("INSERT INTO times(ID, length, date) VALUES(NULL, ?, ?)")
  if err != nil {
    log.Fatal(err)
  }
  _, err = statement.Exec(length, currentDate)
  if err != nil {
    log.Fatal(err)
  }
}
