package DB

import (
	"log"
	"os"
)

type FileDB struct {
	FilePath string
}

func (db *FileDB) WriteFeedbackMessage(msg string) error {

	file, err := os.OpenFile(db.FilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("[Error] Can't close DB file at path %s with error %s\n", db.FilePath, err.Error())
		}
	}()

	stringToWrite := msg + "\n-------------------------------------------\n"

	if _, err = file.WriteString(stringToWrite); err != nil {
		return err
	}

	if err = file.Sync(); err != nil {
		return err
	}

	return nil
}
