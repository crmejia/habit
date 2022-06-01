package habit

import (
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"time"
)

type Habit struct {
	Name      string
	Streak    int
	DueDate   time.Time
	Frequency time.Duration
	Message   string
}

type Store interface {
	Get(name string) (*Habit, error)
	Create(habit *Habit) error
	Update(habit *Habit) error
	AllHabits() []*Habit
}

type MemoryStore struct {
	Habits map[string]*Habit
}

var NilHabitError = errors.New("habit cannot be nil")

func OpenMemoryStore() MemoryStore {
	//here a file store or a db store would get the data from persistence.
	memoryStore := MemoryStore{
		Habits: map[string]*Habit{},
	}
	return memoryStore
}

func (s *MemoryStore) Get(name string) (*Habit, error) {
	habit, ok := s.Habits[name]
	if ok {
		return habit, nil
	}
	return nil, nil
}

func (s *MemoryStore) Create(habit *Habit) error {
	if habit == nil {
		return NilHabitError
	}

	if _, ok := s.Habits[habit.Name]; ok {
		return errors.New("habit already exists")
	}
	s.Habits[habit.Name] = habit
	return nil
}

func (s *MemoryStore) Update(habit *Habit) error {
	if habit == nil {
		return NilHabitError
	}

	if _, ok := s.Habits[habit.Name]; !ok {
		return errors.New("cannot update habit does not exists")
	}

	s.Habits[habit.Name] = habit
	return nil
}

func (s MemoryStore) AllHabits() []*Habit {
	allHabits := make([]*Habit, 0, len(s.Habits))
	for _, h := range s.Habits {
		allHabits = append(allHabits, h)
	}
	return allHabits
}

//DBStore represents a store backed by a SQLite DB
type DBStore struct {
	db *sql.DB
}

//OpenDBStore opens a DBStore
func OpenDBStore(dbSource string) (DBStore, error) {
	db, err := sql.Open("sqlite3", dbSource)
	if err != nil {
		return DBStore{}, err
	}

	const createTable = `
CREATE TABLE IF NOT EXISTS habit(
id INTEGER NOT NULL PRIMARY KEY,
name VARCHAR UNIQUE NOT NULL,
streak INTEGER NOT NULL,
frequency INTEGER NOT NULL,
duedate TEXT NOT NULL );`
	_, err = db.Exec(createTable, nil)

	if err != nil {
		return DBStore{}, err
	}

	return DBStore{db: db}, nil
}

func (s *DBStore) Get(name string) (*Habit, error) {
	const getHabit = `
SELECT name, streak, frequency, duedate FROM habit WHERE name = ?
`
	rows, err := s.db.Query(getHabit, name)
	if err != nil {
		return nil, err
	}
	h := Habit{}

	for rows.Next() {
		var (
			hname         string
			streak        int
			frequency     int64
			duedateString string
		)
		err = rows.Scan(&hname, &streak, &frequency, &duedateString)
		if err != nil {
			return nil, err
		}
		h.Name = hname
		h.Streak = streak
		h.Frequency = time.Duration(frequency)
		dueDate, err := time.Parse("2006-01-02 15:04:05-07:00", duedateString)
		if err != nil {
			return nil, err
		}
		h.DueDate = dueDate
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if h.Name == "" {
		return nil, nil
	}
	return &h, nil
}

func (s *DBStore) Create(h *Habit) error {
	if h == nil {
		return NilHabitError
	}
	const insertHabit = `
INSERT INTO habit(name,streak,frequency,duedate) VALUES(?,?,?,?)
`
	stmt, err := s.db.Prepare(insertHabit)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(h.Name, h.Streak, int64(h.Frequency), h.DueDate)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBStore) Update(h *Habit) error {
	if h == nil {
		return NilHabitError
	}
	const updateHabit = `
UPDATE habit SET streak = ?, frequency = ?, duedate = ? WHERE NAME = ?
`
	stmt, err := s.db.Prepare(updateHabit)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(h.Streak, int64(h.Frequency), h.DueDate, h.Name)
	if err != nil {
		return err
	}
	return nil
}
func (s *DBStore) AllHabits() []*Habit {

	const getAllHabits = `
SELECT name, streak, frequency, duedate FROM habit
`
	rows, err := s.db.Query(getAllHabits)
	if err != nil {
		return nil
	}
	habits := make([]*Habit, 0)

	for rows.Next() {
		var (
			hname         string
			streak        int
			frequency     int64
			duedateString string
		)
		err = rows.Scan(&hname, &streak, &frequency, &duedateString)
		if err != nil {
			return nil
		}
		dueDate, err := time.Parse("2006-01-02 15:04:05-07:00", duedateString)
		if err != nil {
			return nil
		}
		h := Habit{
			Name:      hname,
			Streak:    streak,
			Frequency: time.Duration(frequency),
			DueDate:   dueDate,
		}
		habits = append(habits, &h)
	}

	if err = rows.Err(); err != nil {
		return nil
	}
	return habits
}

type FileStore struct {
	filename string
	habits   map[string]*Habit
}

func OpenFileStore(filename string) (FileStore, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return FileStore{}, err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return FileStore{}, err
	}
	habits := make(map[string]*Habit)

	if len(fileBytes) > 0 {
		err = json.Unmarshal(fileBytes, &habits)
		if err != nil {
			return FileStore{}, err
		}
	}
	fileStore := FileStore{
		filename: filename,
		habits:   habits,
	}
	return fileStore, nil
}

func (s *FileStore) Get(name string) (*Habit, error) {
	habit, ok := s.habits[name]
	if ok {
		return habit, nil
	}
	return nil, nil
}

func (s *FileStore) Create(habit *Habit) error {
	if habit == nil {
		return NilHabitError
	}

	if _, ok := s.habits[habit.Name]; ok {
		return errors.New("habit already exists")
	}
	s.habits[habit.Name] = habit
	err := s.writeFile()
	if err != nil {
		return err
	}

	return nil

}

func (s *FileStore) Update(habit *Habit) error {
	if habit == nil {
		return NilHabitError
	}

	if _, ok := s.habits[habit.Name]; !ok {
		return errors.New("cannot update habit does not exists")
	}

	s.habits[habit.Name] = habit
	err := s.writeFile()
	return err
}

func (s *FileStore) writeFile() error {
	file, err := os.OpenFile(s.filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	fileBytes, err := json.Marshal(s.habits)
	if err != nil {
		return err
	}
	file.Truncate(0)
	file.Seek(0, 0)
	_, err = file.Write(fileBytes)
	return err
}

func (s *FileStore) AllHabits() []*Habit {
	allHabits := make([]*Habit, 0, len(s.habits))
	for _, h := range s.habits {
		allHabits = append(allHabits, h)
	}
	return allHabits
}
