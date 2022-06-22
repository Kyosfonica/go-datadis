package storage

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rubiojr/go-datadis"
)

const (
	folder   = "./"
	database = "datadis.db"
	driver   = "sqlite3"
)

type Sqlite struct {
	db *sql.DB
}

func NewSqlite() (*Sqlite, error) {
	db, err := open()

	return &Sqlite{db}, err
}

func open() (*sql.DB, error) {
	if !fileExists(database) {
		return createDB()
	}

	db, err := sql.Open(driver, folder+database)
	if err != nil {
		return nil, err
	}

	return db, err
}

func (s *Sqlite) Close() {
	err := s.db.Close()

	if err != nil {
		log.Fatal(err)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func createDB() (*sql.DB, error) {
	db, err := sql.Open(driver, folder+database)
	if err != nil {
		return nil, err
	}

	type Measurement struct {
		Cups         string  `json:"cups"`
		Date         string  `json:"date"`
		Time         string  `json:"time"`
		Consumption  float32 `json:"consumptionKWh"`
		ObtainMethod string  `json:"obtainMethod"`
	}

	measurementsTable, err := db.Prepare("CREATE TABLE IF NOT EXISTS measurements (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, cups TEXT COLLATE BINARY, date DATE COLLATE BINARY, time TEXT COLLATE BINARY, consumptionKWh REAL COLLATE BINARY, obtainMethod TEXT COLLATE BINARY)")
	if err != nil {
		return nil, err
	}
	defer measurementsTable.Close()

	_, err = measurementsTable.Exec()
	if err != nil {
		return nil, err
	}

	return db, err
}

func (s *Sqlite) Reader(date string, time string) (*datadis.Measurement, error) {
	row, err := s.db.Prepare("SELECT * FROM measurements WHERE date = ? AND time = ?")
	if err != nil {
		return nil, err
	}

	measurement := datadis.Measurement{}
	err = row.QueryRow(date, time).Scan(&measurement.Cups, &measurement.Date, &measurement.Time, &measurement.Consumption, &measurement.ObtainMethod)
	if err != nil {
		return nil, err
	}

	return &measurement, nil
}

func (s *Sqlite) Writer(measurement *datadis.Measurement) error {
	rows, err := s.db.Prepare("SELECT date, time, consumptionKWh FROM measurements WHERE date = ? AND time = ?")
	if err != nil {
		return err
	}

	var date string
	var time string
	var consumptionKWh float32

	err = rows.QueryRow(measurement.Date, measurement.Time).Scan(&date, &time, &consumptionKWh)

	if err != nil {
		return insert(s, measurement)
	}

	if measurement.Consumption > consumptionKWh {
		return update(s, measurement)
	}

	return err
}

func insert(s *Sqlite, measurement *datadis.Measurement) error {
	var stmt *sql.Stmt
	var err error

	stmt, err = s.db.Prepare("INSERT INTO measurements(cups, date, time, consumptionKWh, ObtainMethod) VALUES (?,?,?,?,?)")

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(measurement.Cups, measurement.Date, measurement.Time, measurement.Consumption, measurement.ObtainMethod)
	return err
}

func update(s *Sqlite, measurement *datadis.Measurement) error {
	var stmt *sql.Stmt
	var err error
	stmt, err = s.db.Prepare("UPDATE measurements SET cups=?, date=?, time=?, consumptionKWh=?, ObtainMethod=? WHERE date =? AND time=?")

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(measurement.Cups, measurement.Date, measurement.Time, measurement.Consumption, measurement.ObtainMethod, measurement.Date, measurement.Time)
	return err
}
