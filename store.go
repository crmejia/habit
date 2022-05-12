package habit

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"time"
)

//things to be done:
//make generic Load, Write that wraps file and db
//create table if not exist
//insert new habit
//update habit

type Storable interface {
	Load() (Tracker, error)
	Write(tracker *Tracker) error
}

type FileStore struct {
	filename string
	tracker  Tracker
}

func NewFileStore(filename string) FileStore {
	return FileStore{filename: filename}
}

func (s FileStore) Load() (Tracker, error) {
	if s.tracker == nil {
		trackerFile, err := os.OpenFile(s.filename, os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			return nil, err
		}
		defer trackerFile.Close()

		fileBytes, err := ioutil.ReadAll(trackerFile)
		if err != nil {
			return nil, err
		}
		ht := Tracker{}
		if len(fileBytes) > 0 {
			err = json.Unmarshal(fileBytes, &ht)
			if err != nil {
				return nil, err
			}
		}
		s.tracker = ht
		return ht, nil
	}
	return s.tracker, nil
}

func (s FileStore) Write(tracker *Tracker) error {
	trackerFile, err := os.OpenFile(s.filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer trackerFile.Close()

	fileBytes, err := json.Marshal(tracker)
	if err != nil {
		return err
	}
	trackerFile.Truncate(0)
	trackerFile.Seek(0, 0)
	_, err = trackerFile.Write(fileBytes)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	trackerFile.Close()
	return nil
}

type DBStore struct {
	dbSource string
}

func NewDBStore(dbSource string) DBStore {
	return DBStore{dbSource: dbSource}
}

func (s DBStore) Load() (Tracker, error) {
	db, err := sql.Open("sqlite3", s.dbSource)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(getAllHabits)
	if err != nil {
		return nil, err
	}
	tracker := Tracker{}
	for rows.Next() {
		var (
			name     string
			streak   int
			interval int64
		)
		err = rows.Scan(&name, &streak, &interval)
		if err != nil {
			return nil, err
		}
		habit := Habit{
			Name:   name,
			Streak: streak,
			//DueDate:  time.Time{},
			Interval: time.Duration(interval),
			//Message:  "",
		}
		tracker[habit.Name] = &habit
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tracker, nil
}

func (s DBStore) Write(tracker *Tracker) error {
	db, err := sql.Open("sqlite3", s.dbSource)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(createTable)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(insertHabit)
	if err != nil {
		return err
	}

	for _, habit := range *tracker {
		_, err := stmt.Exec(habit.Name, habit.Streak, int64(habit.Interval))
		if err != nil {
			return err
		}
	}
	return nil
}

const createTable = `
CREATE TABLE IF NOT EXISTS habit(
id INTEGER NOT NULL PRIMARY KEY,
name VARCHAR NOT NULL,
streak INTEGER NOT NULL,
interval INTEGER NOT NULL
);`

//todo
//dueDate DATETIME NOT NULL,

const getAllHabits = `
SELECT name,streak,interval FROM habit;
`
const insertHabit = `
INSERT INTO habit(name,streak,interval) VALUES(?,?,?)
`
