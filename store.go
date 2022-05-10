package habit

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func (ht *Tracker) LoadFile(filename string) error {
	trackerFile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer trackerFile.Close()

	fileBytes, err := ioutil.ReadAll(trackerFile)
	if err != nil {
		return err
	}
	if len(fileBytes) > 0 {
		err = json.Unmarshal(fileBytes, ht)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ht *Tracker) WriteFile(filename string) error {
	trackerFile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer trackerFile.Close()

	fileBytes, err := json.Marshal(ht)
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

//things to be done:
//make generic Load, Write that wraps file and db
//create table if not exist
//insert new habit
//update habit

type Storable interface {
	Load() (Tracker, error)
	Write(tracker Tracker) error
}

type FileStore struct {
	filename string
}

func NewFileStore(filename string) FileStore {
	return FileStore{filename: filename}
}

func (s FileStore) Load() (Tracker, error) {
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
	return ht, nil
}

func (s FileStore) Write(tracker Tracker) error {
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
