package habit

import (
	"database/sql"
	"encoding/json"
	"errors"
	//SQLite driver package
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"time"
)

//Habit a type representing a habit
type Habit struct {
	Name      string
	Streak    int
	DueDate   time.Time
	Frequency time.Duration
	Message   string
}

//Store is an interface that captures the behavior of a Store
type Store interface {
	Get(name string) (*Habit, error)
	Create(habit *Habit) error
	Update(habit *Habit) error
	GetAllHabits() []*Habit
}

//MemoryStore is a type representing an in-memory store
type MemoryStore struct {
	Habits map[string]*Habit
}

//OpenMemoryStore returns a new MemoryStore. Note that other types returns the interface Store
func OpenMemoryStore() MemoryStore {
	//here a file store or a db store would get the data from persistence.
	memoryStore := MemoryStore{
		Habits: map[string]*Habit{},
	}
	return memoryStore
}

//Get searches Store by name and returns the habit if it exists
func (s *MemoryStore) Get(name string) (*Habit, error) {
	habit, ok := s.Habits[name]
	if ok {
		return habit, nil
	}
	return nil, nil
}

//Create inserts the given habit into the store. It returns an error if the habit already exists
func (s *MemoryStore) Create(habit *Habit) error {
	if habit == nil {
		return ErrNilHabit
	}

	if _, ok := s.Habits[habit.Name]; ok {
		return errors.New("habit already exists")
	}
	s.Habits[habit.Name] = habit
	return nil
}

//Update updates the given habit. It returns an error if the habit does not exist
func (s *MemoryStore) Update(habit *Habit) error {
	if habit == nil {
		return ErrNilHabit
	}

	if _, ok := s.Habits[habit.Name]; !ok {
		return errors.New("cannot update habit does not exists")
	}

	s.Habits[habit.Name] = habit
	return nil
}

//GetAllHabits returns a []*Habits of all the stored habits
func (s MemoryStore) GetAllHabits() []*Habit {
	allHabits := make([]*Habit, 0, len(s.Habits))
	for _, h := range s.Habits {
		allHabits = append(allHabits, h)
	}
	return allHabits
}

//DBStore is a type that wraps a SQLite DB
type DBStore struct {
	db *sql.DB
}

//OpenDBStore opens a connection to the specified dbSource. It takes care of creating a habit table if it does not
//exist.
func OpenDBStore(dbSource string) (Store, error) {
	if dbSource == "" {
		return &DBStore{}, errors.New("empty dbSource string")
	}
	db, err := sql.Open("sqlite3", dbSource)
	if err != nil {
		return &DBStore{}, err
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
		return &DBStore{}, err
	}
	return &DBStore{db: db}, nil
}

//Get queries DBStore by name and returns the habit if it exists
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

//Create inserts the given habit into the store. It returns an error if the habit already exists
func (s *DBStore) Create(h *Habit) error {
	if h == nil {
		return ErrNilHabit
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

//Update updates the given habit. It returns an error if the habit does not exist
func (s *DBStore) Update(h *Habit) error {
	if h == nil {
		return ErrNilHabit
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

//GetAllHabits returns a []*Habits of all the stored habits
func (s *DBStore) GetAllHabits() []*Habit {

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

//FileStore is a type that wraps a JSON encoded file store
type FileStore struct {
	filename string
	habits   map[string]*Habit
}

//OpenFileStore reads the specified file and decodes its content into an unexported map[string]*Habit. It takes care of
//creating a new file if it does not exist.
func OpenFileStore(filename string) (Store, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return &FileStore{}, err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return &FileStore{}, err
	}
	habits := make(map[string]*Habit)

	if len(fileBytes) > 0 {
		err = json.Unmarshal(fileBytes, &habits)
		if err != nil {
			return &FileStore{}, err
		}
	}
	fileStore := FileStore{
		filename: filename,
		habits:   habits,
	}
	return &fileStore, nil
}

//Get searches FileStore by name and returns the habit if it exists
func (s *FileStore) Get(name string) (*Habit, error) {
	habit, ok := s.habits[name]
	if ok {
		return habit, nil
	}
	return nil, nil
}

//Create inserts the given habit into the store. It returns an error if the habit already exists. It triggers
//file io operations.
func (s *FileStore) Create(habit *Habit) error {
	if habit == nil {
		return ErrNilHabit
	}

	if _, ok := s.habits[habit.Name]; ok {
		return errors.New("habit already exists")
	}
	s.habits[habit.Name] = habit
	err := writeHabitsToFile(s.filename, s.habits)
	if err != nil {
		return err
	}

	return nil

}

//Update updates the given habit. It returns an error if the habit does not exist. It triggers file io operations.
func (s *FileStore) Update(habit *Habit) error {
	if habit == nil {
		return ErrNilHabit
	}

	if _, ok := s.habits[habit.Name]; !ok {
		return errors.New("cannot update habit does not exists")
	}

	s.habits[habit.Name] = habit
	err := writeHabitsToFile(s.filename, s.habits)
	return err
}

//GetAllHabits returns a []*Habits of all the stored habit
func (s *FileStore) GetAllHabits() []*Habit {
	allHabits := make([]*Habit, 0, len(s.habits))
	for _, h := range s.habits {
		allHabits = append(allHabits, h)
	}
	return allHabits
}

func writeHabitsToFile(filename string, habits map[string]*Habit) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	fileBytes, err := json.Marshal(habits)
	if err != nil {
		return err
	}
	err = file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = file.Write(fileBytes)
	return err
}

//ErrNilHabit is returned when a habit is nil
var ErrNilHabit = errors.New("habit cannot be nil")
