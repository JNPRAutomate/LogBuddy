package logbuddy

import (
	"fmt"
	"net/url"
	"os"
)

//FileStorage Allows data to be written to a file
type FileStorage struct {
	Location *url.URL     //location of the storage file
	MsgChan  chan Message //Channel to recieve messages from
	LogFile  *os.File     //File the file being used
}

//NewFileStorage Create an initialized NewFileStorage
func NewFileStorage(loc *url.URL) *FileStorage {
	return &FileStorage{Location: loc}
}

//Write Write data to the destination file
func (s *FileStorage) Write(data ...[]byte) error {
	for msg := range data {
		fmt.Println("DATA", data[msg])
		_, err := s.LogFile.Write(data[msg])
		if err != nil {
			return err
		}
	}

	s.LogFile.Sync()
	return nil
}

//Read Reads data from the
func (s *FileStorage) Read(data ...[]byte) error {
	return nil
}

//Open Open the file for writing
func (s *FileStorage) Open() error {
	var err error
	if _, err = os.Stat(s.Location.Path); os.IsNotExist(err) {
		s.LogFile, err = os.Create(s.Location.Path)
		if err != nil {
			return err
		}
	} else {
		s.LogFile, err = os.OpenFile(s.Location.Path, os.O_RDWR, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

//Close Close the file for writing
func (s *FileStorage) Close() error {
	err := s.LogFile.Close()
	return err
}

//SetDest Set the destination URL for writing
func (s *FileStorage) SetDest(dest string) error {
	return nil
}
