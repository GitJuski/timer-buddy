package dbhandler

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Struct for storing queried data
type TimeRow struct {
  ID int
  Length string
  Date string
  Note string
}

const database string = "database.db"

// First run takes a couple of minutes. Create the db and table if they don't exist
func CreateDB() {
  db, err := sql.Open("sqlite3", database)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  _, err = db.Exec("CREATE TABLE IF NOT EXISTS times (ID INTEGER PRIMARY KEY, length TEXT, date TEXT, note TEXT)")
  if err != nil {
    log.Fatal(err)
  }
}

// Saving the time into the database
func InsertTime(length string, currentDate string, note string) {
  db, err := sql.Open("sqlite3", database)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  // I use prepared statements for good practice
  statement, err := db.Prepare("INSERT INTO times(ID, length, date, note) VALUES(NULL, ?, ?, ?)")
  if err != nil {
    log.Fatal(err)
  }
  _, err = statement.Exec(length, currentDate, note)
  if err != nil {
    log.Fatal(err)
  }
}

// Function for querying data filtered by month and year
func GetTimes(month string, year string) []TimeRow {
  db, err := sql.Open("sqlite3", database)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  statement, err := db.Prepare("SELECT * FROM times WHERE strftime('%m', date) = ? AND strftime('%Y', date) = ?")
  if err != nil {
    log.Fatal(err)
  }
  rows, err := statement.Query(month, year)
  if err != nil {
    log.Fatal(err)
  }
  defer rows.Close()

  // A slice for storing struct instances
  var results []TimeRow

  for rows.Next() {
    var timeRow TimeRow // Create a new instance
    err := rows.Scan(&timeRow.ID, &timeRow.Length, &timeRow.Date, &timeRow.Note) // Scan the values into the instance
    if err != nil {
      log.Fatal(err)
    }
    results = append(results, timeRow)
  }
  return results
}
