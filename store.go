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
